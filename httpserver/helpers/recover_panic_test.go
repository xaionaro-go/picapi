package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoverPanic(t *testing.T) {
	fn := func() (err error) {
		defer RecoverPanic(&err)
		panic(`test`)
	}

	assert.Error(t, fn())
}
