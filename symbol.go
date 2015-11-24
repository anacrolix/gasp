package gasp

import "fmt"

type Symbol struct {
	Token Token
	Value string
}

func NewSymbol(s string) Symbol {
	return Symbol{
		Value: s,
	}
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

func (s Symbol) Name() string {
	return s.Value
}
