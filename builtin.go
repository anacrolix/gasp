package gasp

import (
	"fmt"
	"io/ioutil"
	"os"
)

type builtin struct {
	Symbol
	Object
}

var builtins = []builtin{
	{NewSymbol("true"), True},
	{NewSymbol("false"), False},
	{NewSymbol("+"), add},
	{NewSymbol("-"), subtract},
}

func addBuiltinFunc(name string, f func(List) Object) {
	builtins = append(builtins, builtin{NewSymbol(name), builtinCallable{f, name}})
}

func addBuiltinCmpFunc(sym string, pred func(int) bool) {
	addBuiltinFunc(sym, func(l List) Object {
		if pred(Cmp(l.First(), l.Rest().First())) {
			return True
		} else {
			return False
		}
	})
}

func init() {
	addBuiltinFunc("cons", func(l List) Object {
		return l.Rest().First().(List).Cons(l.First())
	})
	addBuiltinFunc("str", func(l List) Object {
		var s string
		for !l.Empty() {
			s += l.First().(String).Value
			l = l.Rest()
		}
		return NewString(s)
	})
	addBuiltinFunc("concat", func(l List) Object {
		ret := EmptyList
		for !l.Empty() {
			r := l.First().(List)
			for !r.Empty() {
				ret = ret.Cons(r.First())
				r = r.Rest()
			}
			l = l.Rest()
		}
		return reverse(ret)
	})
	addBuiltinFunc("empty?", func(l List) Object {
		return isEmpty(l.First())
	})
	addBuiltinFunc("first", func(l List) Object {
		return l.First().(List).First()
	})
	addBuiltinFunc("rest", func(l List) Object {
		return l.First().(List).Rest()
	})
	addBuiltinFunc("print", func(l List) Object {
		if l.Len() != 1 {
			panic(fmt.Sprintf("print expected exactly one argument, got %d", l.Len()))
		}
		fmt.Fprintln(os.Stdout, l.First())
		return nil
	})
	addBuiltinFunc("read", func(l List) Object {
		if !l.Empty() {
			panic("read doesn't take arguments")
		}
		s, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		objs := ReadString(string(s))
		return ListFromSlice(objs)
	})
	// addBuiltinCmpFunc(">", func(cmp int) bool { return cmp > 0 })
	addBuiltinCmpFunc("<", func(cmp int) bool { return cmp < 0 })
}

func isEmpty(obj Object) Object {
	if obj.(List).Empty() {
		return True
	}
	return False
}
