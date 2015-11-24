package gasp

import (
	"fmt"
	"math/big"
)

type Object interface {
	String() string
}

var (
	_ Evaler = List{}
)

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
