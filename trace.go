package trace

import (
	"errors"
	"runtime"
	"strings"
)

// Wrap wraps the error with tracing information
// (file and line #, as well as optional formatted
// context).
// The returned *Err implements the error interface,
// and can be dropped in where error is expected.
//
// Wrap returns nil if cause is nil.
func Wrap(cause error, details ...any) *Err {
	if cause == nil {
		return nil
	}

	_, file, line, _ := runtime.Caller(1) // location of whoever called Wrap()
	loc := ErrorLocation{file, line, details}

	err := &Err{}
	if errors.As(cause, &err) {
		err.trace = append(err.trace, loc)
	} else {
		err.cause = cause
		err.trace = append(make([]ErrorLocation, 0, 10), loc)
	}

	return err
}

type Err struct {
	cause error
	trace []ErrorLocation
}

// Trace extends the error with tracing information
// (file and line #, as well as optional formatted
// context).
//
// If the Err is nil, Trace returns nil
func (e *Err) Trace(details ...any) *Err {
	if e == nil || e.cause == nil {
		return nil
	}

	_, file, line, _ := runtime.Caller(1) // location of whoever called Trace()
	loc := ErrorLocation{file, line, details}
	e.trace = append(e.trace, loc)

	return e
}

func (e *Err) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *Err) Error() string {
	if e == nil || e.cause == nil {
		return ""
	}

	cause := e.cause.Error()
	var s strings.Builder

	// A minor optimization - try to avoid reallocations
	sizeHint := len(cause) + len(" []")
	for _, loc := range e.trace {
		sizeHint += loc.formatSizeHint()
		sizeHint += len(", ")
	}
	s.Grow(sizeHint)

	s.WriteString(cause)
	s.WriteString(" [")

	for i, loc := range e.trace {
		loc.format(&s)
		if i < len(e.trace)-1 {
			s.WriteString(", ")
		}
	}

	s.WriteString("]")

	return s.String()
}

func (e *Err) GetErrorLocations() []ErrorLocation {
	if e == nil || e.cause == nil {
		return []ErrorLocation{}
	}
	return e.trace
}
