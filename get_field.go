package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Field ...
func Field(obj interface{}, fieldName string) (reflect.Value, error) {
	var empty reflect.Value
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}

	val := Value(obj)
	if !isSupportedKind(val.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return empty, fmt.Errorf("no such field: %s", fieldName)
	}

	return field, nil
}

// FieldValue ...
func FieldValue(obj interface{}, fieldName string) (interface{}, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// FieldKind ...
func FieldKind(obj interface{}, fieldName string) (reflect.Kind, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// FieldType returns the kind of the provided obj field.
// The `obj` can either be a structure or pointer to structure.
func FieldType(obj interface{}, fieldName string) (reflect.Type, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Type(), nil
}

// FieldTypeStr ...
func FieldTypeStr(obj interface{}, fieldName string) (string, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return "", err
	}

	return field.Type().String(), nil
}

// EmbedField ...
func EmbedField(obj interface{}, fieldPath string) (reflect.Value, error) {
	var empty reflect.Value
	if obj == nil {
		return empty, errors.New("obj must not be nil")
	}
	if fieldPath == "" {
		return empty, errors.New("field path must not be empty")
	}

	target := Value(obj)
	if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
		return empty, errors.New("obj must be struct")
	}

	fieldNames := strings.Split(fieldPath, ".")
	for i, fieldName := range fieldNames {
		if fieldName == "" {
			return empty, fmt.Errorf("field path:%s is invalid", fieldPath)
		}
		target = target.FieldByName(fieldName)
		if !target.IsValid() {
			return empty, fmt.Errorf("no such field: %s", fieldName)
		}

		if i == len(fieldNames)-1 {
			break
		}
		if target.Kind() == reflect.Pointer {
			if target.IsNil() {
				return empty, fmt.Errorf("field: %s is nil", fieldName)
			}
			target = reflect.ValueOf(target.Interface()).Elem()
		}
		if !isSupportedKind(target.Kind(), []reflect.Kind{reflect.Struct}) {
			return empty, fmt.Errorf("field: %s is not struct", fieldName)
		}
	}
	return target, nil
}

// EmbedFieldValue ...
func EmbedFieldValue(obj interface{}, fieldPath string) (interface{}, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// EmbedFieldKind ...
func EmbedFieldKind(obj interface{}, fieldPath string) (reflect.Kind, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// EmbedFieldType ...
func EmbedFieldType(obj interface{}, fieldPath string) (reflect.Type, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return nil, err
	}
	return field.Type(), nil
}

// EmbedFieldTypeStr ...
func EmbedFieldTypeStr(obj interface{}, fieldPath string) (string, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return "", err
	}

	return field.Type().String(), nil
}

// Fields ...
func Fields(obj interface{}) (map[string]reflect.Value, error) {
	return fields(obj, false, "")
}

// FieldsDeep ...
func FieldsDeep(obj interface{}) (map[string]reflect.Value, error) {
	return fields(obj, true, "")
}

func fields(obj interface{}, deep bool, prefix string) (map[string]reflect.Value, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	typ := Type(obj)
	val := Value(obj)
	if !isSupportedKind(typ.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct")
	}

	res := make(map[string]reflect.Value)
	for i := 0; i < typ.NumField(); i++ {
		ct := typ.Field(i)
		cf := val.Field(i)

		key := ct.Name
		if prefix != "" {
			key = prefix + "." + ct.Name
		}
		res[key] = cf

		if deep {
			// struct
			if cf.Kind() == reflect.Struct {
				m, err := fields(cf.Interface(), deep, key)
				if err != nil {
					return nil, err
				}
				for k, v := range m {
					res[k] = v
				}
				continue
			}

			// struct pointer
			if cf.Kind() == reflect.Ptr && !cf.IsNil() {
				cf = cf.Elem()
				m, err := fields(cf.Interface(), deep, key)
				if err != nil {
					return nil, err
				}
				for k, v := range m {
					res[k] = v
				}
				continue
			}
		}
	}
	return res, nil
}

// SelectFields ...
func SelectFields(obj interface{}, f func(string, reflect.StructField, reflect.Value) bool) (map[string]reflect.Value,
	error) {
	return selectFields(obj, f, false, "")
}

// SelectFieldsDeep ...
func SelectFieldsDeep(obj interface{}, f func(string, reflect.StructField,
	reflect.Value) bool) (map[string]reflect.Value,
	error) {
	return selectFields(obj, f, true, "")
}

func selectFields(obj interface{}, f func(string, reflect.StructField, reflect.Value) bool,
	deep bool, prefix string) (map[string]reflect.Value, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	typ := Type(obj)
	val := Value(obj)
	if !isSupportedKind(typ.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct")
	}

	res := make(map[string]reflect.Value)
	for i := 0; i < typ.NumField(); i++ {
		ct := typ.Field(i)
		cf := val.Field(i)

		key := ct.Name
		if prefix != "" {
			key = prefix + "." + ct.Name
		}
		if f(key, ct, cf) {
			res[key] = cf
		}

		if deep {
			// struct
			if cf.Kind() == reflect.Struct {
				m, err := selectFields(cf.Interface(), f, deep, key)
				if err != nil {
					return nil, err
				}
				for k, v := range m {
					res[k] = v
				}
				continue
			}

			// struct pointer
			if cf.Kind() == reflect.Ptr && !cf.IsNil() {
				cf = cf.Elem()
				m, err := selectFields(cf.Interface(), f, deep, key)
				if err != nil {
					return nil, err
				}
				for k, v := range m {
					res[k] = v
				}
				continue
			}
		}
	}
	return res, nil
}

// RangeFields ...
func RangeFields(obj interface{}, f func(string, reflect.StructField, reflect.Value) bool) error {
	return rangeFields(obj, f, false, "")
}

// RangeFieldsDeep ...
func RangeFieldsDeep(obj interface{}, f func(string, reflect.StructField, reflect.Value) bool) error {
	return rangeFields(obj, f, true, "")
}

func rangeFields(obj interface{}, f func(string, reflect.StructField, reflect.Value) bool,
	deep bool, prefix string) error {
	if obj == nil {
		return errors.New("obj must not be nil")
	}

	typ := Type(obj)
	val := Value(obj)
	if !isSupportedKind(typ.Kind(), []reflect.Kind{reflect.Struct}) {
		return errors.New("obj must be struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		ct := typ.Field(i)
		cf := val.Field(i)

		key := ct.Name
		if prefix != "" {
			key = prefix + "." + ct.Name
		}
		if !f(key, ct, cf) {
			return nil
		}

		if deep {
			// struct
			if cf.Kind() == reflect.Struct {
				err := rangeFields(cf.Interface(), f, deep, key)
				if err != nil {
					return err
				}
				continue
			}

			// struct pointer
			if cf.Kind() == reflect.Ptr && !cf.IsNil() {
				cf = cf.Elem()
				err := rangeFields(cf.Interface(), f, deep, key)
				if err != nil {
					return err
				}
				continue
			}
		}
	}
	return nil
}
