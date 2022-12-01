package error_handling

import (
	"fmt"
	"runtime/debug"
)

type NovaError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func (err NovaError) Error() string {
	return err.Message
}

func WrapError(err error, messageF string, msgArgs ...interface{}) NovaError {
	return NovaError{
		Inner:      err, //<1>
		Message:    fmt.Sprintf(messageF, msgArgs...),
		StackTrace: string(debug.Stack()),        // <2>
		Misc:       make(map[string]interface{}), // <3>
	}
}
