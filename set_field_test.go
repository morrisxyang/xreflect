package xreflect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetField(t *testing.T) {
	p := &Person{
		Name:  "",
		Age:   0,
		phone: "",
	}
	err := SetField(nil, "Name", "John")
	assert.EqualError(t, err, "obj must not be nil")

	err = SetField(*p, "Name", "John")
	assert.EqualError(t, err, "obj must be struct pointer")

	err = SetField(p, "Name1", "John")
	assert.EqualError(t, err, "field: Name1 is invalid")

	s := "str"
	err = SetField(&s, "Name", "John")
	assert.EqualError(t, err, "obj must be struct pointer")

	err = SetField(p, "phone", "123")
	assert.EqualError(t, err, "field: phone can not set")

	err = SetField(p, "Name", "John")
	assert.Equal(t, err, nil)
	assert.Equal(t, p.Name, "John")

	err = SetField(p, "PtrPerson", &Person{
		Name: "Mike",
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, p.PtrPerson.Name, "Mike")
}

func TestSetPrivateField(t *testing.T) {
	p := &Person{
		Name:  "",
		Age:   0,
		phone: "",
	}
	err := SetPrivateField(nil, "Name", "John")
	assert.EqualError(t, err, "obj must not be nil")

	err = SetPrivateField(*p, "Name", "John")
	assert.EqualError(t, err, "obj must be struct pointer")

	err = SetPrivateField(p, "Name1", "John")
	assert.EqualError(t, err, "field: Name1 is invalid")

	s := "str"
	err = SetPrivateField(&s, "Name", "John")
	assert.EqualError(t, err, "obj must be struct pointer")

	err = SetPrivateField(p, "phone", "123")
	assert.Equal(t, err, nil)
	assert.Equal(t, p.phone, "123")

	err = SetPrivateField(p, "Name", "John")
	assert.Equal(t, err, nil)
	assert.Equal(t, p.Name, "John")

	err = SetPrivateField(p, "PtrPerson", &Person{
		Name: "Mike",
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, p.PtrPerson.Name, "Mike")
}

func TestSetEmbedField(t *testing.T) {
	// first level
	country := newCountry()
	err := SetEmbedField(&country, "ID1", 1)
	assert.EqualError(t, err, "field: ID1 is invalid")

	err = SetEmbedField(nil, "ID", 1)
	assert.EqualError(t, err, "obj must not be nil")

	err = SetEmbedField("123", "ID", 1)
	assert.EqualError(t, err, "obj must be pointer")

	err = SetEmbedField(&country, "", 1)
	assert.EqualError(t, err, "field path must not be empty")

	err = SetEmbedField(&country, ".Town.Int", 1)
	assert.EqualError(t, err, "field path:.Town.Int is invalid")

	err = SetEmbedField(&country, "ID", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.ID, 1)

	err = SetEmbedField(&country, "Name", "B country")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.Name, "B country")

	type myString string
	err = SetEmbedField(&country, "Name", myString("C country"))
	assert.Equal(t, err, nil)
	assert.Equal(t, country.Name, "C country")

	err = SetEmbedField(&country, "City", City{
		PtrTown: nil,
		Town:    Town{Int: 1},
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City, City{
		PtrTown: nil,
		Town:    Town{Int: 1},
	})

	err = SetEmbedField(&country, "PtrCity", &City{
		PtrTown: nil,
		Town:    Town{Int: 1},
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity, &City{
		PtrTown: nil,
		Town:    Town{Int: 1},
	})

	// three level struct
	country = newCountry()
	err = SetEmbedField(&country, "City.Town.Int", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.Town.Int, 1)

	err = SetEmbedField(&country, "City.Town.Str", "Now")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.Town.Str, "Now")

	err = SetEmbedField(&country, "City.Town.Bool", true)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.Town.Bool, true)

	err = SetEmbedField(&country, "City.Town.Strs", []string{"A", "B"})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.Town.Strs, []string{"A", "B"})

	// three level ptr
	c := &Country{}
	err = SetEmbedField(c, "PtrCity.PtrTown.Int", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, c.PtrCity.PtrTown.Int, 1)

	country = newCountry()
	err = SetEmbedField(&country, "PtrCity.PtrTown.Int", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.PtrTown.Int, 1)

	err = SetEmbedField(&country, "PtrCity.PtrTown.Str", "Now")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.PtrTown.Str, "Now")

	err = SetEmbedField(&country, "PtrCity.PtrTown.Bool", true)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.PtrTown.Bool, true)

	err = SetEmbedField(&country, "PtrCity.PtrTown.Strs", []string{"A", "B"})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.PtrTown.Strs, []string{"A", "B"})

	// three level mix struct and ptr
	country = newCountry()
	err = SetEmbedField(&country, "City.PtrTown.Int", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.PtrTown.Int, 1)

	err = SetEmbedField(&country, "City.PtrTown.Str", "Now")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.PtrTown.Str, "Now")

	err = SetEmbedField(&country, "City.PtrTown.Bool", true)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.PtrTown.Bool, true)

	err = SetEmbedField(&country, "City.PtrTown.Strs", []string{"A", "B"})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.PtrTown.Strs, []string{"A", "B"})

	country = newCountry()
	err = SetEmbedField(&country, "PtrCity.Town.Int", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.Town.Int, 1)

	err = SetEmbedField(&country, "PtrCity.Town.Str", "Now")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.Town.Str, "Now")

	err = SetEmbedField(&country, "PtrCity.Town.Bool", true)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.Town.Bool, true)

	err = SetEmbedField(&country, "PtrCity.Town.Strs", []string{"A", "B"})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.Town.Strs, []string{"A", "B"})
}
