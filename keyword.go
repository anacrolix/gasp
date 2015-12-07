package gasp

type Keyword struct {
	S string
}

func NewKeyword(s string) Keyword {
	return Keyword{s}
}

func (k Keyword) String() string {
	return ":" + k.S
}
