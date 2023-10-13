// Package xreflect is a reflection utility library.
//
// The xreflect package aims to provide developers with high-level abstractions over the Go standard reflect library.
// This library's API is often considered low-level and unintuitive, making simple tasks like accessing structure
// field values or tags more complex than necessary.

package xreflect

import (
	"fmt"
	"reflect"
)

// NewInstance returns a new instance of the same type as the input value.
// The returned value will contain the zero value of the type.
// If obj type is a slice, chan, etc. , it will create an instance with the same capacity.
func NewInstance(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	entity := reflect.ValueOf(obj)

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

// Type returns the reflection type of obj.
// If obj is a pointer, it will be automatically dereferenced once.
func Type(obj interface{}) reflect.Type {
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

// TypePenetrateElem performs the same functionality as Type, but it will parse through all pointers
// until the final type is reached.
func TypePenetrateElem(obj interface{}) reflect.Type {
	if obj == nil {
		return nil
	}
	ty := Type(obj)
	for ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	return ty
}

// Value returns the reflection value of obj.
// If obj is a pointer, it will be automatically dereferenced once.
func Value(obj interface{}) reflect.Value {
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

// ValuePenetrateElem performs the same functionality as Value, but it will parse through all pointers
// until the final type is reached.
func ValuePenetrateElem(obj interface{}) reflect.Value {
	var empty reflect.Value
	if obj == nil {
		return empty
	}
	ty := Value(obj)
	for ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	return ty
}

// GetPkgPath returns the package path of obj.
func GetPkgPath(obj interface{}) string {
	ty := Type(obj)
	if ty == nil {
		return ""
	}
	return ty.PkgPath()
}

// Implements returns whether obj implements the given interface in.
func Implements(obj interface{}, in interface{}) bool {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return false
	}

	interfaceType := reflect.TypeOf(in).Elem()
	return objType.Implements(interfaceType)
}

// IsInterfaceNil returns whether obj is actually nil.
func IsInterfaceNil(obj interface{}) bool {
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		return value.IsNil()
	} else if value.Kind() == reflect.Invalid {
		return true
	}
	return false
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
