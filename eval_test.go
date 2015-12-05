package gasp

import (
	"testing"

	_ "github.com/anacrolix/envpprof"
	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	env := NewStandardEnv()
	env.NS = env.NS.Assoc(NewSymbol("+"), add).Assoc(NewSymbol("*"), multiply).Assoc(NewSymbol("-"), subtract)
	assert.EqualValues(t, `6`, EvalString(env, `(+ 1 (+ 2 3))`).String())
	assert.EqualValues(t, `6`, EvalString(env, `(* (+ 1 2) (- 4 2))`).String())
	assert.EqualValues(t, `()`, EvalString(env, `()`).String())
	assert.EqualValues(t, `42`, EvalString(env, `(if (> 2 1) 42)`).String())
	assert.EqualValues(t, `nil`, EvalString(env, `(if (<= 2 1) 42)`).String())
	assert.EqualValues(t, `false`, EvalString(env, `(not (< 1 2))`).String())
	assert.EqualValues(t, `true`, EvalString(env, `(not (> 1 2))`).String())
}
