package usage

import (
	"errors"
	"strings"
	"testing"
)

type usageErrorErrorTester struct {
	oString string
}

func (tester usageErrorErrorTester) assertString() func(*testing.T) {
	return func(t *testing.T) {
		_, err, _ := strings.Cut(tester.oString, ": ")
		sampleUsageError := &UsageError{errors.New(err)}
		if got := sampleUsageError.Error(); got != tester.oString {
			t.Errorf("error is %q but should be %q", got, tester.oString)
		}
	}
}

type usageErrorIsTester struct {
	iTarget error
	oBool   bool
}

func (tester usageErrorIsTester) assertBool() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsageError := &UsageError{errors.New("foo")}
		if got := sampleUsageError.Is(tester.iTarget); got != tester.oBool {
			t.Errorf("result is %t but should be %t", got, tester.oBool)
		}
	}
}

type usageErrorUnwrapTester struct {
	oErr error
}

func (tester usageErrorUnwrapTester) assertError() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsageError := &UsageError{tester.oErr}
		if got := sampleUsageError.Unwrap(); !errors.Is(got, tester.oErr) {
			t.Errorf("error is %q but should be %q", got.Error(), tester.oErr.Error())
		}
	}
}

func (tester usageErrorUnwrapTester) assertNil() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsageError := &UsageError{}
		if got := sampleUsageError.Unwrap(); got != nil {
			t.Errorf("error is %q but should be nil", got.Error())
		}
	}
}

func TestUsageErrorError(t *testing.T) {
	t.Run("baseline", usageErrorErrorTester{
		oString: "usage: foo",
	}.assertString())
}

func TestUsageErrorIs(t *testing.T) {
	t.Run("baseline", usageErrorIsTester{
		iTarget: errors.New("usage: foo"),
		oBool:   true,
	}.assertBool())
	t.Run("is not", usageErrorIsTester{
		iTarget: errors.New("usage: bar"),
		oBool:   false,
	}.assertBool())
}

func TestUsageErrorUnwrap(t *testing.T) {
	t.Run("baseline", usageErrorUnwrapTester{
		oErr: errors.New("foo"),
	}.assertError())
	t.Run("nil wrapped error", usageErrorUnwrapTester{}.assertNil())
}
