package gasp

import "math/big"

type Int struct {
	Token Token
	Value *big.Int
}

func (i Int) String() string {
	return i.Value.String()
}

func (i Int) Cmp(o Object) int {
	return i.Value.Cmp(o.(Int).Value)
}

func (i Int) True() bool {
	return i.Value.Int64() != 0
}
