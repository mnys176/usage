package usage

import (
	"strings"
	"testing"
)

func assertArgSlice(t *testing.T, got, want argSlice) {
	if len(got) != len(want) {
		t.Fatalf("%d args returned but wanted %d", len(got), len(want))
	}
	for i, gotArg := range got {
		if gotArg != want[i] {
			t.Errorf("arg is %q but should be %q", gotArg, want[i])
		}
	}
}

// func assertArgSliceString(t *testing.T, got string) {
// 	rawPattern := `^((<.+?>)( <.+?>)*)?$`
// 	if matched, _ := regexp.MatchString(rawPattern, got); !matched {
// 		t.Errorf("arg slice string %q does not match pattern %q", got, rawPattern)
// 	}
// }

type argSliceStringTester struct {
	oString string
}

func (tester argSliceStringTester) assertString() func(*testing.T) {
	return func(t *testing.T) {
		sampleArgSlice := argSlice(strings.Split(tester.oString[1:len(tester.oString)-1], "> <"))
		if got := sampleArgSlice.String(); got != tester.oString {
			t.Errorf("string is %q but should be %q", got, tester.oString)
		}
	}
}

func (tester argSliceStringTester) assertEmptyString() func(*testing.T) {
	return func(t *testing.T) {
		sampleArgSlice := argSlice{}
		if got := sampleArgSlice.String(); got != tester.oString {
			t.Errorf("string is %q but should be empty", got)
		}
	}
}

func TestArgSliceString(t *testing.T) {
	t.Run("baseline", argSliceStringTester{
		oString: "<foo>",
	}.assertString())
	t.Run("multiple args", argSliceStringTester{
		oString: "<foo> <bar> <baz>",
	}.assertString())
	t.Run("no args", argSliceStringTester{}.assertEmptyString())
}
