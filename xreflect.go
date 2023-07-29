// Package xreflect 反射工具库
package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// NewInstance returns a new instance of the same type as the input value.
// The returned value will contain the zero value of the type.
func NewInstance(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	entity := reflect.ValueOf(value)

	switch entity.Kind() {
	case reflect.Ptr:
		entity = reflect.New(entity.Elem().Type())
		break
	case reflect.Chan:
		entity = reflect.MakeChan(entity.Type(), entity.Cap())
		break
	case reflect.Map:
		entity = reflect.MakeMap(entity.Type())
		break
	case reflect.Slice:
		entity = reflect.MakeSlice(entity.Type(), 0, entity.Cap())
		break
	default:
		entity = reflect.New(entity.Type()).Elem()
	}

	return entity.Interface()
}

// SetField 设置 field
func SetField(obj interface{}, fieldName string, fieldValue interface{}) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}
	if !isSupportedType(obj, []reflect.Kind{reflect.Pointer}) {
		return errors.New("obj must be pointer")
	}
	if fieldName == "" {
		return errors.New("field name must not be empty")
	}

	target := reflect.ValueOf(obj).Elem()
	target = target.FieldByName(fieldName)
	if err := checkField(target); err != nil {
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
		return errors.New("obj must be pointer")
	}
	if fieldName == "" {
		return errors.New("field name must not be empty")
	}

	target := reflect.ValueOf(obj).Elem()
	target = target.FieldByName(fieldName)
	if !target.IsValid() {
		return fmt.Errorf("field: %s is invalid", fieldName)
	}
	// private
	target = reflect.NewAt(target.Type(), unsafe.Pointer(target.UnsafeAddr())).Elem()

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
	if !isSupportedType(obj, []reflect.Kind{reflect.Pointer}) {
		return errors.New("obj must be pointer")
	}
	if fieldPath == "" {
		return errors.New("field path must not be empty")
	}

	target := reflect.ValueOf(obj).Elem()
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

		if err := checkField(target); err != nil {
			return err
		}
		if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
			return fmt.Errorf("field %s is not struct", target)
		}

		target = target.FieldByName(fieldName)
	}

	if err := checkField(target); err != nil {
		return err
	}

	actualValue := reflect.ValueOf(fieldValue)
	if target.Type() != actualValue.Type() {
		actualValue = actualValue.Convert(target.Type())
	}
	target.Set(actualValue)
	return nil
}

// GetField ...
func GetField(obj interface{}, name string) (reflect.Value, error) {
	var empty reflect.Value
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}

	objValue := GetValue(obj)
	if !isSupportedKind(objValue.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}

	field := objValue.FieldByName(name)
	if !field.IsValid() {
		return empty, fmt.Errorf("no such field: %s", name)
	}

	return field, nil
}

// GetFieldValue ...
func GetFieldValue(obj interface{}, name string) (interface{}, error) {
	field, err := GetField(obj, name)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// GetFieldKind ...
func GetFieldKind(obj interface{}, name string) (reflect.Kind, error) {
	field, err := GetField(obj, name)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// GetFieldType ...
func GetFieldType(obj interface{}, name string) (reflect.Type, error) {
	field, err := GetField(obj, name)
	if err != nil {
		return nil, err
	}

	return field.Type(), nil
}

// GetFieldTypeStr ...
func GetFieldTypeStr(obj interface{}, name string) (string, error) {
	field, err := GetField(obj, name)
	if err != nil {
		return "", err
	}

	return field.Type().String(), nil
}

// GetFieldTag returns the provided obj field tag value.
// The `obj` parameter can either be a structure or pointer to structure.
func GetFieldTag(obj interface{}, fieldName, tagKey string) (string, error) {
	if obj == nil {
		return "", errors.New("obj must not be nil")
	}

	objValue := GetValue(obj)
	if !isSupportedKind(objValue.Kind(), []reflect.Kind{reflect.Struct}) {
		return "", errors.New("obj must be struct")
	}

	structField, ok := objValue.Type().FieldByName(fieldName)
	if !ok {
		return "", fmt.Errorf("no such field: %s in obj", fieldName)
	}

	return structField.Tag.Get(tagKey), nil
}

// GetType ...
func GetType(obj interface{}) reflect.Type {
	if obj == nil {
		return nil
	}
	if v, ok := obj.(reflect.Type); ok {
		return v
	}
	if v, ok := obj.(reflect.Value); ok {
		return v.Type()
	}
	return reflect.TypeOf(obj)
}

// GetTypePenetrateElem ...
func GetTypePenetrateElem(obj interface{}) reflect.Type {
	if obj == nil {
		return nil
	}
	ty := GetType(obj)
	for ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	return ty
}

// GetValue ...
func GetValue(obj interface{}) reflect.Value {
	var empty reflect.Value
	if obj == nil {
		return empty
	}
	if v, ok := obj.(reflect.Value); ok {
		return v
	}
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		return reflect.ValueOf(obj).Elem()
	}
	return reflect.ValueOf(obj)
}

// GetValuePenetrateElem ...
func GetValuePenetrateElem(obj interface{}) reflect.Value {
	var empty reflect.Value
	if obj == nil {
		return empty
	}
	ty := GetValue(obj)
	for ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	return ty
}

func checkField(field reflect.Value) error {
	if !field.IsValid() {
		return fmt.Errorf("field %s is invalid", field)
	}
	if !field.CanSet() {
		return fmt.Errorf("field %s can not set", field)
	}

	return nil
}

func isSupportedKind(k reflect.Kind, kinds []reflect.Kind) bool {
	for _, v := range kinds {
		if k == v {
			return true
		}
	}

	return false
}

func isSupportedType(obj interface{}, types []reflect.Kind) bool {
	for _, t := range types {
		if reflect.TypeOf(obj).Kind() == t {
			return true
		}
	}

	return false
}
