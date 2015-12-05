package gasp

var (
	Nil = _nil{}

	_ Truer  = Nil
	_ Object = Nil
)

type _nil struct{}

func (_nil) True() bool {
	return false
}

func (_nil) String() string {
	return "nil"
}
