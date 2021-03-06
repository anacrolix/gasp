package gasp

import (
	"testing"

	_ "github.com/anacrolix/envpprof"
	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	env := NewStandardEnv()
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
		{`3`, `(-> '(1 ((2 3) 4)) second first second)`},
		{`3`, `(infix 1 + 2)`},
		{`6`, `(let [a 2]
				  (let [b 3]
				  	  (* a b)))`},
		{`6`, `(let [a (+ 1 1) b 3]
				   (* a b))`},
		{`(+ (* a b) - 4)`, `(concat '(+ (* a b)) '(- 4))`},
		{`5`, `(let [a (+ 1 1) b 3]
				   (-> (* a b) (+ 4) (- 5)))`},
		{`6`, `(-> 1 (+ 2) (+ 3))`},
		{`8`, `(let [c (+ 1 2)
		           		d 5
		           		e 6]
		            (-> (+ d e) (- c))
		            )`},
		{`6`, `((partial * 2) 3)`},
		{`5`, `((comp (fn [x] (- x 4)) (partial * 3) (partial + 2)) 1)`},
		{`true`, `(<> 1 2)`},
		{`false`, `(<> 1 1)`},
		{`false`, `(any ())`},
		{`false`, `(any '(0))`},
		{`true`, `(any '(1))`},
		{`true`, `(any '(1 0))`},
		{`true`, `(any '(0 1))`},
		{`false`, `(any '(0 0))`},
		{`()`, `(zip '(0) ())`},
		{`(2 4 6)`, `(map (partial * 2) '(1 2 3))`},
		{`(-1 -2 -3)`, `(map - '(1 2 3) '(2 4 6))`},
		{`((1 3) (2 4))`, `(zip '(1 2) '(3 4))`},
	} {
		// log.Println("run", _case.Input)
		assert.EqualValues(t, _case.Input, EvalString(env, _case.Output).String())
	}
}
