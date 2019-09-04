package helpers

import (
	"fmt"
)

// RecoverPanic is a handler to be used in `defer`-s to recover panics and
// store them into `err`.
func RecoverPanic(err *error) {
	recoverResult := recover()
	if recoverResult == nil {
		return
	}

	if err == nil {
		return
	}

	switch recoverResult := recoverResult.(type) {
	case error:
		*err = recoverResult
	default:
		*err = fmt.Errorf("%v", recoverResult)
	}

	// TODO: add a stack trace
}
