Tracederror
================

A simple error wrapper to add caller information.

## Usage ##
Wrap errors upon return to add caller information. Instead of
~~~~
	err := somethingThatCanFail(arguments)
	return err
~~~~

Do
~~~~
	err := somethingThatCanFail(arguments)
	return tracederror.New(err)
~~~~

Since tracederror.New is idempotent, calls to tracederror.New can be inserted anywhere
there is an error without extra wrapping:

~~~~~
func odd(a int) error{
	if a%2 == 0 {
		return tracederror.New(fmt.Errorf("Input was not odd!"))
	}

	return nil
}

func complicated(a int) error {
	err := odd(a)

	// This works even if err is nil or already a traced error.
	// If err is nil, then tracederror.New will also return nil, so no changes need to be made
	// to error handling code.
	return tracederror.New(err)
}
~~~~~