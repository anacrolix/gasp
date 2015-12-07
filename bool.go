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
		return "true"
	} else {
		return "false"
	}
}

func (me _bool) True() bool {
	return me.truth
}

func NewBool(b bool) Object {
	if b {
		return True
	} else {
		return False
	}
}
