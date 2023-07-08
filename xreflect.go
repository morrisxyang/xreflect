// Package xreflect ...
package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// SetEmbedStructField 设置嵌套的结构体字段, obj 必须是指针
func SetEmbedStructField(obj interface{}, path string, value interface{}) error {
	if obj == nil || path == "" {
		return nil
	}
	if reflect.TypeOf(obj).Kind() != reflect.Pointer {
		return nil
	}

	target := reflect.ValueOf(obj).Elem()
	fieldNames := strings.Split(path, ".")
	for _, fieldName := range fieldNames {
		if fieldName == "" {
			return fmt.Errorf("path:%s is invalid", path)
		}
		if target.Kind() == reflect.Pointer {
			target = reflect.ValueOf(target.Interface()).Elem()
		}
		if !target.IsValid() || !target.CanSet() || target.Kind() != reflect.Struct {
			return errors.New("set operation is invalid")
		}
		target = target.FieldByName(fieldName)
	}

	if !target.IsValid() || !target.CanSet() {
		return fmt.Errorf("cannot set %s field value", path)
	}

	targetType := target.Type()
	actualValue := reflect.ValueOf(value)
	if targetType != actualValue.Type() {
		actualValue = actualValue.Convert(targetType)
	}

	target.Set(actualValue)
	return nil
}
