package gasp

import (
	"fmt"
	"io"
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
	{NewSymbol("*"), multiply},
}

type goBuiltin struct {
	Symbol string
	Object interface{}
}

var goBuiltins = []goBuiltin{
	{"fmt.Printf", fmt.Printf},
}

func addGoBuiltins() {
	for _, gb := range goBuiltins {
		builtins = append(builtins, builtin{
			NewSymbol("go:" + gb.Symbol),
			WrapGo(gb.Object),
		})
	}
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

func addSpecialBuiltins() {
	addBuiltinFunc("nil?", func(l List) Object {
		if l.First() == nil {
			return True
		} else {
			return False
		}
	})
	addBuiltinFunc("cons", func(l List) Object {
		return l.Rest().First().(List).Cons(l.First())
	})
	addBuiltinFunc("str", func(l List) Object {
		var s string
		for !l.Empty() {
			switch v := l.First().(type) {
			case String:
				s += v.Value
			default:
				s += v.String()
			}
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
	addBuiltinFunc("list?", func(l List) Object {
		_, ok := l.First().(List)
		return NewBool(ok)
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
		r := NewReader(os.Stdin)
		o, err := r.Read()

		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			if o == nil {
				o = NewKeyword("eof")
			}
		}
		return o
	})
	addBuiltinCmpFunc("<", func(cmp int) bool { return cmp < 0 })
	addBuiltinFunc("apply", func(l List) Object {
		f := l.First()
		l = reverse(l.Rest())
		args := l.First().(List)
		l = l.Rest()
		for !l.Empty() {
			args = args.Cons(l.First())
			l = l.Rest()
		}
		return Call(f, args)
	})
}

func init() {
	addSpecialBuiltins()
	addGoBuiltins()
}

func isEmpty(obj Object) Object {
	if obj.(List).Empty() {
		return True
	}
	return False
}
