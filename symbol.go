package gasp

import "fmt"

type Symbol struct {
	Token Token
	Value string
}

var (
	_ Evaler = Symbol{}
)

func NewSymbol(s string) Symbol {
	return Symbol{
		Value: s,
	}
}

func (s Symbol) Eval(env *Env) Object {
	ret := env.Lookup(s)
	if ret == nil {
		panic(fmt.Sprintf("symbol %q not found in %s", s.Value, env.String()))
	}
	return ret
}

func (s Symbol) String() string {
	return fmt.Sprintf("%s", s.Value)
}

func (s Symbol) Name() string {
	return s.Value
}
