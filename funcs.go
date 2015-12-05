package gasp

import (
	"fmt"
	"math/big"
)

type Caller interface {
	Call(List) Object
}

type builtinCallable struct {
	f    func(l List) Object
	name string
}

var (
	_ Caller = builtinCallable{}
)

func (me builtinCallable) String() string {
	return fmt.Sprintf("#(%s)", me.name)
}

func (me builtinCallable) Call(args List) Object {
	return me.f(args)
}

func (me builtinCallable) True() bool {
	return me.f != nil
}

var add = builtinCallable{
	f: func(l List) Object {
		var ret big.Int
		for !l.Empty() {
			ret.Add(&ret, l.First().(Int).Value)
			l = l.Rest()
		}
		return Int{
			Value: &ret,
		}
	},
	name: "add",
}

var multiply = builtinCallable{
	f: func(l List) Object {
		var ret big.Int
		ret.Set(l.First().(Int).Value)
		l = l.Rest()
		for !l.Empty() {
			ret.Mul(&ret, l.First().(Int).Value)
			l = l.Rest()
		}
		return Int{
			Value: &ret,
		}
	},
	name: "multiply",
}

var subtract = builtinCallable{
	f: func(l List) Object {
		var ret big.Int
		ret.Set(l.First().(Int).Value)
		l = l.Rest()
		for !l.Empty() {
			ret.Sub(&ret, l.First().(Int).Value)
			l = l.Rest()
		}
		return Int{
			Value: &ret,
		}
	},
	name: "subtract",
}
