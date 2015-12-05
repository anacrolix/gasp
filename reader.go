package gasp

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/anacrolix/missinggo"
)

type Reader struct {
	tr  *Tokenizer
	t   Token
	err error
}

func NewReader(r io.Reader) *Reader {
	ret := &Reader{
		tr: NewTokenizer(r),
	}
	ret.advance()
	return ret
}

func (r *Reader) Read() (obj Object, err error) {
	obj, err = r.readObject()
	if err == io.EOF {
		return
	}
	if err != nil {
		err = fmt.Errorf("error reading form at %d:%d: %s", r.tr.line, r.tr.lineOff, err)
	}
	return
}

func unescapeStr(s string) (ret string, err error) {
	escaped := false
	for _, c := range s {
		if escaped {
			switch c {
			case '"', '\\':
			case 'n':
				c = '\n'
			default:
				err = fmt.Errorf("invalid escape: \\%c", c)
				return
			}
			ret += string(c)
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		ret += string(c)
	}
	if escaped {
		err = fmt.Errorf("unexpected end of string: %q", ret)
	}
	return
}

func (r *Reader) readObject() (obj Object, err error) {
	if r.err != nil {
		err = r.err
		return
	}
	switch r.t.Type {
	case LParen:
		r.advance()
		obj, err = r.readList()
	case Str:
		s := String{
			Token: r.t,
		}
		s.Value, err = unescapeStr(r.t.Value)
		if err != nil {
			return
		}
		obj = s
		r.advance()
	case TokenTypeSymbol:
		s := Symbol{
			Token: r.t,
			Value: r.t.Value,
		}
		r.advance()
		switch s.Value {
		case "'":
			obj, err = r.readObject()
			obj = EmptyList.Cons(obj).Cons(NewSymbol("quote"))
		default:
			obj = s
		}
	case TokenTypeInt:
		i := Int{
			Token: r.t,
			Value: new(big.Int),
		}
		i.Value.SetString(r.t.Value, 0)
		obj = i
		r.advance()
	case TokenTypeComment:
		r.advance()
	default:
		err = fmt.Errorf("unexpected token: %s", r.t)
	}
	return
}

func (r *Reader) advance() {
	r.t, r.err = r.tr.Read()
}

func sprintsep(a ...interface{}) string {
	s := fmt.Sprintln(a...)
	return s[:len(s)-1]
}

func (r *Reader) readList() (ret List, err error) {
	var objs []Object
	for r.t.Type != RParen {
		var obj Object
		obj, err = r.readObject()
		if err != nil {
			err = fmt.Errorf("while reading list (%s: %s", sprintsep(missinggo.ConvertToSliceOfEmptyInterface(objs)...), err)
			return
		}
		objs = append(objs, obj)
	}
	r.advance()
	l := &EmptyList
	for i := len(objs) - 1; i >= 0; i-- {
		l = &List{
			Value: objs[i],
			Next:  l,
		}
	}
	ret = *l
	return
}

func ReadString(s string) (ret []Object) {
	r := NewReader(bytes.NewReader([]byte(s)))
	for {
		obj, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Errorf("error reading string %s", err))
		}
		ret = append(ret, obj)
	}
	return
}
