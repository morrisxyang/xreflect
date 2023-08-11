package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Field ...
func Field(obj interface{}, fieldName string) (reflect.Value, error) {
	var empty reflect.Value
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}

	objValue := Value(obj)
	if !isSupportedKind(objValue.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}

	field := objValue.FieldByName(fieldName)
	if !field.IsValid() {
		return empty, fmt.Errorf("no such field: %s", fieldName)
	}

	return field, nil
}

// FieldValue ...
func FieldValue(obj interface{}, fieldName string) (interface{}, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// FieldKind ...
func FieldKind(obj interface{}, fieldName string) (reflect.Kind, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// FieldType returns the kind of the provided obj field.
// The `obj` can either be a structure or pointer to structure.
func FieldType(obj interface{}, fieldName string) (reflect.Type, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Type(), nil
}

// FieldTypeStr ...
func FieldTypeStr(obj interface{}, fieldName string) (string, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return "", err
	}

	return field.Type().String(), nil
}

// EmbedField ...
func EmbedField(obj interface{}, fieldPath string) (reflect.Value, error) {
	var empty reflect.Value
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}
	target := Value(obj)
	if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}
	if fieldPath == "" {
		return empty, errors.New("field path must not be empty")
	}

	fieldNames := strings.Split(fieldPath, ".")
	for _, fieldName := range fieldNames {
		if fieldName == "" {
			return empty, fmt.Errorf("field path:%s is invalid", fieldPath)
		}
		if target.Kind() == reflect.Pointer {
			target = reflect.ValueOf(target.Interface()).Elem()
		}
		if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
			return empty, fmt.Errorf("field %s is not struct", target)
		}

		target = target.FieldByName(fieldName)
		if !target.IsValid() {
			return empty, fmt.Errorf("no such field: %s", fieldName)
		}
	}

	return target, nil
}

// EmbedFieldValue ...
func EmbedFieldValue(obj interface{}, fieldPath string) (interface{}, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// EmbedFieldKind ...
func EmbedFieldKind(obj interface{}, fieldPath string) (reflect.Kind, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// EmbedFieldType ...
func EmbedFieldType(obj interface{}, fieldPath string) (reflect.Type, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return nil, err
	}
	return field.Type(), nil
}

// EmbedFieldTypeStr ...
func EmbedFieldTypeStr(obj interface{}, fieldPath string) (string, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return "", err
	}

	return field.Type().String(), nil
}

// Fields ...
func Fields(obj interface{}) ([]reflect.Value, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	objValue := Value(obj)
	if !isSupportedKind(objValue.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct")
	}

	var res []reflect.Value
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		res = append(res, field)
	}
	return res, nil
}

// SelectFields ...
func SelectFields(obj interface{}, f func(int, reflect.Value) bool) ([]reflect.Value, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	objValue := Value(obj)
	if !isSupportedKind(objValue.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct")
	}

	var res []reflect.Value
	for i := 0; i < objValue.NumField(); i++ {
		if f(i, objValue.Field(i)) {
			res = append(res, objValue.Field(i))
		}
	}
	return res, nil
}

// RangeFields ...
func RangeFields(obj interface{}, f func(int, reflect.Value) bool) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}

	objValue := Value(obj)
	if !isSupportedKind(objValue.Kind(), []reflect.Kind{reflect.Struct}) {
		return errors.New("obj must be struct")
	}

	for i := 0; i < objValue.NumField(); i++ {
		if !f(i, objValue.Field(i)) {
			break
		}
	}
	return nil
}
