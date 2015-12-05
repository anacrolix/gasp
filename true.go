package gasp

type Truer interface {
	True() bool
}

func IsTrue(obj Object) bool {
	if obj == nil {
		return false
	}
	if t, ok := obj.(Truer); ok {
		return t.True()
	}
	return true
}
