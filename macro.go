package gasp

import "fmt"

type Macro struct {
	obj Object
}

func NewMacro(obj Object) Macro {
	return Macro{obj}
}

// func (m Macro) Eval(env *Env) Object {
// 	ret := Eval(m.obj, env)
// 	ret = Eval(ret, env)
// 	return ret
// }

func (m Macro) String() string {
	return fmt.Sprintf("(macro %s)", m.obj.String())
}
