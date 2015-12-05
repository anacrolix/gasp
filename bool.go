package gasp

var (
	True  = _bool{true}
	False = _bool{false}

	_ Truer = _bool{}
)

type _bool struct {
	truth bool
}

func (me _bool) String() string {
	if me.truth {
		return "#t"
	} else {
		return "#f"
	}
}

func (me _bool) True() bool {
	return me.truth
}
