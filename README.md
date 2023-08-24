# xreflect

![Static Badge](https://img.shields.io/badge/License-BSD2-Green)
![Static Badge](https://img.shields.io/badge/go%20verion-%3E%3D1.15-blue)

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

Call a function

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

## Core Methods


## FAQ