package xreflect

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	Heart  interface{}
	Person struct {
		Name      string  `json:"name"`
		Age       int     `json:"age"`
		PtrPerson *Person `json:"ptr_person,omitempty"`

		// private field
		phone string `json:"phone"`

		// anonymous fields
		string
		int
		*Person
		Country
		Heart // interface
		inner
	}
	inner struct {
		innerString string
	}

	Country struct {
		ID   int
		Name string

		City    City
		PtrCity *City
	}

	City struct {
		ID      int
		PtrTown *Town
		Town    Town
	}

	Town struct {
		Int  int
		Str  string
		Bool bool
		Strs []string
	}
)

func newCountry() Country {
	town := Town{
		Int:  0,
		Str:  "Str",
		Bool: false,
		Strs: []string{"Str"},
	}

	city := City{
		ID:      0,
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
			result := Type(tc.obj)
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
			result := TypePenetrateElem(tc.obj)
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
		{name: "reflect.Value",
			input:    reflect.ValueOf(1),
			expected: reflect.ValueOf(1),
		},
		{name: "Nil",
			input:    nil,
			expected: reflect.Value{},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Value(tc.input)
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
		{"nil", nil, reflect.Value{}},
		{"***int", i3, reflect.ValueOf(i1).Elem()},
		{"int", i0, reflect.ValueOf(i0)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ValuePenetrateElem(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, actual)
			}
		})
	}
}

func F() {

}

func TestGetPkgPath(t *testing.T) {

	tests := []struct {
		name string
		obj  interface{}
		want string
	}{
		{
			name: "struct",
			obj:  &Person{},
			want: "github.com/morrisxyang/xreflect",
		},
		{
			name: "[]string{}",
			obj:  []string{},
			want: "",
		},
		{
			name: "Func",
			obj:  F,
			want: "",
		},
		{
			name: "assert.NoError",
			obj:  assert.NoError,
			want: "",
		},
		{
			name: "io.EOF",
			obj:  io.EOF,
			want: "errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetPkgPath(tt.obj), "GetPkgPath(%v)", tt.obj)
		})
	}
}

func TestImplements(t *testing.T) {
	type error interface {
		Error() string
	}

	type args struct {
		obj interface{}
		in  interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Nil",
			args: args{
				obj: nil,
				in:  nil,
			},
			want: false,
		},
		{
			name: "error",
			args: args{
				obj: errors.New(""),
				in:  (*error)(nil),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Implements(tt.args.obj, tt.args.in), "Implements(%v, %v)", tt.args.obj, tt.args.in)
		})
	}
}

func TestIsInterfaceNil(t *testing.T) {
	var inter interface{} = (*int)(nil)
	var inter1 interface{} = nil
	t.Run("inter not nil", func(t *testing.T) {
		assert.Equal(t, false, inter == nil)
	})
	t.Run("inter is nil", func(t *testing.T) {
		assert.Equal(t, true, IsInterfaceNil(inter))
	})
	t.Run("inter1 is nil", func(t *testing.T) {
		assert.Equal(t, true, IsInterfaceNil(inter1))
	})
}
