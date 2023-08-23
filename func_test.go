package xreflect

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var addFunc = func(a int, b int) int {
	return a + b
}

var addReturnErrorFunc = func(a int, b int) (int, error) {
	return 0, errors.New("add error")
}

var onlyReturnErrorFunc = func() error {
	return errors.New("error")
}

var emptyFunc = func() {
	fmt.Println("empty")
	return
}

var changeMapFunc = func(m map[string]string) {
	m["1"] = "1"
}

var interfaceParamFunc = func(a, b interface{}) int {
	return a.(int) + b.(int)
}

var sliceParamFunc = func(ss []int) int {
	var sum int
	for _, s := range ss {
		sum += s
	}
	return sum
}

var pointFunc = func(p *int) {
	*p = 1
}

var pureVariadicFunc = func(ps ...*int) int {
	var sum int
	for _, p := range ps {
		sum += *p
	}
	return sum
}

var variadicFunc = func(pi *int, ps ...*int) int {
	var sum int
	for _, p := range ps {
		if p != nil {
			sum += *p
		}
	}
	if pi != nil {
		return sum + *pi
	}
	return sum
}

func TestCallFunc(t *testing.T) {
	_, err := CallFunc(nil, nil)
	assert.EqualError(t, err, "fn must not be nil")

	_, err = CallFunc("", nil)
	assert.EqualError(t, err, "fn must be func")

	_, err = CallFunc(addFunc, 1)
	assert.EqualError(t, err, "fn params num is 2, but got 1")

	_, err = CallFunc(addFunc, 1, 2, 3)
	assert.EqualError(t, err, "fn params num is 2, but got 3")

	res, err := CallFunc(addReturnErrorFunc, 1, 2)
	assert.EqualError(t, err, "add error")
	assert.Equal(t, 1, len(res))
	assert.Equal(t, 0, res[0].Interface())

	res, err = CallFunc(onlyReturnErrorFunc)
	assert.EqualError(t, err, "error")
	assert.Equal(t, 0, len(res))

	res, err = CallFunc(emptyFunc)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(res))

	res, err = CallFunc(addFunc, 1, 2)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, res[0].Interface())

	res, err = CallFunc(interfaceParamFunc, 1, 2)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, res[0].Interface())

	res, err = CallFunc(sliceParamFunc, []int{1, 2})
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, res[0].Interface())

	m := make(map[string]string)
	res, err = CallFunc(changeMapFunc, m)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(m))

	var p int
	res, err = CallFunc(pointFunc, &p)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, p)

	res, err = CallFunc(pureVariadicFunc)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, res[0].Interface())

	vp1 := 1
	vp2 := 2
	vp3 := 3
	res, err = CallFunc(pureVariadicFunc, &vp1, &vp2)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, res[0].Interface())

	res, err = CallFunc(pureVariadicFunc, &vp1, &vp2, &vp3)
	assert.Equal(t, nil, err)
	assert.Equal(t, 6, res[0].Interface())

	res, err = CallFunc(variadicFunc)
	assert.EqualError(t, err, "fn params num is 1 at least, but got 0")

	res, err = CallFunc(variadicFunc, &vp1, &vp2, &vp3)
	assert.Equal(t, nil, err)
	assert.Equal(t, 6, res[0].Interface())

	res, err = CallFunc(variadicFunc, &vp1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, res[0].Interface())

	res, err = CallFunc(variadicFunc, &vp1, nil, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, res[0].Interface())

	res, err = CallFunc(variadicFunc, nil, nil, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, res[0].Interface())

	res, err = CallFuncSlice(pureVariadicFunc)
	assert.EqualError(t, err, "use reflect.CallSlice, fn params num should be 1, but got 0")

	res, err = CallFuncSlice(variadicFunc, &vp1, []*int{nil, nil, nil})
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, res[0].Interface())

	res, err = CallFuncSlice(variadicFunc, &vp1, []*int{&vp2, &vp3})
	assert.Equal(t, nil, err)
	assert.Equal(t, 6, res[0].Interface())

}

type A struct {
	Int int
}

func (a A) GetInt() int {
	return a.Int
}

func (a *A) AddOne() int {
	a.Int += 1
	return a.Int
}

func (a *A) AddInts(is ...int) int {
	for _, i := range is {
		a.Int += i
	}
	return a.Int
}

func TestCallMethod(t *testing.T) {
	_, err := CallMethod(nil, "")
	assert.EqualError(t, err, "obj must not be nil")

	_, err = CallMethod(&Person{}, "ABC")
	assert.EqualError(t, err, "method: ABC not found")

	res, err := CallMethod(&A{1}, "GetInt")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, res[0].Interface())

	res, err = CallMethod(A{1}, "GetInt")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, res[0].Interface())

	res, err = CallMethod(&A{1}, "AddOne")
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, res[0].Interface())

	res, err = CallMethod(A{1}, "AddOne")
	assert.EqualError(t, err, "method: AddOne not found")

	res, err = CallMethod(&A{1}, "AddInts")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, res[0].Interface())

	res, err = CallMethod(&A{1}, "AddInts", 1, 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, res[0].Interface())

	is := []interface{}{1, 1}
	res, err = CallMethod(&A{1}, "AddInts", is...)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, res[0].Interface())

	_, err = CallMethodSlice(nil, "")
	assert.EqualError(t, err, "obj must not be nil")

	_, err = CallMethodSlice(&Person{}, "ABC")
	assert.EqualError(t, err, "method: ABC not found")

	_, err = CallMethodSlice(&A{1}, "AddInts")
	assert.EqualError(t, err, "use reflect.CallSlice, fn params num should be 1, but got 0")

	res, err = CallMethodSlice(&A{1}, "AddInts", []int{1, 1})
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, res[0].Interface())
}
