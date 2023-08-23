package xreflect

import (
	"errors"
	"fmt"
	"reflect"
)

// CallFunc invokes a function using reflection and returns the result.
// It supports variadic arguments and uses the Call method of reflect.Value underneath.
// Variadic arguments need to be flattened, otherwise should use CallFuncSlice.
// There is no need to parse errors from the returned []reflect.Value. If the called function's last return value
// is an error, it will be extracted and returned as the last return value of CallFunc.
func CallFunc(fn interface{}, args ...interface{}) ([]reflect.Value, error) {
	if fn == nil {
		return nil, errors.New("fn must not be nil")
	}

	typ := Type(fn)
	val := Value(fn)
	if !isSupportedKind(val.Kind(), []reflect.Kind{reflect.Func}) {
		return nil, errors.New("fn must be func")
	}
	if !typ.IsVariadic() && len(args) != typ.NumIn() {
		return nil, fmt.Errorf("fn params num is %d, but got %d", typ.NumIn(), len(args))
	}
	if typ.IsVariadic() {
		if len(args) < typ.NumIn()-1 {
			return nil, fmt.Errorf("fn params num is %d at least, but got %d", typ.NumIn()-1, len(args))
		}
	}

	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		// If the argument is nil, use zero value
		if arg == nil {
			// Handling variadic parameters, whose type is a slice.
			// We need to create the zero value of its element type here.
			if typ.IsVariadic() && i >= typ.NumIn()-1 {
				reflectArgs[i] = reflect.New(typ.In(typ.NumIn() - 1).Elem()).Elem()
			} else {
				reflectArgs[i] = reflect.New(typ.In(i)).Elem()
			}
		} else {
			reflectArgs[i] = reflect.ValueOf(arg)
		}
	}

	var retValues []reflect.Value
	retValues = val.Call(reflectArgs)

	if len(retValues) > 0 {
		// If the last return value of the function is an error and not empty, extract it
		if errResult := retValues[len(retValues)-1].Interface(); errResult != nil {
			if err, ok := errResult.(error); ok {
				return retValues[0 : len(retValues)-1], err
			}
		}
	}

	return retValues, nil
}

// CallFuncSlice has the same functionality as CallFunc, it must have variadic parameters, but it uses
// the CallSlice method of reflect.Value under the hood.
// CallSlice calls the variadic function v with the input arguments in,
// assigning the slice in[len(in)-1] to v's final variadic argument.
// For example, if len(in) == 3, v.CallSlice(in) represents the Go call v(in[0], in[1], in[2]...).
func CallFuncSlice(fn interface{}, args ...interface{}) ([]reflect.Value, error) {
	if fn == nil {
		return nil, errors.New("fn must not be nil")
	}

	typ := Type(fn)
	val := Value(fn)
	if !isSupportedKind(val.Kind(), []reflect.Kind{reflect.Func}) {
		return nil, errors.New("fn must be func")
	}
	if !typ.IsVariadic() {
		return nil, errors.New("fn must be variadic")
	}

	if len(args) != typ.NumIn() {
		return nil, fmt.Errorf("use reflect.CallSlice, fn params num should be %d, but got %d",
			typ.NumIn(), len(args))
	}

	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		// If the argument is nil, use zero value
		if arg == nil {
			reflectArgs[i] = reflect.New(typ.In(i)).Elem()
		} else {
			reflectArgs[i] = reflect.ValueOf(arg)
		}
	}

	var retValues []reflect.Value
	retValues = val.CallSlice(reflectArgs)

	if len(retValues) > 0 {
		// If the last return value of the function is an error and not empty, extract it.
		if errResult := retValues[len(retValues)-1].Interface(); errResult != nil {
			if err, ok := errResult.(error); ok {
				return retValues[0 : len(retValues)-1], err
			}
		}
	}

	return retValues, nil
}

// CallMethod calls the method `method` of the `obj` object and returns the result, supporting variadic parameters.
// The `obj` object must be a struct or a struct pointer.
// It internally uses CallFunc, see CallFunc for more details.
func CallMethod(obj interface{}, method string, params ...interface{}) ([]reflect.Value, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	typ := Type(obj)
	if !isSupportedKind(typ.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct or struct pointer")
	}

	val := reflect.ValueOf(obj)
	methodValue := val.MethodByName(method)
	if !methodValue.IsValid() {
		return nil, fmt.Errorf("method: %s not found", method)
	}

	return CallFunc(methodValue.Interface(), params...)
}

// CallMethodSlice has the same functionality as CallMethod, and it uses CallFuncSlice as its underlying implementation.
// For more details, refer to the CallFuncSlice documentation.
func CallMethodSlice(obj interface{}, method string, params ...interface{}) ([]reflect.Value, error) {
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}

	typ := Type(obj)
	if !isSupportedKind(typ.Kind(), []reflect.Kind{reflect.Struct}) {
		return nil, errors.New("obj must be struct or struct pointer")
	}

	val := reflect.ValueOf(obj)
	methodValue := val.MethodByName(method)
	if !methodValue.IsValid() {
		return nil, fmt.Errorf("method: %s not found", method)
	}

	return CallFuncSlice(methodValue.Interface(), params...)
}
