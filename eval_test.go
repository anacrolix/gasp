package gasp

import (
	"log"
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
		{`(3 2 1)`, `(reduce (flip cons) () '(1 2 3))`},
		{`(1 1 9 9 8 8 8)`, `(list 1 1 9 9 8 8 8)`},
		{`(1 2 3 4)`, `(concat '(1 2) '(3 4))`},
		{`(1 2)`, `'(1 2)`},
		{`10`, `(apply + '(1 2 3 4))`},
		{`(1 2 3 4)`, `(conj '(1 2) 3 4)`},
		{`false`, `(and (< 1 2) (> 3 4))`},
	} {
		log.Println("run", _case.Input)
		assert.EqualValues(t, _case.Input, EvalString(env, _case.Output).String())
	}
}
