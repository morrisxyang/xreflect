package xreflect

import (
	"errors"
	"fmt"
	"reflect"
)

// CallFunc 通过反射调用函数并返回结果.
// 支持可变参数, 底层使用 reflect.Value 的 Call 方法, 可变参数需拉平, 否则需使用 CallFuncSlice.
// 无需从返回值 []reflect.Value 中解析错误, 如果被调用函数最后一个返回值返回错误, 错误将被提取出来作为 CallFunc 的最后一个返回值返回.
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
			// 处理可变参数, 可变参数的类型是 Slice, 这里需要创建其元素类型的零值.
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
		// 如果函数最后一个返回值为 error 且不为空, 提取出来
		if errResult := retValues[len(retValues)-1].Interface(); errResult != nil {
			if err, ok := errResult.(error); ok {
				return retValues[0 : len(retValues)-1], err
			}
		}
	}

	return retValues, nil
}

// CallFuncSlice 功能同 CallFunc 一致, 支持可变参数, 但底层使用 reflect.Value 的 CallSlice 方法.
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
		// 如果函数最后一个返回值为 error 且不为空, 提取
		if errResult := retValues[len(retValues)-1].Interface(); errResult != nil {
			if err, ok := errResult.(error); ok {
				return retValues[0 : len(retValues)-1], err
			}
		}
	}

	return retValues, nil
}

// CallMethod 调用 obj 对象的 method 方法并返回结果, 支持可变参数.
// obj 对象必须是结构体或结构体指针.
// 底层使用 CallFunc, 详见 CallFunc.
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

// CallMethodSlice 作用同 CallMethod, 底层使用 CallFuncSlice, 详见 CallFuncSlice.
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
