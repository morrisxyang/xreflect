package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// GetStructField 获取结构体的字段
func GetStructField(obj interface{}, fieldName string) (reflect.StructField, error) {
	var empty reflect.StructField
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}

	ty := GetType(obj)
	if !isSupportedKind(ty.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}

	field, ok := ty.FieldByName(fieldName)
	if !ok {
		return empty, fmt.Errorf("no such field: %s in obj", fieldName)
	}
	return field, nil
}

// GetStructFieldKind ...
func GetStructFieldKind(obj interface{}, fieldName string) (reflect.Kind, error) {
	field, err := GetStructField(obj, fieldName)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Type.Kind(), nil
}

// GetStructFieldType ...
func GetStructFieldType(obj interface{}, fieldName string) (reflect.Type, error) {
	field, err := GetStructField(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Type, nil
}

// GetStructFieldTypeStr ...
func GetStructFieldTypeStr(obj interface{}, fieldName string) (string, error) {
	field, err := GetStructField(obj, fieldName)
	if err != nil {
		return "", err
	}

	return field.Type.String(), nil
}

// HasField checks if the provided `obj` struct has field named `name`.
// The `obj` can either be a structure or pointer to structure.
func HasField(obj interface{}, fieldName string) (bool, error) {
	_, err := GetStructField(obj, fieldName)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetStructFieldTag returns the provided obj field tag value.
// The `obj` parameter can either be a structure or pointer to structure.
func GetStructFieldTag(obj interface{}, fieldName, tagKey string) (string, error) {
	structField, err := GetStructField(obj, fieldName)
	if err != nil {
		return "", err
	}

	return structField.Tag.Get(tagKey), nil
}

// GetStructFields 获取结构体的字段
func GetStructFields(obj interface{}) ([]reflect.StructField, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	ty := GetType(obj)
	if !isSupportedKind(ty.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct")
	}

	var res []reflect.StructField
	for i := 0; i < ty.NumField(); i++ {
		res = append(res, ty.Field(i))
	}
	return res, nil
}

// SelectStructFields ...
func SelectStructFields(obj interface{}, f func(int, reflect.StructField) bool) ([]reflect.StructField, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	ty := GetType(obj)
	if !isSupportedKind(ty.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct")
	}

	var res []reflect.StructField
	for i := 0; i < ty.NumField(); i++ {
		if f(i, ty.Field(i)) {
			res = append(res, ty.Field(i))
		}
	}
	return res, nil
}

// RangeStructFields ...
func RangeStructFields(obj interface{}, f func(int, reflect.StructField) bool) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}

	ty := GetType(obj)
	if !isSupportedKind(ty.Kind(), []reflect.Kind{reflect.Struct}) {
		return errors.New("obj must be struct")
	}

	for i := 0; i < ty.NumField(); i++ {
		if !f(i, ty.Field(i)) {
			break
		}
	}
	return nil
}

// GetAnonymousStructFields 获取匿名结构体字段
func GetAnonymousStructFields(obj interface{}) ([]reflect.StructField, error) {
	return SelectStructFields(obj, func(i int, field reflect.StructField) bool {
		return field.Anonymous
	})
}

// GetEmbedStructField ...
func GetEmbedStructField(obj interface{}, fieldPath string) (reflect.StructField, error) {
	var empty reflect.StructField
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}
	target := GetType(obj)
	if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}
	if fieldPath == "" {
		return empty, errors.New("field path must not be empty")
	}

	fieldNames := strings.Split(fieldPath, ".")
	for i, fieldName := range fieldNames {
		if fieldName == "" {
			return empty, fmt.Errorf("field path:%s is invalid", fieldPath)
		}
		if target.Kind() == reflect.Pointer {
			target = target.Elem()
		}
		if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
			return empty, fmt.Errorf("field %s is not struct", target)
		}

		structField, ok := target.FieldByName(fieldName)
		if !ok {
			return empty, fmt.Errorf("no such field: %s", fieldName)
		}
		target = structField.Type
		if i == len(fieldNames)-1 {
			return structField, nil
		}
	}
	return empty, nil
}

// GetEmbedStructFieldKind ...
func GetEmbedStructFieldKind(obj interface{}, fieldPath string) (reflect.Kind, error) {
	field, err := GetEmbedStructField(obj, fieldPath)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Type.Kind(), nil
}

// GetEmbedStructFieldType ...
func GetEmbedStructFieldType(obj interface{}, fieldPath string) (reflect.Type, error) {
	field, err := GetEmbedStructField(obj, fieldPath)
	if err != nil {
		return nil, err
	}

	return field.Type, nil
}

// GetEmbedStructFieldTypeStr ...
func GetEmbedStructFieldTypeStr(obj interface{}, fieldPath string) (string, error) {
	field, err := GetEmbedStructField(obj, fieldPath)
	if err != nil {
		return "", err
	}

	return field.Type.String(), nil
}
