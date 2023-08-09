package xreflect

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStructFieldXMethods(t *testing.T) {
	_, err := StructField(nil, "Name")
	assert.EqualError(t, err, "obj must not be nil")

	_, err = StructField("", "Name")
	assert.EqualError(t, err, "obj must be struct")

	p := &Person{}
	_, err = StructField(p, "Name1")
	assert.EqualError(t, err, "no such field: Name1 in obj")

	st, err := StructField(p, "Name")
	assert.Equal(t, nil, err)
	assert.Equal(t, "Name", st.Name)

	st, err = StructField(p, "phone")
	assert.Equal(t, nil, err)
	assert.Equal(t, "phone", st.Name)
	assert.Equal(t, "github.com/morrisxyang/xreflect", st.PkgPath)

	st, err = StructField(p, "int")
	assert.Equal(t, nil, err)
	assert.Equal(t, "int", st.Name)
	assert.Equal(t, "github.com/morrisxyang/xreflect", st.PkgPath)

	st, err = StructField(p, "string")
	assert.Equal(t, nil, err)
	assert.Equal(t, "string", st.Name)
	assert.Equal(t, "github.com/morrisxyang/xreflect", st.PkgPath)

	st, err = StructField(p, "Person")
	assert.Equal(t, nil, err)
	assert.Equal(t, "Person", st.Name)
	assert.Equal(t, "", st.PkgPath)

	k, err := StructFieldKind(p, "Name")
	assert.Equal(t, nil, err)
	assert.Equal(t, reflect.String, k)

	ty, err := StructFieldType(p, "Age")
	assert.Equal(t, nil, err)
	assert.Equal(t, reflect.Int, ty.Kind())

	ts, err := StructFieldTypeStr(p, "Age")
	assert.Equal(t, nil, err)
	assert.Equal(t, "int", ts)

	b, err := HasStructField(p, "Age")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, b)

	_, err = StructFields(nil)
	assert.EqualError(t, err, "obj must not be nil")

	_, err = StructFields("123")
	assert.EqualError(t, err, "obj must be struct")

	sfs, err := StructFields(p)
	assert.Equal(t, nil, err)
	assert.Equal(t, 7, len(sfs))

	sfs, err = SelectStructFields(nil, nil)
	assert.EqualError(t, err, "obj must not be nil")

	sfs, err = SelectStructFields("123", nil)
	assert.EqualError(t, err, "obj must be struct")

	sfs, err = SelectStructFields(p, func(i int, field reflect.StructField) bool {
		return true
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 7, len(sfs))

	sfs, err = AnonymousStructFields(p)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(sfs))

	err = RangeStructFields(nil, nil)
	assert.EqualError(t, err, "obj must not be nil")

	err = RangeStructFields("123", nil)
	assert.EqualError(t, err, "obj must be struct")

	err = RangeStructFields(p, func(i int, field reflect.StructField) bool {
		return true
	})
	assert.Equal(t, nil, err)

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
			got, err := StructFieldTagValue(tt.args.obj, tt.args.fieldName, tt.args.tagKey)
			if !tt.wantErr(t, err, fmt.Sprintf("StructFieldTagValue(%v, %v, %v)", tt.args.obj, tt.args.fieldName, tt.args.tagKey)) {
				return
			}
			assert.Equalf(t, tt.want, got, "StructFieldTagValue(%v, %v, %v)", tt.args.obj, tt.args.fieldName, tt.args.tagKey)
		})
	}
}

func TestGetEmbedStructFieldXMethods(t *testing.T) {
	to := Town{
		Int:  1,
		Str:  "Town",
		Bool: true,
		Strs: []string{"Str"},
	}
	ci := City{
		PtrTown: &to,
		Town:    to,
	}
	ct := Country{
		ID:      1,
		Name:    "Country",
		City:    ci,
		PtrCity: &ci,
	}

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
				obj:  ct,
				name: "Name",
			},
			want:    "Country",
			wantErr: assert.NoError,
		},
		{
			name: "No such field",
			args: args{
				obj:  ct,
				name: "Address",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.EqualError(t, err, "no such field: Address")
				return false
			},
		},
		{
			name: "Nil",
			args: args{
				obj:  nil,
				name: "Name",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.EqualError(t, err, "obj must not be nil")
				return false
			},
		},
		{
			name: "Not a struct",
			args: args{
				obj:  "test",
				name: "Name",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.EqualError(t, err, "obj must be struct")
				return false
			},
		},
		{
			name: "City.Town.Int",
			args: args{
				obj:  ct,
				name: "",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.EqualError(t, err, "field path must not be empty")
				return false
			},
		},
		{
			name: ".Town.Int",
			args: args{
				obj:  ct,
				name: ".Town.Int",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.EqualError(t, err, "field path:.Town.Int is invalid")
				return false
			},
		},
		{
			name: "City.Town.Int",
			args: args{
				obj:  ct,
				name: "City.Town.Int",
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "City.Town.Int",
			args: args{
				obj:  &ct,
				name: "City.Town.Int",
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "PtrCity.PtrTown.Int",
			args: args{
				obj:  &ct,
				name: "PtrCity.PtrTown.Int",
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "City.PtrTown.Int",
			args: args{
				obj:  &ct,
				name: "City.PtrTown.Int",
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "City.Town.Bool",
			args: args{
				obj:  &ct,
				name: "City.Town.Bool",
			},
			want:    true,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EmbedStructFieldKind(tt.args.obj, tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("EmbedStructFieldKind(%v, %v)", tt.args.obj, tt.args.name)) {
				return
			}
			assert.Equalf(t, reflect.TypeOf(tt.want).Kind(), got, "EmbedStructFieldKind(%v, %v)", tt.args.obj, tt.args.name)
		})

		t.Run(tt.name, func(t *testing.T) {
			got, err := EmbedStructFieldType(tt.args.obj, tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("EmbedStructFieldType(%v, %v)", tt.args.obj, tt.args.name)) {
				return
			}
			assert.Equalf(t, reflect.TypeOf(tt.want), got, "EmbedStructFieldType(%v, %v)", tt.args.obj, tt.args.name)
		})

		t.Run(tt.name, func(t *testing.T) {
			got, err := EmbedStructFieldTypeStr(tt.args.obj, tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("EmbedStructFieldTypeStr(%v, %v)", tt.args.obj, tt.args.name)) {
				return
			}
			assert.Equalf(t, reflect.TypeOf(tt.want).String(), got, "EmbedStructFieldTypeStr(%v, %v)", tt.args.obj, tt.args.name)

		})
	}
}
