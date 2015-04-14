package tracederror

import (
	"fmt"
	"runtime"
)

var (
	// If debug is false, no stack or file/line information is included with the wrapped error.
	Debug     = true

	// If this is false, then only the file and line where the error was created is printed, instead 
	// of a full stack trace.
	ShowStack = true
)

type TracedError struct {
	WrappedError error
	stack        []byte
	line         int
	file         string
}

// Creates a new traced error. Calling this on an instance of traced error is idempotent, 
// it just returns the original traced error. Calling this on nil returns nil.
func New(wrappedError error) error {
	if wrappedError == nil {
		return nil
	}

	_, ok := wrappedError.(*TracedError)
	if ok {
		return wrappedError
	}

	enErr := &TracedError{
		WrappedError: wrappedError,
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
