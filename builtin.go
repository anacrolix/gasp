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
