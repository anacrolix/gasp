package gasp

import (
	"fmt"
	"log"
	"strings"
)

var EmptyList = List{}

type List struct {
	Next  *List
	Value Object
}

var (
	_ Evaler = List{}
	_ Truer  = List{}
)

func (l List) True() bool {
	return !l.Empty()
}

func (l List) Len() (ret int) {
	for !l.Empty() {
		ret++
		l = l.Rest()
	}
	return
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
	return l == EmptyList
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

func (l List) Nth(i int) Object {
	if i < 0 {
		return nil
	}
	if i == 0 {
		return l.First()
	}
	return l.Rest().Nth(i - 1)
}

func (l List) Eval(env *Env) Object {
	if s, ok := l.First().(Symbol); ok {
		switch s.Name() {
		case "import":
			return nil
		case "if":
			res := Eval(l.Nth(1), env)
			if IsTrue(res) {
				return Eval(l.Nth(2), env)
			}
			if l.Len() < 4 {
				return Nil
			}
			return Eval(l.Nth(3), env)
		case "def":
			l = l.Rest()
			name := l.First()
			l = l.Rest()
			val := Eval(l.First(), env)
			if !l.Rest().Empty() {
				panic(fmt.Sprintf("extraneous arguments to def: %s", l.Rest()))
			}
			env.Define(name, val)
			return nil
		case "fn":
			return NewFunc(l.Nth(1).(List), l.Drop(2), env)
		}
	}
	el := EmptyList
	for !l.Empty() {
		el = el.Cons(Eval(l.First(), env))
		l = l.Rest()
	}
	el = reverse(el)
	log.Println(el)
	if el.Empty() {
		return EmptyList
	}
	first, ok := el.First().(Caller)
	if !ok {
		panic(fmt.Sprintf("not callable: %s in %s", el.First(), l))
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

func (l List) Drop(i int) List {
	for ; i > 0; i-- {
		l = l.Rest()
	}
	return l
}

func reverse(l List) List {
	ret := EmptyList
	for !l.Empty() {
		ret = ret.Cons(l.First())
		l = l.Rest()
	}
	return ret
}

func ListFromSlice(sl []Object) (ret List) {
	for i := len(sl) - 1; i >= 0; i-- {
		ret = ret.Cons(sl[i])
	}
	return
}
