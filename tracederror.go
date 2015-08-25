package tracederror

import (
	"fmt"
	"runtime"
)

var (
	// If debug is false, no stack or file/line information is included with the wrapped error.
	Debug = true

	// If this is false, then only the file and line where the error was created is printed, instead
	// of a full stack trace.
	ShowStack = true
)

type TracedError struct {
	WrappedError error
	Context      interface{}
	stack        []byte
	line         int
	file         string
}

var onError func(error, []byte, string, int, interface{}) = nil

func OnError(handler func(error, []byte, string, int, interface{})) {
	onError = handler
}

// Creates a new traced error. Calling this on an instance of traced error is idempotent,
// it just returns the original traced error. Calling this on nil returns nil.
func New(wrappedError error) error {
	return NewWithContext(wrappedError, nil)
}

func NewWithContext(wrappedError error, context interface{}) error {
	if wrappedError == nil {
		return nil
	}

	_, ok := wrappedError.(*TracedError)
	if ok {
		return wrappedError
	}

	enErr := &TracedError{
		WrappedError: wrappedError,
		Context:      context,
	}

	if Debug {
		if ShowStack {
			if _, ok := wrappedError.(*TracedError); !ok {
				enErr.stack = stack()
			}
		}

		if _, file, line, ok := runtime.Caller(1); ok {
			enErr.file = file
			enErr.line = line
		}

		if onError != nil {
			onError(wrappedError, enErr.stack, enErr.file, enErr.line, context)
		}
	}

	return enErr
}

// Returns the original error's Error(), with optional caller information
func (enErr *TracedError) Error() string {
	if Debug {
		if ShowStack {
			return fmt.Sprintf("Error: %s\nStack: %s", enErr.WrappedError.Error(), string(enErr.stack))
		} else {
			return fmt.Sprintf("Error: %s\nFile: %s\nLine: %d", enErr.WrappedError.Error(), enErr.file, enErr.line)
		}
	} else {
		return enErr.WrappedError.Error()
	}
}

// Returns the wrapped error, or nil if it doesn't exist.
// If it is not a traced error, just return the error by itself.
func Inner(e error) error {
	if e == nil {
		return nil
	}

	self := e.(*TracedError)
	if self != nil {
		return self.WrappedError
	}

	return e
}

func stack() []byte {
	buf := make([]byte, 32)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			break
		}
		buf = make([]byte, len(buf)*2)
	}
	return buf
}
