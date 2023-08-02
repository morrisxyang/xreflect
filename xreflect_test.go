package xreflect

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	PtrPerson *Person

	// private field
	phone string `json:"phone"`

	// anonymous fields
	int
	string
	*Person
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
	err := SetEmbedField(&country, "ID", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, country.ID, 1)

	err = SetEmbedField(&country, "Name", "B country")
	assert.Equal(t, err, nil)
	assert.Equal(t, country.Name, "B country")

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

func TestGetStructField(t *testing.T) {
	_, err := GetStructField(nil, "Name")
	assert.EqualError(t, err, "obj must not be nil")

	_, err = GetStructField("", "Name")
	assert.EqualError(t, err, "obj must be struct")

	p := &Person{}
	_, err = GetStructField(p, "Name1")
	assert.EqualError(t, err, "no such field: Name1 in obj")

	st, err := GetStructField(p, "Name")
	assert.Equal(t, nil, err)
	assert.Equal(t, "Name", st.Name)

	st, err = GetStructField(p, "phone")
	assert.Equal(t, nil, err)
	assert.Equal(t, "phone", st.Name)
	assert.Equal(t, "github.com/morrisxyang/xreflect", st.PkgPath)

	st, err = GetStructField(p, "int")
	assert.Equal(t, nil, err)
	assert.Equal(t, "int", st.Name)
	assert.Equal(t, "github.com/morrisxyang/xreflect", st.PkgPath)

	st, err = GetStructField(p, "string")
	assert.Equal(t, nil, err)
	assert.Equal(t, "string", st.Name)
	assert.Equal(t, "github.com/morrisxyang/xreflect", st.PkgPath)

	st, err = GetStructField(p, "Person")
	assert.Equal(t, nil, err)
	assert.Equal(t, "Person", st.Name)
	assert.Equal(t, "", st.PkgPath)
}

func TestGetStructFieldTag(t *testing.T) {
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
			got, err := GetStructFieldTag(tt.args.obj, tt.args.fieldName, tt.args.tagKey)
			if !tt.wantErr(t, err, fmt.Sprintf("GetStructFieldTag(%v, %v, %v)", tt.args.obj, tt.args.fieldName, tt.args.tagKey)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetStructFieldTag(%v, %v, %v)", tt.args.obj, tt.args.fieldName, tt.args.tagKey)
		})
	}
}

func TestGetType(t *testing.T) {
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
		{
			name:     "Testing with nil",
			obj:      nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetType(tc.obj)
			if result != tc.expected {
				t.Errorf("Expected type %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGetTypePenetrateElem(t *testing.T) {
	var i3 ***int
	i0 := 1
	i1 := &i0
	i2 := &i1
	i3 = &i2

	testCases := []struct {
		name     string
		obj      interface{}
		expected reflect.Type
	}{
		{
			name:     "***int and *int",
			obj:      i3,
			expected: reflect.TypeOf(i1).Elem(),
		},
		{
			name:     "***int and int",
			obj:      i3,
			expected: reflect.TypeOf(i0),
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
		{
			name:     "Testing with nil",
			obj:      nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetTypePenetrateElem(tc.obj)
			if result != tc.expected {
				t.Errorf("Expected type %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	ii := &[]int{1, 2, 3}
	testCases := []struct {
		name     string
		input    interface{}
		expected reflect.Value
	}{
		{name: "Int",
			input:    42,
			expected: reflect.ValueOf(42),
		},
		{name: "String",
			input:    "hello",
			expected: reflect.ValueOf("hello"),
		},
		{name: "&[]int{1, 2, 3}",
			input:    ii,
			expected: reflect.ValueOf(ii).Elem(),
		},
		{name: "Nil",
			input:    nil,
			expected: reflect.Value{},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetValue(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected reflect value %v, but got %v", tc.expected, actual)
			}
		})
	}
}

func TestGetValuePenetrateElem(t *testing.T) {
	var i3 ***int
	i0 := 1
	i1 := &i0
	i2 := &i1
	i3 = &i2

	testCases := []struct {
		name     string
		input    interface{}
		expected reflect.Value
	}{
		{"***int", i3, reflect.ValueOf(i1).Elem()},
		{"int", i0, reflect.ValueOf(i0)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetValuePenetrateElem(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, actual)
			}
		})
	}
}
