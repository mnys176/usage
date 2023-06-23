package usage

import (
	"errors"
	"strings"
	"testing"
)

type usageErrorErrorTester struct {
	oErr string
}

func (tester usageErrorErrorTester) assertErrorString() func(*testing.T) {
	return func(t *testing.T) {
		_, err, _ := strings.Cut(tester.oErr, ": ")
		sampleUsageError := &UsageError{errors.New(err)}
		got := sampleUsageError.Error()
		assertErrorString(t, got, tester.oErr)
	}
}

type usageErrorIsTester struct {
	iTarget   error
	oEquality bool
}

func (tester usageErrorIsTester) assertErrorEquality() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsageError := &UsageError{errors.New("foo")}
		got := sampleUsageError.Is(tester.iTarget)
		assertErrorEquality(t, got, tester.oEquality)
	}
}

type usageErrorUnwrapTester struct {
	oErr error
}

func (tester usageErrorUnwrapTester) assertError() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsageError := &UsageError{tester.oErr}
		got := sampleUsageError.Unwrap()
		assertError(t, got, tester.oErr)
	}
}

func (tester usageErrorUnwrapTester) assertNil() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsageError := &UsageError{}
		got := sampleUsageError.Unwrap()
		assertNilError(t, got)
	}
}

func TestUsageErrorError(t *testing.T) {
	t.Run("baseline", usageErrorErrorTester{
		oErr: "usage: foo",
	}.assertErrorString())
}

func TestUsageErrorIs(t *testing.T) {
	t.Run("baseline", usageErrorIsTester{
		iTarget:   errors.New("usage: foo"),
		oEquality: true,
	}.assertErrorEquality())
	t.Run("is not", usageErrorIsTester{
		iTarget:   errors.New("usage: bar"),
		oEquality: false,
	}.assertErrorEquality())
}

func TestUsageErrorUnwrap(t *testing.T) {
	t.Run("baseline", usageErrorUnwrapTester{
		oErr: errors.New("foo"),
	}.assertError())
	t.Run("nil wrapped error", usageErrorUnwrapTester{}.assertNil())
}
