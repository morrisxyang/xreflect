package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// StructField returns the reflect.StructField of the provided obj field.
// The obj can either be a structure or pointer to structure.
func StructField(obj interface{}, fieldName string) (reflect.StructField, error) {
	var empty reflect.StructField
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}

	typ := Type(obj)
	if !isSupportedKind(typ.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}

	field, ok := typ.FieldByName(fieldName)
	if !ok {
		return empty, fmt.Errorf("no such field: %s in obj", fieldName)
	}
	return field, nil
}

// StructFieldKind returns the reflect.Kind of the provided obj field.
// The obj can either be a structure or pointer to structure.
func StructFieldKind(obj interface{}, fieldName string) (reflect.Kind, error) {
	field, err := StructField(obj, fieldName)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Type.Kind(), nil
}

// StructFieldType returns the reflect.Type of the provided obj field.
// The obj can either be a structure or pointer to structure.
func StructFieldType(obj interface{}, fieldName string) (reflect.Type, error) {
	field, err := StructField(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Type, nil
}

// StructFieldTypeStr returns the string of reflect.Type of the provided obj field.
// The obj can either be a structure or pointer to structure.
func StructFieldTypeStr(obj interface{}, fieldName string) (string, error) {
	field, err := StructField(obj, fieldName)
	if err != nil {
		return "", err
	}

	return field.Type.String(), nil
}

// HasStructField checks if the provided obj struct has field named fieldName.
// The obj can either be a structure or pointer to structure.
func HasStructField(obj interface{}, fieldName string) (bool, error) {
	_, err := StructField(obj, fieldName)
	if err != nil {
		return false, err
	}

	return true, nil
}

// StructFieldTag returns the reflect.StructTag of the provided obj field.
// The obj parameter can either be a structure or pointer to structure.
func StructFieldTag(obj interface{}, fieldName string) (reflect.StructTag, error) {
	structField, err := StructField(obj, fieldName)
	if err != nil {
		return "", err
	}

	return structField.Tag, nil
}

// StructFieldTagValue returns the provided obj field tag value.
// The obj parameter can either be a structure or pointer to structure.
func StructFieldTagValue(obj interface{}, fieldName, tagKey string) (string, error) {
	tag, err := StructFieldTag(obj, fieldName)
	if err != nil {
		return "", err
	}

	return tag.Get(tagKey), nil
}

// EmbedStructField returns the reflect.Value of a field in the
// nested structure of obj based on the specified fieldPath.
// The obj can either be a structure or a pointer to a structure.
func EmbedStructField(obj interface{}, fieldPath string) (reflect.StructField, error) {
	var empty reflect.StructField
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}
	if fieldPath == "" {
		return empty, errors.New("field path must not be empty")
	}

	target := Type(obj)
	if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}

	fieldNames := strings.Split(fieldPath, ".")
	for i, fieldName := range fieldNames {
		if fieldName == "" {
			return empty, fmt.Errorf("field path: %s is invalid", fieldPath)
		}
		structField, ok := target.FieldByName(fieldName)
		if !ok {
			return empty, fmt.Errorf("no such field: %s", fieldName)
		}
		target = structField.Type
		if i == len(fieldNames)-1 {
			return structField, nil
		}

		if target.Kind() == reflect.Pointer {
			target = target.Elem()
		}
		if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
			return empty, fmt.Errorf("field: %s is not struct", fieldName)
		}
	}
	return empty, nil
}

// EmbedStructFieldKind returns the reflect.Kind of a field in the
// nested structure of obj based on the specified fieldPath.
// The obj can either be a structure or a pointer to a structure.
func EmbedStructFieldKind(obj interface{}, fieldPath string) (reflect.Kind, error) {
	field, err := EmbedStructField(obj, fieldPath)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Type.Kind(), nil
}

// EmbedStructFieldType returns the reflect.Type of a field in the
// nested structure of obj based on the specified fieldPath.
// The obj can either be a structure or a pointer to a structure.
func EmbedStructFieldType(obj interface{}, fieldPath string) (reflect.Type, error) {
	field, err := EmbedStructField(obj, fieldPath)
	if err != nil {
		return nil, err
	}

	return field.Type, nil
}

// EmbedStructFieldTypeStr returns the string of the reflect.Type of a field in the nested structure
// of obj based on the specified fieldPath.
// The obj can either be a structure or a pointer to a structure.
func EmbedStructFieldTypeStr(obj interface{}, fieldPath string) (string, error) {
	field, err := EmbedStructField(obj, fieldPath)
	if err != nil {
		return "", err
	}

	return field.Type.String(), nil
}

// StructFields returns a slice of reflect.StructField containing all the fields of the obj.
// The obj can either be a structure or a pointer to a structure.
func StructFields(obj interface{}) ([]reflect.StructField, error) {
	return structFields(obj, false)
}

// StructFieldsFlatten returns "flattened" struct fields.
// Note that StructFieldsFlatten treats fields from anonymous inner structs as normal fields.
func StructFieldsFlatten(obj interface{}) ([]reflect.StructField, error) {
	return structFields(obj, true)
}

func structFields(obj interface{}, flatten bool) ([]reflect.StructField, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	typ := Type(obj)
	if !isSupportedKind(typ.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct")
	}

	var res []reflect.StructField
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !flatten {
			res = append(res, field)
			continue
		}

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			subFields, err := structFields(field.Type, flatten)
			if err != nil {
				return nil, fmt.Errorf("cannot get fields in %s: %w", field.Name, err)
			}
			res = append(res, subFields...)
		} else {
			res = append(res, field)
		}
	}
	return res, nil
}

// SelectStructFields has the same functionality as StructFields, but only the
// fields for which the function f returns true will be returned.
// The obj can either be a structure or pointer to structure.
func SelectStructFields(obj interface{}, f func(int, reflect.StructField) bool) ([]reflect.StructField, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	typ := Type(obj)
	if !isSupportedKind(typ.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct")
	}

	var res []reflect.StructField
	for i := 0; i < typ.NumField(); i++ {
		if f(i, typ.Field(i)) {
			res = append(res, typ.Field(i))
		}
	}
	return res, nil
}

// RangeStructFields iterates over all struct fields of obj and calls function f on each field.
// If function f returns false, the iteration stops.
// The obj can either be a structure or pointer to structure.
func RangeStructFields(obj interface{}, f func(int, reflect.StructField) bool) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}

	typ := Type(obj)
	if !isSupportedKind(typ.Kind(), []reflect.Kind{reflect.Struct}) {
		return errors.New("obj must be struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		if !f(i, typ.Field(i)) {
			break
		}
	}
	return nil
}

// AnonymousStructFields returns the slice of reflect.StructField of all anonymous fields in obj.
// The obj can either be a structure or pointer to structure.
func AnonymousStructFields(obj interface{}) ([]reflect.StructField, error) {
	return SelectStructFields(obj, func(i int, field reflect.StructField) bool {
		return field.Anonymous
	})
}
