package gasp

import "math/big"

type Object interface {
	String() string
}

var (
	_ Evaler = List{}
)

type Int struct {
	Token Token
	Value *big.Int
}

func (i Int) String() string {
	return i.Value.String()
}
