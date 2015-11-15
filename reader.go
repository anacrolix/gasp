package gasp

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
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
	return r.readObject()
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
		obj = String{
			Token: r.t,
			Value: r.t.Value,
		}
		r.advance()
	case TokenTypeSymbol:
		obj = Symbol{
			Token: r.t,
			Value: r.t.Value,
		}
		r.advance()
	case TokenTypeInt:
		i := Int{
			Token: r.t,
			Value: new(big.Int),
		}
		i.Value.SetString(r.t.Value, 0)
		obj = i
		r.advance()
	default:
		err = fmt.Errorf("unexpected token type: %d", r.t.Type)
	}
	return
}

func (r *Reader) advance() {
	r.t, r.err = r.tr.Read()
}

func (r *Reader) readList() (l *List, err error) {
	var objs []Object
	for r.t.Type != RParen {
		var obj Object
		obj, err = r.readObject()
		if err != nil {
			err = fmt.Errorf("while reading list: %s", err)
			return
		}
		objs = append(objs, obj)
	}
	r.advance()
	l = &EmptyList
	for i := len(objs) - 1; i >= 0; i-- {
		l = &List{
			Value: objs[i],
			Next:  l,
		}
	}
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
			panic(err)
		}
		ret = append(ret, obj)
	}
	return
}
