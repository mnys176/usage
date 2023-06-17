package usage

import (
	"strings"
	"testing"
)

/***** Helpers ************************************************/

func assertArgSlice(t *testing.T, got, want ArgSlice) {
	if len(got) != len(want) {
		t.Fatalf("%d args returned but wanted %d", len(got), len(want))
	}
	for i, gotArg := range got {
		if gotArg != want[i] {
			t.Errorf("arg is %q but should be %q", gotArg, want[i])
		}
	}
}

/***** Testers ************************************************/

type argSliceStringTester struct {
	oString string
}

func (tester argSliceStringTester) assertString() func(*testing.T) {
	return func(t *testing.T) {
		sampleArgSlice := ArgSlice(strings.Split(tester.oString[1:len(tester.oString)-1], "> <"))
		if got := sampleArgSlice.String(); got != tester.oString {
			t.Errorf("string is %q but should be %q", got, tester.oString)
		}
	}
}

func (tester argSliceStringTester) assertEmptyString() func(*testing.T) {
	return func(t *testing.T) {
		sampleArgSlice := ArgSlice{}
		if got := sampleArgSlice.String(); got != tester.oString {
			t.Errorf("string is %q but should be empty", got)
		}
	}
}

/***** Test Cases *********************************************/

func TestArgSliceString(t *testing.T) {
	t.Run("baseline", argSliceStringTester{
		oString: "<foo>",
	}.assertString())
	t.Run("multiple args", argSliceStringTester{
		oString: "<foo> <bar> <baz>",
	}.assertString())
	t.Run("no args", argSliceStringTester{}.assertEmptyString())
}
