package gasp

import (
	"fmt"
	"log"
)

type Caller interface {
	Call(List) Object
}

func Call(obj Object, args List) (ret Object) {
	log.Println("call", obj, args)
	c, ok := obj.(Caller)
	if !ok {
		panic(fmt.Sprintf("not callable: %s", obj))
	}
	ret = c.Call(args)
	return
}
