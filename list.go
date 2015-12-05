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
	if l.Empty() {
		return l
	}
	first := l.First()
	if s, ok := first.(Symbol); ok {
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
		case "quote":
			return l.Rest().First()
		case "macro":
			return NewMacro(Eval(l.Rest().First(), env))
		case "eval":
			// This is here because I don't propagate env to calls.
			return Eval(Eval(l.Rest().First(), env), env)
		}
	}
	first = Eval(first, env)
	l = l.Rest()
	if m, ok := first.(Macro); ok {
		log.Println("macro", m.obj)
		return Eval(Call(m.obj, l), env)
	}
	var args List
	for !l.Empty() {
		args = args.Cons(Eval(l.First(), env))
		l = l.Rest()
	}
	args = reverse(args)
	return first.(Caller).Call(args)
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
