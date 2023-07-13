// Package xreflect ...
package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

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

	targetType := target.Type()
	actualValue := reflect.ValueOf(fieldValue)
	if targetType != actualValue.Type() {
		actualValue = actualValue.Convert(targetType)
	}

	target.Set(actualValue)
	return nil
}
