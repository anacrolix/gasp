package gasp

import (
	"bytes"
	"io"
	"testing"

	_ "github.com/anacrolix/envpprof"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
cases:
	for _, _case := range []struct {
		Input    string
		Expected []string
		Err      error
	}{
		{
			"",
			nil,
			io.EOF,
		}, {
			`"hello"`,
			[]string{`"hello"`},
			io.EOF,
		}, {
			`()`,
			[]string{`()`},
			io.EOF,
		},
		{
			`(a b) (b c)`,
			[]string{`(a b)`, `(b c)`},
			io.EOF,
		},
		{
			`(* (+ 1 2) (- 4 2))`,
			[]string{`(* (+ 1 2) (- 4 2))`},
			io.EOF,
		},
		{
			`(Sprintf "%s, %s" "hello" "world")`,
			[]string{`(Sprintf "%s, %s" "hello" "world")`},
			io.EOF,
		},
	} {
		r := NewReader(bytes.NewReader([]byte(_case.Input)))
		var err error
		for _, expObj := range _case.Expected {
			var obj Object
			obj, err = r.Read()
			if err != nil {
				t.Errorf("expected %q: %s", expObj, err)
				continue cases
			}
			assert.EqualValues(t, expObj, obj.String())
		}
		obj, err := r.Read()
		assert.Nil(t, obj)
		assert.EqualValues(t, _case.Err, err)
	}
}
