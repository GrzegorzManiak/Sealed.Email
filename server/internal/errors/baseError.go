package errors

import "fmt"

type BaseError struct {
	message string
	title   string
	fault   string
	code    int
}

func (e BaseError) Fault() string {
	return e.fault
}

func (e BaseError) Message() string {
	return e.message
}

func (e BaseError) Title() string {
	return e.title
}

func (e BaseError) Code() int {
	return e.code
}

func (e BaseError) Error() string {
	return fmt.Sprintf("[%s]: %s (%s)", e.title, e.message, e.fault)
}
