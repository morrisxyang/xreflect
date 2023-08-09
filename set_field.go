package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// SetField 设置 field
func SetField(obj interface{}, fieldName string, fieldValue interface{}) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}
	if !isSupportedType(obj, []reflect.Kind{reflect.Pointer}) {
		return errors.New("obj must be struct pointer")
	}
	if fieldName == "" {
		return errors.New("field name must not be empty")
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

// SetPrivateField 设置私有字段
func SetPrivateField(obj interface{}, fieldName string, fieldValue interface{}) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}
	if !isSupportedType(obj, []reflect.Kind{reflect.Pointer}) {
		return errors.New("obj must be struct pointer")
	}
	if fieldName == "" {
		return errors.New("field name must not be empty")
	}

	target := Value(obj)
	if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
		return errors.New("obj must be struct pointer")
	}

	target = target.FieldByName(fieldName)
	if !target.IsValid() {
		return fmt.Errorf("field: %s is invalid", fieldName)
	}
	// private field
	target = reflect.NewAt(target.Type(), unsafe.Pointer(target.UnsafeAddr())).Elem()

	actualValue := reflect.ValueOf(fieldValue)
	if target.Type() != actualValue.Type() {
		actualValue = actualValue.Convert(target.Type())
	}
	target.Set(actualValue)
	return nil
}

// SetEmbedField 设置嵌套的结构体字段, obj 必须是指针
func SetEmbedField(obj interface{}, fieldPath string, fieldValue interface{}) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}
	if !isSupportedType(obj, []reflect.Kind{reflect.Pointer}) {
		return errors.New("obj must be pointer")
	}
	if fieldPath == "" {
		return errors.New("field path must not be empty")
	}

	target := Value(obj)
	fieldNames := strings.Split(fieldPath, ".")
	for _, fieldName := range fieldNames {
		if fieldName == "" {
			return fmt.Errorf("field path:%s is invalid", fieldPath)
		}
		if target.Kind() == reflect.Pointer {
			// 	结构体指针为空则自行创建
			if target.IsNil() {
				target.Set(reflect.New(target.Type().Elem()).Elem().Addr())
			}
			target = reflect.ValueOf(target.Interface()).Elem()
		}
		if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
			return fmt.Errorf("field %s is not struct", target)
		}

		target = target.FieldByName(fieldName)
		if err := checkField(target, fieldName); err != nil {
			return err
		}
	}

	actualValue := reflect.ValueOf(fieldValue)
	if target.Type() != actualValue.Type() {
		actualValue = actualValue.Convert(target.Type())
	}
	target.Set(actualValue)
	return nil
}
