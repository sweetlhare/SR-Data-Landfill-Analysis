package custom

import "fmt"

type Error string

// Error ...
func (e Error) Error() string {
	return string(e)
}

// Description ...
func (e Error) AddDescription(s string) Error {
	return Error(fmt.Sprintf("%s (%s)", string(e), s))
}

// AddPrefix ...
func (e Error) AddPrefix(s string) Error {
	return Error(fmt.Sprintf("%s: %s", s, string(e)))
}
