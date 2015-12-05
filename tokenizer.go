package gasp

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

	"github.com/bradfitz/iter"
)

type Tokenizer struct {
	r       *bufio.Reader
	line    int64
	lineOff int64
	buf     []byte
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		r:    bufio.NewReader(r),
		line: 1,
	}
}

func (me *Tokenizer) Read() (t Token, err error) {
	for {
		err = me.buffer()
		tt, ms, ok := me.match()
		if ok && tt == Whitespace {
			me.advance(len(ms[0]))
			continue
		}
		if ok && (len(ms[0]) < len(me.buf) || err == io.EOF) {
			t.Type = tt
			switch len(ms) {
			case 1:
			case 2:
				t.Value = string(ms[1])
			default:
				panic(fmt.Sprintf("%d groups captured: %+q", len(ms), ms))
			}
			t.Line = me.line
			t.LineOff = me.lineOff
			me.advance(len(ms[0]))
			err = nil
			return
		}
		if err == io.EOF && len(me.buf) != 0 {
			err = fmt.Errorf("can't tokenize %q", me.buf)
		}
		if err != nil {
			return
		}
	}
}

func (me *Tokenizer) advance(n int) {
	for i := range iter.N(n) {
		if me.buf[i] == '\n' {
			me.line++
			me.lineOff = 0
		} else {
			me.lineOff++
		}
	}
	me.buf = me.buf[n:]
}

func (me *Tokenizer) buffer() (err error) {
	b := make([]byte, 4096)
	n, err := me.r.Read(b)
	me.buf = append(me.buf, b[:n]...)
	return
}

var tokens = []struct {
	Type   TokenType
	Regexp *regexp.Regexp
}{
	{
		Type:   Str,
		Regexp: regexp.MustCompile(`^"((?:|(?:[^"]|\\")*[^\\]))"`),
	},
	{
		Type:   LParen,
		Regexp: regexp.MustCompile(`^(\(|\[)`),
	},
	{
		Type:   RParen,
		Regexp: regexp.MustCompile(`^(\)|\])`),
	},
	{
		Type:   Whitespace,
		Regexp: regexp.MustCompile(`^\s+`),
	},
	{
		Type:   TokenTypeSymbol,
		Regexp: regexp.MustCompile(`^([a-zA-Z*+-/.:<>=&?']+)`),
	},
	{
		Type:   TokenTypeInt,
		Regexp: regexp.MustCompile(`^([0-9]+)`),
	},
}

func (me *Tokenizer) match() (tt TokenType, ms [][]byte, ok bool) {
	for _, t := range tokens {
		ms1 := t.Regexp.FindSubmatch(me.buf)
		if ms1 == nil {
			continue
		}
		if ok {
			panic(fmt.Sprintf("ambiguous token: %s %q and %s %q", tt, ms[0], t.Type, ms1[0]))
		}
		ms = ms1
		tt = t.Type
		ok = true
	}
	return
}
