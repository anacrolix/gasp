package gasp

type Evaler interface {
	Eval(*Env) Object
}

func Eval(obj Object, env *Env) (ret Object) {
	if e, ok := obj.(Evaler); ok {
		// defer func() { log.Println("eval", obj, "->", ret) }()
		ret = e.Eval(env)
		return
	}
	ret = obj
	return
}

func EvalString(env *Env, s string) (ret Object) {
	objs := ReadString(s)
	// log.Println(objs)
	for _, o := range objs {
		ret = Eval(o, env)
	}
	return
}
