package gasp

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cases = []struct {
	Input  string
	Tokens []Token
	Err    error
}{
	{"", nil, io.EOF},
	{"(", []Token{Token{
		Type:  LParen,
		Line:  1,
		Value: "(",
	}}, io.EOF},
	{" (", []Token{Token{
		Type:    LParen,
		Line:    1,
		LineOff: 1,
		Value:   "(",
	}}, io.EOF},
	{"\n(", []Token{Token{
		Type:  LParen,
		Line:  2,
		Value: "(",
	}}, io.EOF},
	{
		"\n(\n)",
		[]Token{
			Token{
				Type:  LParen,
				Line:  2,
				Value: "(",
			},
			Token{
				Type:  RParen,
				Line:  3,
				Value: ")",
			},
		},
		io.EOF,
	},
	{
		`"\""`,
		[]Token{
			Token{
				Type:  Str,
				Line:  1,
				Value: "\\\"",
			},
		},
		io.EOF,
	},
	{
		`"hello, world"`,
		[]Token{
			Token{
				Type:  Str,
				Line:  1,
				Value: "hello, world",
			},
		},
		io.EOF,
	},
	{
		`(Sprintf "%s, %s" "hello" "world")`,
		[]Token{
			Token{
				Type:    LParen,
				Line:    1,
				LineOff: 0,
				Value:   "(",
			},
			Token{
				Type:    TokenTypeSymbol,
				Line:    1,
				LineOff: 1,
				Value:   "Sprintf",
			},
			Token{
				Type:    Str,
				Line:    1,
				LineOff: 9,
				Value:   "%s, %s",
			},
			Token{
				Type:    Str,
				Line:    1,
				LineOff: 18,
				Value:   "hello",
			},
			Token{
				Type:    Str,
				Line:    1,
				LineOff: 26,
				Value:   "world",
			},
			Token{
				Type:    RParen,
				Line:    1,
				LineOff: 33,
				Value:   ")",
			},
		},
		io.EOF,
	},
}

func readAllTokens(tr TokenReader) (ts []Token, err error) {
	for {
		var t Token
		t, err = tr.Read()
		if err != nil {
			return
		}
		ts = append(ts, t)
	}
}

func TestTokenizer(t *testing.T) {
	for _, _case := range cases {
		tr := NewTokenizer(bytes.NewReader([]byte(_case.Input)))
		ts, err := readAllTokens(tr)
		assert.EqualValues(t, _case.Err, err, "%q", _case.Input)
		assert.EqualValues(t, _case.Tokens, ts, "%q", _case.Input)
	}
}
