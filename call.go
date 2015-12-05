package gasp

import (
	"fmt"
	"log"
)

type Caller interface {
	Call(List) Object
}

func Call(obj Object, args List) Object {
	log.Println("call", obj, args)
	c, ok := obj.(Caller)
	if !ok {
		panic(fmt.Sprintf("not callable: %s", obj))
	}
	return c.Call(args)
}
