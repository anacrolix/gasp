package gasp

import (
	"fmt"
	"log"
	"reflect"
)

type GoObject struct {
	v reflect.Value
}

func (me GoObject) String() string {
	return fmt.Sprint(me)
}

type Goer interface {
	ToGo(reflect.Type) interface{}
}

func WrapGo(i interface{}) Object {
	return GoObject{reflect.ValueOf(i)}
}

func ToGo(obj Object, typ reflect.Type) interface{} {
	return obj.(Goer).ToGo(typ)
}

func FromGo(i interface{}) Object {
	switch v := i.(type) {
	case string:
		return NewString(v)
	}
	panic(i)
}

func inType(call reflect.Type, nthArg int) reflect.Type {
	if call.IsVariadic() && nthArg >= call.NumIn()-1 {
		return call.In(call.NumIn() - 1).Elem()
	}
	return call.In(nthArg)
}

func (me GoObject) Call(args List) Object {
	if me.v.Kind() != reflect.Func {
		panic(fmt.Sprintf("not callable: %s", me))
	}
	var in []reflect.Value
	var i int
	for !args.Empty() {
		in = append(in, reflect.ValueOf(ToGo(args.First(), inType(me.v.Type(), i))))
		args = args.Rest()
		i++
	}
	log.Println(len(in), in)
	out := me.v.Call(in)
	log.Println(len(out), out)
	if len(out) == 1 {
		return FromGo(out[0].Interface())
	}
	ret := EmptyList
	for _, o := range out {
		ret = ret.Cons(FromGo(o))
	}
	return ret
}
