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
	if l.Empty() {
		return EmptyList
	}
	if s, ok := l.First().(Symbol); ok {
		switch s.Name() {
		case "import":
			return nil
		}
	}
	el := EmptyList
	for !l.Empty() {
		el = el.Cons(Eval(l.First(), env))
		l = l.Rest()
	}
	el = reverse(el)
	first, ok := el.First().(Caller)
	if !ok {
		panic(fmt.Sprint("not caller: %s in %s", first))
	}
	return first.Call(el.Rest())
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
