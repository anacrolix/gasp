package gasp

type Evaler interface {
	Eval(Env) Object
}

func Eval(obj Object, env Env) Object {
	if e, ok := obj.(Evaler); ok {
		return e.Eval(env)
	}
	return obj
}

func EvalString(env Env, s string) Object {
	obj := ReadString(s)
	return Eval(obj, env)
}
