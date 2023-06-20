package main

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
		_, msg, _ := strings.Cut(tester.oString, ": ")
		sampleUsageError := &UsageError{msg}
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
		sampleUsageError := &UsageError{"foo"}
		if got := sampleUsageError.Is(tester.iTarget); got != tester.oBool {
			t.Errorf("result is %t but should be %t", got, tester.oBool)
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
