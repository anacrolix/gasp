package gasp

import (
	"testing"

	_ "github.com/anacrolix/envpprof"
	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	env := NewStandardEnv()
	env.NS = env.NS.Assoc(NewSymbol("+"), add).Assoc(NewSymbol("*"), multiply).Assoc(NewSymbol("-"), subtract)
	for _, _case := range []struct {
		Input, Output string
	}{
		{`6`, `(+ 1 (+ 2 3))`},
		{`6`, `(+ 1 (+ 2 3))`},
		{`6`, `(* (+ 1 2) (- 4 2))`},
		{`()`, `()`},
		{`42`, `(if (> 2 1) 42)`},
		{`nil`, `(if (<= 2 1) 42)`},
		{`false`, `(not (< 1 2))`},
		{`true`, `(not (> 1 2))`},
		{`13`, `(reduce + 3 '(1 2 3 4))`},
	} {
		assert.EqualValues(t, _case.Input, EvalString(env, _case.Output).String())
	}
}
