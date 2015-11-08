package gasp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	env := Env{
		NS: NewMap().Assoc(NewSymbol("+"), add).Assoc(NewSymbol("*"), multiply).Assoc(NewSymbol("-"), subtract),
	}
	assert.EqualValues(t, `6`, EvalString(env, `(+ 1 (+ 2 3))`).String())
	assert.EqualValues(t, `6`, EvalString(env, `(* (+ 1 2) (- 4 2))`).String())
	assert.EqualValues(t, `()`, EvalString(env, `()`).String())
}
