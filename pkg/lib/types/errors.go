package types

import (
	"fmt"
	"strings"
)

type ErrorSource int

const (
	InputError ErrorSource = iota
	OutputError
	SupervisorError
	UnknownErrorSource
)

type Errors []error

func (e *Errors) Add(err interface{}) *Errors {
	if err == nil {
		return e
	}
	switch val := err.(type) {
	case error:
		*e = append(*e, val)
	case Errors:
		*e = append(*e, val...)
	default:
		*e = append(*e, fmt.Errorf("%v", val))
	}
	return e
}
func (e Errors) Count() int {
	return len(e)
}
func (e Errors) IsEmpty() bool {
	return e.Count() == 0
}
func (e Errors) Wrap() error {
	if e.Count() == 0 {
		return nil
	}
	errs := make([]string, e.Count())
	for i, err := range e {
		errs[i] = err.Error()
	}
	return fmt.Errorf("%v", strings.Join(errs, "\n"))
}
