package gasp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringString(t *testing.T) {
	s := NewString(`"`)
	assert.Equal(t, `"\""`, s.String())
}
