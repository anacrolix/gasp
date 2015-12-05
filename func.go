package gasp

import "github.com/bradfitz/iter"

type Func struct {
	args  List
	body  List
	outer *Env
}

var (
	_ Caller = Func{}
	// _ Evaler = Func{}
)

func NewFunc(args List, body List, outer *Env) Func {
	return Func{args, body, outer}
}

func (f Func) Call(args List) (ret Object) {
	env := Env{
		Outer: f.outer,
		NS:    NewMap(),
	}
	if args.Len() != f.args.Len() {
		panic("argument count mismatch")
	}
	for i := range iter.N(f.args.Len()) {
		env.Define(f.args.Nth(i), args.Nth(i))
	}
	body := f.body
	for !body.Empty() {
		ret = Eval(body.First(), &env)
		body = body.Rest()
	}
	return
}

func (f Func) String() string {
	return f.body.Cons(f.args).Cons(NewSymbol("fn")).String()
}
