package xreflect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Country struct {
	ID   int
	Name string

	City    City
	PtrCity *City
}

type City struct {
	PtrTown *Town
	Town    Town
}

type Town struct {
	Int  int
	Str  string
	Bool bool
	Strs []string
}

func newCountry() Country {
	town := Town{
		Int:  0,
		Str:  "Str",
		Bool: false,
		Strs: []string{"Str"},
	}

	city := City{
		Town:    town,
		PtrTown: &town,
	}

	country := Country{
		ID:      0,
		Name:    "A country",
		City:    city,
		PtrCity: &city,
	}
	return country
}

func TestNewInstance(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  interface{}
	}{
		{"int", int(1), 0},
		{"float", float32(1), float32(0)},
		{"string", "1", ""},
		{"struct", Country{ID: 1}, Country{}},
		{"struct ptr", &Country{ID: 1}, &Country{}},
		{"[]string", []string{"1"}, []string{}},
		{"1 array", [1]string{}, [1]string{}},
		{"2 array", [2]string{}, [2]string{}},
		{"map[string]string", make(map[string]string), make(map[string]string)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewInstance(tt.value), "NewInstance(%v)", tt.value)
		})
	}
}

func Test_SetEmbedStructField(t *testing.T) {
	// first level
	country := newCountry()
	err := SetEmbedStructField(&country, "ID", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.ID, 1)

	err = SetEmbedStructField(&country, "Name", "B country")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.Name, "B country")

	err = SetEmbedStructField(&country, "City", City{
		PtrTown: nil,
		Town:    Town{Int: 1},
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City, City{
		PtrTown: nil,
		Town:    Town{Int: 1},
	})

	err = SetEmbedStructField(&country, "PtrCity", &City{
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
	err = SetEmbedStructField(&country, "City.Town.Int", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.Town.Int, 1)

	err = SetEmbedStructField(&country, "City.Town.Str", "Now")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.Town.Str, "Now")

	err = SetEmbedStructField(&country, "City.Town.Bool", true)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.Town.Bool, true)

	err = SetEmbedStructField(&country, "City.Town.Strs", []string{"A", "B"})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.Town.Strs, []string{"A", "B"})

	// three level ptr
	country = newCountry()
	err = SetEmbedStructField(&country, "PtrCity.PtrTown.Int", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.PtrTown.Int, 1)

	err = SetEmbedStructField(&country, "PtrCity.PtrTown.Str", "Now")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.PtrTown.Str, "Now")

	err = SetEmbedStructField(&country, "PtrCity.PtrTown.Bool", true)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.PtrTown.Bool, true)

	err = SetEmbedStructField(&country, "PtrCity.PtrTown.Strs", []string{"A", "B"})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.PtrTown.Strs, []string{"A", "B"})

	// three level mix struct and ptr
	country = newCountry()
	err = SetEmbedStructField(&country, "City.PtrTown.Int", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.PtrTown.Int, 1)

	err = SetEmbedStructField(&country, "City.PtrTown.Str", "Now")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.PtrTown.Str, "Now")

	err = SetEmbedStructField(&country, "City.PtrTown.Bool", true)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.PtrTown.Bool, true)

	err = SetEmbedStructField(&country, "City.PtrTown.Strs", []string{"A", "B"})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.City.PtrTown.Strs, []string{"A", "B"})

	country = newCountry()
	err = SetEmbedStructField(&country, "PtrCity.Town.Int", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.Town.Int, 1)

	err = SetEmbedStructField(&country, "PtrCity.Town.Str", "Now")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.Town.Str, "Now")

	err = SetEmbedStructField(&country, "PtrCity.Town.Bool", true)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.Town.Bool, true)

	err = SetEmbedStructField(&country, "PtrCity.Town.Strs", []string{"A", "B"})
	assert.Equal(t, err, nil)
	assert.Equal(t, country.PtrCity.Town.Strs, []string{"A", "B"})
}
