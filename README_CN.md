# xreflect

![Static Badge](https://img.shields.io/badge/License-BSD2-Green)
![Static Badge](https://img.shields.io/badge/go%20verion-%3E%3D1.15-blue)

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

调用函数

```go
var addFunc = func(nums ...int) int {
		var sum int
		for _, num := range nums {
			sum += num
		}
		return sum
}

res, _ := CallFunc(addFunc1, 1, 2, 3)

// 6
fmt.Println(res[0].Interface())
```



## 核心方法


## FAQ

