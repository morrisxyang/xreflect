package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// GetField ...
func GetField(obj interface{}, fieldName string) (reflect.Value, error) {
	var empty reflect.Value
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}

	objValue := GetValue(obj)
	if !isSupportedKind(objValue.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}

	field := objValue.FieldByName(fieldName)
	if !field.IsValid() {
		return empty, fmt.Errorf("no such field: %s", fieldName)
	}

	return field, nil
}

// GetFieldValue ...
func GetFieldValue(obj interface{}, fieldName string) (interface{}, error) {
	field, err := GetField(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// GetFieldKind ...
func GetFieldKind(obj interface{}, fieldName string) (reflect.Kind, error) {
	field, err := GetField(obj, fieldName)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// GetFieldType ...
func GetFieldType(obj interface{}, fieldName string) (reflect.Type, error) {
	field, err := GetField(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Type(), nil
}

// GetFieldTypeStr ...
func GetFieldTypeStr(obj interface{}, fieldName string) (string, error) {
	field, err := GetField(obj, fieldName)
	if err != nil {
		return "", err
	}

	return field.Type().String(), nil
}

// GetEmbedField ...
func GetEmbedField(obj interface{}, fieldPath string) (reflect.Value, error) {
	var empty reflect.Value
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}
	target := GetValue(obj)
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

// GetEmbedFieldValue ...
func GetEmbedFieldValue(obj interface{}, fieldPath string) (interface{}, error) {
	field, err := GetEmbedField(obj, fieldPath)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// GetEmbedFieldKind ...
func GetEmbedFieldKind(obj interface{}, fieldPath string) (reflect.Kind, error) {
	field, err := GetEmbedField(obj, fieldPath)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// GetEmbedFieldType ...
func GetEmbedFieldType(obj interface{}, fieldPath string) (reflect.Type, error) {
	field, err := GetEmbedField(obj, fieldPath)
	if err != nil {
		return nil, err
	}

	return field.Type(), nil
}

// GetEmbedFieldTypeStr ...
func GetEmbedFieldTypeStr(obj interface{}, fieldPath string) (string, error) {
	field, err := GetEmbedField(obj, fieldPath)
	if err != nil {
		return "", err
	}

	return field.Type().String(), nil
}
