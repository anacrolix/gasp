package gasp

import "github.com/bradfitz/iter"

type Func struct {
	params List
	rest   Object
	body   List
	outer  *Env
}

var (
	_ Caller = Func{}
	// _ Evaler = Func{}
)

func parseParams(ps List) (pos List, rest Object) {
	for !ps.Empty() {
		s := ps.First().(Symbol)
		if s.Value == "&" {
			rest = ps.Rest().First()
			break
		}
		pos = pos.Cons(s)
		ps = ps.Rest()
	}
	pos = reverse(pos)
	return
}

func NewFunc(args List, body List, outer *Env) Func {
	f := Func{
		body:  body,
		outer: outer,
	}
	f.params, f.rest = parseParams(args)
	return f
}

func (f Func) Call(args List) (ret Object) {
	env := Env{
		Outer: f.outer,
		NS:    NewMap(),
	}
	for i := range iter.N(f.params.Len()) {
		env.Define(f.params.Nth(i), args.Nth(i))
	}
	if f.rest != nil {
		env.Define(f.rest, args.Drop(f.params.Len()))
	}
	body := f.body
	for !body.Empty() {
		ret = Eval(body.First(), &env)
		body = body.Rest()
	}
	return
}

func (f Func) String() string {
	return f.body.Cons(f.params).Cons(NewSymbol("fn")).String()
}
