package gasp

import (
	"fmt"
	"reflect"
)

type String struct {
	Token Token
	Value string
}

func (s String) String() string {
	return fmt.Sprintf("%q", s.Value)
}

func (s String) ToGo(typ reflect.Type) interface{} {
	switch typ.Kind() {
	case reflect.String:
		return s.Value
	case reflect.Interface:
		return s.Value
	default:
		panic(typ.Kind())
	}
}

func NewString(s string) Object {
	return String{Value: s}
}
