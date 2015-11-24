package gasp

import (
	"fmt"
	"strings"
)

var EmptyList = List{}

type List struct {
	Next  *List
	Value Object
}

func (l List) String() string {
	var strs []string
	for {
		if l.Value == nil {
			break
		}
		strs = append(strs, l.Value.String())
		l = *l.Next
	}
	return fmt.Sprintf("(%s)", strings.Join(strs, " "))
}

func (l List) Empty() bool {
	return l.Value == nil
}

func (l List) Cons(obj Object) List {
	return List{
		Value: obj,
		Next:  &l,
	}
}

func (l List) First() Object {
	return l.Value
}

func (l List) Rest() List {
	return *l.Next
}

func (l List) Eval(env Env) Object {
	el := EmptyList
	for !l.Empty() {
		el = el.Cons(Eval(l.First(), env))
		l = l.Rest()
	}
	el = reverse(el)
	if el.Empty() {
		return EmptyList
	}
	return el.First().(Caller).Call(el.Rest())
}

func (l List) AsSlice() (ret []Object) {
	for !l.Empty() {
		ret = append(ret, l.First())
		l = l.Rest()
	}
	return
}

func reverse(l List) List {
	ret := EmptyList
	for !l.Empty() {
		ret = ret.Cons(l.First())
		l = l.Rest()
	}
	return ret
}
