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
		Type: LParen,
		Line: 1,
	}}, io.EOF},
	{" (", []Token{Token{
		Type:    LParen,
		Line:    1,
		LineOff: 1,
	}}, io.EOF},
	{"\n(", []Token{Token{
		Type: LParen,
		Line: 2,
	}}, io.EOF},
	{
		"\n(\n)",
		[]Token{
			Token{
				Type: LParen,
				Line: 2,
			},
			Token{
				Type: RParen,
				Line: 3,
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
		assert.EqualValues(t, _case.Tokens, ts)
	}
}
