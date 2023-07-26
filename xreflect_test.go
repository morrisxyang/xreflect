package xreflect

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	phone string `json:"phone"`
}

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
	s := "1"
	tests := []struct {
		name  string
		value interface{}
		want  interface{}
	}{
		{"int", int(1), 0},
		{"float", float32(1), float32(0)},
		{"complex", complex(1, 1), complex(0, 0)},
		{"string", "1", ""},
		{"struct", Country{ID: 1}, Country{}},
		{"struct ptr", &Country{ID: 1}, &Country{}},
		{"[]string", []string{"1"}, []string{}},
		{"[]*string", []*string{&s}, []*string{}},
		{"1 array", [1]string{}, [1]string{}},
		{"2 array", [2]string{}, [2]string{}},
		{"map[string]string", make(map[string]string), make(map[string]string)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewInstance(tt.value), "NewInstance(%v)", tt.value)
		})
	}

	// test chan
	ci1 := make(chan int)
	ci2 := NewInstance(ci1).(chan int)
	assert.Equal(t, 0, cap(ci1))
	assert.Equal(t, 0, cap(ci2))
	go func() {
		assert.Equal(t, 1, <-ci2)
	}()
	ci2 <- 1

	ci3 := make(chan int, 3)
	ci4 := NewInstance(ci3).(chan int)
	assert.Equal(t, 3, cap(ci4))
	assert.Equal(t, 0, len(ci4))
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

func TestGetField(t *testing.T) {
	type args struct {
		obj  interface{}
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get",
			args: args{
				obj:  Person{Name: "John", Age: 30},
				name: "Name",
			},
			want:    "John",
			wantErr: assert.NoError,
		},
		{
			name: "No such field",
			args: args{
				obj:  Person{Name: "John", Age: 30},
				name: "Address",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "no such field: Address")
			},
		},
		{
			name: "nil",
			args: args{
				obj:  nil,
				name: "Name",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "obj must not be nil")
			},
		},
		{
			name: "Not a struct",
			args: args{
				obj:  "test",
				name: "Name",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "obj must be struct")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFieldValue(tt.args.obj, tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("GetFieldValue(%v, %v)", tt.args.obj, tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetFieldValue(%v, %v)", tt.args.obj, tt.args.name)
		})
	}
}

func TestGetFieldTag(t *testing.T) {
	type args struct {
		obj       interface{}
		fieldName string
		tagKey    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Struct json tag",
			args: args{
				obj:       Person{},
				fieldName: "Name",
				tagKey:    "json",
			},
			want:    "name",
			wantErr: assert.NoError,
		},
		{
			name: "Struct ptr json tag",
			args: args{
				obj:       &Person{},
				fieldName: "Name",
				tagKey:    "json",
			},
			want:    "name",
			wantErr: assert.NoError,
		},
		{
			name: "Struct no exist field",
			args: args{
				obj:       &Person{},
				fieldName: "Name1",
				tagKey:    "json",
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "no such field: Name1 in obj")
			},
		},
		{
			name: "Struct no exist tag",
			args: args{
				obj:       &Person{},
				fieldName: "Name",
				tagKey:    "json1",
			},
			want:    "",
			wantErr: assert.NoError,
		},
		{
			name: "Struct private tag",
			args: args{
				obj:       &Person{},
				fieldName: "phone",
				tagKey:    "json",
			},
			want:    "phone",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFieldTag(tt.args.obj, tt.args.fieldName, tt.args.tagKey)
			if !tt.wantErr(t, err, fmt.Sprintf("GetFieldTag(%v, %v, %v)", tt.args.obj, tt.args.fieldName, tt.args.tagKey)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetFieldTag(%v, %v, %v)", tt.args.obj, tt.args.fieldName, tt.args.tagKey)
		})
	}
}

func Test_getType(t *testing.T) {
	testCases := []struct {
		name     string
		obj      interface{}
		expected reflect.Type
	}{
		{
			name:     "Testing with reflect.Type",
			obj:      reflect.TypeOf("test"),
			expected: reflect.TypeOf("test"),
		},
		{
			name:     "Testing with reflect.Value",
			obj:      reflect.ValueOf(10),
			expected: reflect.TypeOf(10),
		},
		{
			name:     "Testing with other types",
			obj:      "test",
			expected: reflect.TypeOf("test"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getType(tc.obj)
			if result != tc.expected {
				t.Errorf("Expected type %v, got %v", tc.expected, result)
			}
		})
	}

	var i ***int
	assert.Equal(t, reflect.Ptr, getType(i).Kind())
	assert.Equal(t, "***int", getType(i).String())

	ty := getType(i)
	for ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	assert.Equal(t, "int", ty.String())
}
