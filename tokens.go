package gasp

import "fmt"

type TokenReader interface {
	Read() (Token, error)
}

type Token struct {
	Type    TokenType
	Value   string
	Line    int64 // 1-indexed
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
	TokenTypeComment
)

func (t Token) String() string {
	return fmt.Sprintf("%q:%d:%d", t.Value, t.Line, t.LineOff)
}
