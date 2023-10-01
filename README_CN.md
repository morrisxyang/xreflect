# xreflect

[![Go Report Card](https://goreportcard.com/badge/github.com/morrisxyang/xreflect)](https://goreportcard.com/report/github.com/morrisxyang/xreflect)
[![Coverage Status](https://coveralls.io/repos/github/morrisxyang/xreflect/badge.svg)](https://coveralls.io/github/morrisxyang/xreflect)
![Static Badge](https://img.shields.io/badge/License-BSD2-Green)
![Static Badge](https://img.shields.io/badge/go%20verion-%3E%3D1.18-blue)

一个简单的, 易用的反射工具库.

主要支持如下特性:

- 设置结构体字段值, 支持通过路径比如`A.B.C`设置**嵌套结构体**字段的值

- 获取结构体字段的值, 类型, Tag 等.

- 遍历结构体所有字段, 支持 `select` 模式和 `range` 模式, 如果使用**深度遍历**方法比如 `FieldsDeep` 将遍历所有嵌套结构.

- 函数调用, 方法调用, 支持可变参数.

- 新建实例, 判断接口实现等等.

## 安装和文档

安装命令 `go get github.com/morrisxyang/xreflect`.

文档见 https://pkg.go.dev/github.com/morrisxyang/xreflect

## 快速开始

设置嵌套结构体字段值

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

获取 json tag 

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

筛选实例字段(深度遍历)

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


调用函数

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



## 核心方法

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

### `Field` 和 `StrcutField` 的区别是?

`Field` 返回  reflect.Value, `StrcutField` 返回 reflect.StrcutField. 

