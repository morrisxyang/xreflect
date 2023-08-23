package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// SetField sets the fieldName field of the obj object according to the fieldValue parameter.
// The obj can either be a structure or a pointer to a structure.
// The type of fieldValue must be compatible with the type of the fieldName field, otherwise it will panic.
func SetField(obj interface{}, fieldName string, fieldValue interface{}) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}
	if fieldName == "" {
		return errors.New("field name must not be empty")
	}

	if !isSupportedType(obj, []reflect.Kind{reflect.Pointer}) {
		return errors.New("obj must be struct pointer")
	}
	target := Value(obj)
	if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
		return errors.New("obj must be struct pointer")
	}

	target = target.FieldByName(fieldName)
	if err := checkField(target, fieldName); err != nil {
		return err
	}

	actualValue := reflect.ValueOf(fieldValue)
	if target.Type() != actualValue.Type() {
		actualValue = actualValue.Convert(target.Type())
	}
	target.Set(actualValue)
	return nil
}

// SetPrivateField is similar to SetField, but it allows you to set private fields of an object.
// The obj can be either a structure or a pointer to a structure.
func SetPrivateField(obj interface{}, fieldName string, fieldValue interface{}) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}
	if fieldName == "" {
		return errors.New("field name must not be empty")
	}

	if !isSupportedType(obj, []reflect.Kind{reflect.Pointer}) {
		return errors.New("obj must be struct pointer")
	}
	target := Value(obj)
	if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
		return errors.New("obj must be struct pointer")
	}

	target = target.FieldByName(fieldName)
	if !target.IsValid() {
		return fmt.Errorf("field: %s is invalid", fieldName)
	}
	// deal private field
	target = reflect.NewAt(target.Type(), unsafe.Pointer(target.UnsafeAddr())).Elem()

	actualValue := reflect.ValueOf(fieldValue)
	if target.Type() != actualValue.Type() {
		actualValue = actualValue.Convert(target.Type())
	}
	target.Set(actualValue)
	return nil
}

// SetEmbedField sets a nested struct field using fieldPath. The rest of the functionality is the same as SetField.
// For example, fieldPath can be "FieldA.FieldB.FieldC", where FieldA and FieldB must be structures or pointers
// to structures. If FieldB does not exist, it will be automatically created.
// The obj can either be a structure or pointer to structure.
func SetEmbedField(obj interface{}, fieldPath string, fieldValue interface{}) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}
	if fieldPath == "" {
		return errors.New("field path must not be empty")
	}

	if !isSupportedType(obj, []reflect.Kind{reflect.Pointer}) {
		return errors.New("obj must be struct pointer")
	}
	target := Value(obj)
	if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
		return errors.New("obj must be struct pointer")
	}

	fieldNames := strings.Split(fieldPath, ".")
	for i, fieldName := range fieldNames {
		if fieldName == "" {
			return fmt.Errorf("field path:%s is invalid", fieldPath)
		}
		target = target.FieldByName(fieldName)
		if err := checkField(target, fieldName); err != nil {
			return err
		}

		if i == len(fieldNames)-1 {
			break
		}
		if target.Kind() == reflect.Pointer {
			// If the structure pointer is nil, create it.
			if target.IsNil() {
				target.Set(reflect.New(target.Type().Elem()).Elem().Addr())
			}
			target = reflect.ValueOf(target.Interface()).Elem()
		}
		if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
			return fmt.Errorf("field: %s is not struct", fieldName)
		}
	}

	actualValue := reflect.ValueOf(fieldValue)
	if target.Type() != actualValue.Type() {
		actualValue = actualValue.Convert(target.Type())
	}
	target.Set(actualValue)
	return nil
}
