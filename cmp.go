package gasp

type Cmper interface {
	Cmp(Object) int
}

func Cmp(a, b Object) int {
	if ac, ok := a.(Cmper); ok {
		return ac.Cmp(b)
	}
	panic("uncomparable")
}
