package gasp

import (
	"fmt"
	"math/big"
	"strings"
)

type Object interface {
	String() string
}

type String struct {
	Token Token
	Value string
}

func (s String) String() string {
	return fmt.Sprintf("%q", s.Value)
}

type List struct {
	Next  *List
	Value Object
}

var (
	_ Evaler = List{}

	EmptyList = List{}
)

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

func reverse(l List) List {
	ret := EmptyList
	for !l.Empty() {
		ret = ret.Cons(l.First())
		l = l.Rest()
	}
	return ret
}

type Symbol struct {
	Token Token
	Value string
}

func (s Symbol) Eval(env Env) Object {
	ret := env.NS.Get(s)
	if ret == nil {
		panic(fmt.Sprintf("symbol not found: %s", s.Value))
	}
	return ret
}

func (s Symbol) String() string {
	return fmt.Sprintf("%s", s.Value)
}

type Int struct {
	Token Token
	Value *big.Int
}

func (i Int) String() string {
	return i.Value.String()
}

func NewSymbol(s string) Symbol {
	return Symbol{
		Value: s,
	}
}
