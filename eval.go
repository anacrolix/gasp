package gasp

import "log"

type Evaler interface {
	Eval(*Env) Object
}

func Eval(obj Object, env *Env) Object {
	if e, ok := obj.(Evaler); ok {
		return e.Eval(env)
	}
	return obj
}

func EvalString(env *Env, s string) (ret Object) {
	objs := ReadString(s)
	log.Println(objs)
	for _, o := range objs {
		ret = Eval(o, env)
	}
	return
}
