package xreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Field returns the reflect.Value of the provided obj field.
// The obj can either be a structure or pointer to structure.
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

// FieldValue returns the actual value of the provided obj field.
// The obj can either be a structure or pointer to structure.
func FieldValue(obj interface{}, fieldName string) (interface{}, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// FieldKind returns the reflect.Kind of the provided obj field.
// The obj can either be a structure or pointer to structure.
func FieldKind(obj interface{}, fieldName string) (reflect.Kind, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// FieldType returns the reflect.Type of the provided obj field.
// The obj can either be a structure or pointer to structure.
func FieldType(obj interface{}, fieldName string) (reflect.Type, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Type(), nil
}

// FieldTypeStr returns the string of reflect.Type of the provided obj field.
// The obj can either be a structure or pointer to structure.
func FieldTypeStr(obj interface{}, fieldName string) (string, error) {
	field, err := Field(obj, fieldName)
	if err != nil {
		return "", err
	}

	return field.Type().String(), nil
}

// EmbedField returns the reflect.Value of a field in the nested structure of obj based on the specified fieldPath.
// The obj can either be a structure or a pointer to a structure.
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

// EmbedFieldValue returns the actual value of a field in the nested structure of obj based on the specified fieldPath.
// The obj can either be a structure or a pointer to a structure.
func EmbedFieldValue(obj interface{}, fieldPath string) (interface{}, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// EmbedFieldKind returns the reflect.Kind of a field in the nested structure of obj based on the specified fieldPath.
// The obj can either be a structure or a pointer to a structure.
func EmbedFieldKind(obj interface{}, fieldPath string) (reflect.Kind, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return reflect.Invalid, err
	}

	return field.Kind(), nil
}

// EmbedFieldType returns the reflect.Type of a field in the nested structure of obj based on the specified fieldPath.
// The obj can either be a structure or a pointer to a structure.
func EmbedFieldType(obj interface{}, fieldPath string) (reflect.Type, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return nil, err
	}
	return field.Type(), nil
}

// EmbedFieldTypeStr returns the reflect.Type of a field in the nested structure of obj based on the specified fieldPath.
// The obj can either be a structure or a pointer to a structure.
func EmbedFieldTypeStr(obj interface{}, fieldPath string) (string, error) {
	field, err := EmbedField(obj, fieldPath)
	if err != nil {
		return "", err
	}

	return field.Type().String(), nil
}

// Fields returns a map of reflect.Value containing all the fields of the obj, with the field names as keys.
// The obj can either be a structure or a pointer to a structure.
func Fields(obj interface{}) (map[string]reflect.Value, error) {
	return fields(obj, false, "")
}

// FieldsDeep traverses the obj deeply, including all nested structures, and returns all fields as reflect.Value
// in the form of a map, where the key is the path of the field.
// The obj can either be a structure or pointer to structure.
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
			// deal struct
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

			// deal struct pointer
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

// SelectFields has the same functionality as Fields, but only the fields for which the function f returns true
// will be returned.
// The obj can either be a structure or pointer to structure.
func SelectFields(obj interface{}, f func(string, reflect.StructField, reflect.Value) bool) (map[string]reflect.Value,
	error) {
	return selectFields(obj, f, false, "")
}

// SelectFieldsDeep has the same functionality as FieldsDeep, but only the fields for which the function f returns true
// will be returned.
// The obj can either be a structure or pointer to structure.
func SelectFieldsDeep(obj interface{},
	f func(string, reflect.StructField, reflect.Value) bool) (map[string]reflect.Value,
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
			// deal struct
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

			// deal struct pointer
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

// RangeFields iterates over all fields of obj and calls function f on each field.
// If function f returns false, the iteration stops.
// The obj can either be a structure or pointer to structure.
func RangeFields(obj interface{}, f func(string, reflect.StructField, reflect.Value) bool) error {
	return rangeFields(obj, f, false, "")
}

// RangeFieldsDeep performs a deep traversal of obj and its nested structures, and calls function f on each field.
// If the function f returns false, the iteration stops.
// The obj can either be a structure or pointer to structure.
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
			// deal struct
			if cf.Kind() == reflect.Struct {
				err := rangeFields(cf.Interface(), f, deep, key)
				if err != nil {
					return err
				}
				continue
			}

			// deal struct pointer
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
