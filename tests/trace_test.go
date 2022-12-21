package tests

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gabehardgrave/trace"
	"github.com/stretchr/testify/require"
)

func TestTraceIsNOOPForNil(t *testing.T) {
	err := trace.Wrap(nil)
	require.Nil(t, err)

	err = err.Trace()
	require.Nil(t, err)
}

func TestTraceCapturesFileAndLineNumber(t *testing.T) {
	ogError := errors.New("original error")

	err := trace.Wrap(ogError)
	expectedTrace := "original error [" +
		pwd(t, "/trace_test.go:24") +
		"]"
	require.Equal(t, expectedTrace, err.Error())

	err = err.Trace()
	expectedTrace = ("original error [" +
		pwd(t, "/trace_test.go:24, ") +
		pwd(t, "/trace_test.go:30") +
		"]")
	require.Equal(t, expectedTrace, err.Error())

	wrappedErr := fmt.Errorf("%w", err)
	tracedErr := trace.Wrap(wrappedErr)

	expectedTrace = ("original error [" +
		pwd(t, "/trace_test.go:24, ") +
		pwd(t, "/trace_test.go:30, ") +
		pwd(t, "/trace_test.go:38") +
		"]")
	require.Equal(t, expectedTrace, tracedErr.Error())

	unwrappedErr := errors.Unwrap(tracedErr)
	require.EqualError(t, unwrappedErr, "original error")
	require.ErrorIs(t, tracedErr, ogError)
}

func TestTraceCapturesExtraDetails(t *testing.T) {
	ogError := errors.New("original error")

	foo := struct {
		id   int
		name string
	}{
		id:   7,
		name: "yeet",
	}

	err := trace.Wrap(ogError, "foo=%v", foo)
	err = err.Trace(999, "bizz")

	expectedTrace := "original error [" +
		pwd(t, "/trace_test.go:63{ foo={7 yeet} }, ") +
		pwd(t, "/trace_test.go:64{ [999 bizz] }") +
		"]"

	require.Equal(t, expectedTrace, err.Error())
}

func pwd(t *testing.T, s string) string {
	dir, err := os.Getwd()
	require.NoError(t, err)
	return dir + s
}
