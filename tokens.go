package gasp

type TokenReader interface {
	Read() (Token, error)
}

type Token struct {
	Type    TokenType
	Value   string
	Line    int64
	LineOff int64
}

type TokenType byte

const (
	LParen = iota + 1
	RParen
	Str
	Whitespace
	TokenTypeSymbol
	TokenTypeInt
)

func (t *Token) String() string {
	switch t.Type {
	case LParen:
		return "("
	case RParen:
		return ")"
	default:
		return "?"
	}
}
