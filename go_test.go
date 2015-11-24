package gasp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromGo(t *testing.T) {
	obj := FromGo("hello")
	assert.IsType(t, String{}, obj)
	assert.EqualValues(t, `"hello"`, obj.String())
}

func TestGoCall(t *testing.T) {
	env := Env{
		NS: NewMap().Assoc(NewSymbol("Sprintf"), WrapGo(fmt.Sprintf)),
	}
	assert.EqualValues(t, `"hello, world"`, EvalString(env, `(Sprintf "%s, %s" "hello" "world")`).String())
}
