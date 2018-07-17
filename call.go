package gasp

import "fmt"

type Caller interface {
	Call(List) Object
}

func Call(obj Object, args List) (ret Object) {
	c, ok := obj.(Caller)
	if !ok {
		panic(fmt.Sprintf("not callable: %s", obj))
	}
	ret = c.Call(args)
	// log.Println("call", obj, "with", args, "returns", ret)
	return
}
