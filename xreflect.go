// Package xreflect 反射工具库
package xreflect

import (
	"fmt"
	"reflect"
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
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		return reflect.TypeOf(obj).Elem()
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

func checkField(field reflect.Value, name string) error {
	if !field.IsValid() {
		return fmt.Errorf("field: %s is invalid", name)
	}
	if !field.CanSet() {
		return fmt.Errorf("field: %s can not set", name)
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
