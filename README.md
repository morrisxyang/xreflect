# xreflect

[![Go Report Card](https://goreportcard.com/badge/github.com/morrisxyang/xreflect)](https://goreportcard.com/report/github.com/morrisxyang/xreflect)
[![Coverage Status](https://coveralls.io/repos/github/morrisxyang/xreflect/badge.svg)](https://coveralls.io/github/morrisxyang/xreflect)
![Static Badge](https://img.shields.io/badge/License-BSD2-Green)
![Static Badge](https://img.shields.io/badge/go%20verion-%3E%3D1.18-blue)

A simple and user-friendly reflection utility library.

The xreflect package aims to provide developers with high-level abstractions over the Go standard reflect library.
This library's API is often considered low-level and unintuitive, making simple tasks like setting structure
field values more complex than necessary.

The main features supported are:

- Setting the values of structure fields, supporting **nested structure** field values by using paths such as `A.B.C`.

- Getting the values, types, tags, etc., of structure fields.

- Traversing all fields of a structure, supporting both `select` mode and `range` mode. If a **deep traversal** method like `FieldsDeep` is used, it will traverse all nested structures.

- Function calls and method calls, supporting variadic parameters.

- Creating new instances, checking interface implementations, and more.

## Installation and Docs

Install using `go get github.com/morrisxyang/xreflect`.

Full documentation is available at https://pkg.go.dev/github.com/morrisxyang/xreflect

## Quick Start

Set nested struct field value

```go
person := &Person{
	Name: "John",
	Age:  20,
	Country: Country{
		ID:   0,
		Name: "Perk",
	},
}

_ = SetEmbedField(person, "Country.ID", 1)

// Perk's ID: 1 
fmt.Printf("Perk's ID: %d \n", person.Country.ID)
```

Find json tag

```go
type Person struct {
	Name string `json:"name" xml:"_name"`
}
p := &Person{}
// json:"name" xml:"_name"
fmt.Println(StructFieldTag(p, "Name"))
// name <nil>
fmt.Println(StructFieldTagValue(p, "Name", "json"))
// _name <nil>
fmt.Println(StructFieldTagValue(p, "Name", "xml"))
```

Filter instance fields (deep traversal)

```go
type Person struct {
	id   string
	Age  int    `json:"int"`
	Name string `json:"name"`
	Home struct {
		Address string `json:"address"`
	}
}

p := &Person{}
fields, _ := SelectFieldsDeep(p, func(s string, field reflect.StructField, value reflect.Value) bool {
	return field.Tag.Get("json") != ""
})
// key: Age type: int
// key: Name type: string
// key: Home.Address type: string
for k, v := range fields {
	fmt.Printf("key: %s type: %v\n", k, v.Type())
}
```


Call a function

```go
var addFunc = func(nums ...int) int {
		var sum int
		for _, num := range nums {
			sum += num
		}
		return sum
}

res, _ := CallFunc(addFunc, 1, 2, 3)

// 6
fmt.Println(res[0].Interface())
```

## Core Methods


### FieldX

- [func Field(obj interface{}, fieldName string) (reflect.Value, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#Field)
- [func FieldValue(obj interface{}, fieldName string) (interface{}, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#FieldValue)
- [func EmbedField(obj interface{}, fieldPath string) (reflect.Value, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#EmbedField)
- [func EmbedFieldValue(obj interface{}, fieldPath string) (interface{}, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#EmbedFieldValue)
- [func Fields(obj interface{}) (map[string]reflect.Value, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#Fields)
- [func FieldsDeep(obj interface{}) (map[string]reflect.Value, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#FieldsDeep)
- [func RangeFields(obj interface{}, f func(string, reflect.StructField, reflect.Value) bool) error](https://pkg.go.dev/github.com/morrisxyang/xreflect#RangeFields)
- [func SelectFields(obj interface{}, f func(string, reflect.StructField, reflect.Value) bool) (map[string]reflect.Value, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#SelectFields)
- etc.

### SetX

- [func SetEmbedField(obj interface{}, fieldPath string, fieldValue interface{}) error](https://pkg.go.dev/github.com/morrisxyang/xreflect#SetEmbedField)

- [func SetField(obj interface{}, fieldName string, fieldValue interface{}) error](https://pkg.go.dev/github.com/morrisxyang/xreflect#SetField)
- [func SetPrivateField(obj interface{}, fieldName string, fieldValue interface{}) error](https://pkg.go.dev/github.com/morrisxyang/xreflect#SetPrivateField)
- etc.

### StrcutFieldX

- [func StructField(obj interface{}, fieldName string) (reflect.StructField, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#StructField)
- [func StructFieldTagValue(obj interface{}, fieldName, tagKey string) (string, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#StructFieldTagValue)
- [func EmbedStructField(obj interface{}, fieldPath string) (reflect.StructField, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#EmbedStructField)

- [func StructFields(obj interface{}) ([]reflect.StructField, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#StructFields)
- [func StructFieldsFlatten(obj interface{}) ([]reflect.StructField, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#StructFieldsFlatten)

- [func SelectStructFields(obj interface{}, f func(int, reflect.StructField) bool) ([]reflect.StructField, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#SelectStructFields)

- [func RangeStructFields(obj interface{}, f func(int, reflect.StructField) bool) error](https://pkg.go.dev/github.com/morrisxyang/xreflect#RangeStructFields)

- etc.

### FuncX

- [func CallFunc(fn interface{}, args ...interface{}) ([]reflect.Value, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#CallFunc)
- [func CallMethod(obj interface{}, method string, params ...interface{}) ([]reflect.Value, error)](https://pkg.go.dev/github.com/morrisxyang/xreflect#CallMethod)

- etc.

### Others

- [func NewInstance(obj interface{}) interface{}](https://pkg.go.dev/github.com/morrisxyang/xreflect#NewInstance)
- [func Type(obj interface{}) reflect.Type](https://pkg.go.dev/github.com/morrisxyang/xreflect#Type)
- [func TypePenetrateElem(obj interface{}) reflect.Type](https://pkg.go.dev/github.com/morrisxyang/xreflect#TypePenetrateElem)
- [func Value(obj interface{}) reflect.Value](https://pkg.go.dev/github.com/morrisxyang/xreflect#Value)
- [func ValuePenetrateElem(obj interface{}) reflect.Value](https://pkg.go.dev/github.com/morrisxyang/xreflect#ValuePenetrateElem)
- etc.

## FAQ

### What is the difference between `Field` and `StructField`?

`Field` returns `reflect.Value`, while `StructField` returns `reflect.StructField`.
