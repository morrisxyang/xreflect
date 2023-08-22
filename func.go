package xreflect

import (
	"errors"
	"fmt"
	"reflect"
)

// CallFunc ...
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
	// var sliceType reflect.Type
	for i, arg := range args {
		// If the argument is nil, use zero value
		if arg == nil {
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
		// 如果函数最后一个返回值为 error 且不为空, 提取
		if errResult := retValues[len(retValues)-1].Interface(); errResult != nil {
			if err, ok := errResult.(error); ok {
				return retValues[0 : len(retValues)-1], err
			}
		}
	}

	return retValues, nil
}

// CallFuncSlice ...
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
	// var sliceType reflect.Type
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
		// 如果函数最后一个返回值为 error 且不为空, 提取
		if errResult := retValues[len(retValues)-1].Interface(); errResult != nil {
			if err, ok := errResult.(error); ok {
				return retValues[0 : len(retValues)-1], err
			}
		}
	}

	return retValues, nil
}

// CallMethod ...
func CallMethod(obj interface{}, method string, params ...interface{}) ([]reflect.Value, error) {
	if obj == nil {
		return nil, errors.New("fn must not be nil")
	}

	val := reflect.ValueOf(obj)
	methodValue := val.MethodByName(method)
	if methodValue.IsZero() {
		return nil, fmt.Errorf("method: %s not found", method)
	}

	return CallFunc(methodValue.Interface(), params...)
}

// CallMethodSlice ...
func CallMethodSlice(obj interface{}, method string, params ...interface{}) ([]reflect.Value, error) {
	if obj == nil {
		return nil, errors.New("fn must not be nil")
	}

	val := reflect.ValueOf(obj)
	methodValue := val.MethodByName(method)
	if methodValue.IsZero() {
		return nil, fmt.Errorf("method: %s not found", method)
	}

	return CallFuncSlice(methodValue.Interface(), params...)
}
