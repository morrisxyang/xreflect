// Package xreflect ...
package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// SetField 设置 field
func SetField(obj interface{}, fieldName string, fieldValue interface{}) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}
	if fieldName == "" {
		return errors.New("field name must not be empty")
	}
	if reflect.TypeOf(obj).Kind() != reflect.Pointer {
		return errors.New("obj must be pointer")
	}

	target := reflect.ValueOf(obj).Elem()
	target = target.FieldByName(fieldName)
	if !target.IsValid() {
		return fmt.Errorf("field: %s is invalid", fieldName)
	}
	if !target.CanSet() {
		return fmt.Errorf("field: %s cannot set", fieldName)
	}

	actualValue := reflect.ValueOf(fieldValue)
	if target.Type() != actualValue.Type() {
		actualValue = actualValue.Convert(target.Type())
	}
	target.Set(actualValue)
	return nil
}

// SetEmbedStructField 设置嵌套的结构体字段, obj 必须是指针
func SetEmbedStructField(obj interface{}, fieldPath string, fieldValue interface{}) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}
	if fieldPath == "" {
		return errors.New("field path must not be empty")
	}
	if reflect.TypeOf(obj).Kind() != reflect.Pointer {
		return errors.New("obj must be pointer")
	}

	target := reflect.ValueOf(obj).Elem()
	fieldNames := strings.Split(fieldPath, ".")
	for _, fieldName := range fieldNames {
		if fieldName == "" {
			return fmt.Errorf("field path:%s is invalid", fieldPath)
		}
		if target.Kind() == reflect.Pointer {
			if target.IsNil() {
				target.Set(reflect.New(target.Type().Elem()).Elem().Addr())
			}
			target = reflect.ValueOf(target.Interface()).Elem()
		}
		if !target.IsValid() || !target.CanSet() || target.Kind() != reflect.Struct {
			return errors.New("set operation is invalid")
		}
		target = target.FieldByName(fieldName)
	}

	if !target.IsValid() || !target.CanSet() {
		return fmt.Errorf("%s cannot be set", fieldPath)
	}

	actualValue := reflect.ValueOf(fieldValue)
	if target.Type() != actualValue.Type() {
		actualValue = actualValue.Convert(target.Type())
	}
	target.Set(actualValue)
	return nil
}
