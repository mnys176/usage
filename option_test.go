package usage

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

/***** Helpers ************************************************/

func stringToOption(str string) *Option {
	aliasesAndArgs, choppedDescription, _ := strings.Cut(str, "\n"+strings.Repeat(Indent, 2))
	aliasesAndArgs = strings.TrimPrefix(aliasesAndArgs, Indent)

	aliases := make([]string, 0)
	argSlc := make(argSlice, 0)
	for _, token := range strings.Split(aliasesAndArgs, " ") {
		if strings.HasPrefix(token, "-") {
			aliases = append(aliases, strings.Trim(token, "-,"))
		} else if strings.HasPrefix(token, "<") {
			argSlc = append(argSlc, strings.Trim(token, "<>"))
		}
	}

	var description strings.Builder
	for _, line := range strings.Split(choppedDescription, "\n"+strings.Repeat(Indent, 2)) {
		if line == "" {
			description.WriteString("\n\n")
		} else {
			description.WriteString(" " + line)
		}
	}

	return &Option{
		aliases:     aliases,
		Description: strings.Replace(strings.TrimSpace(description.String()), "> <", fmt.Sprintf("> %s <", strings.Repeat("a", 72)), 1),
		args:        argSlc,
	}
}

func stringToMultipleOptions(str string) []Option {
	splitter := regexp.MustCompile(`[\n:]\n` + Indent)
	optionStrings := splitter.Split(str, -1)
	options := make([]Option, 0)
	for _, o := range optionStrings[1:] {
		options = append(options, *stringToOption(strings.TrimSpace(o)))
	}
	return options
}

func assertAliasSlice(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d aliases returned but wanted %d", len(got), len(want))
	}
	for i, gotAlias := range got {
		if gotAlias != want[i] {
			t.Errorf("alias is %q but should be %q", gotAlias, want[i])
		}
	}
}

func assertOptionStruct(t *testing.T, got, want *Option) {
	assertAliasSlice(t, got.aliases, want.aliases)
	if got.Description != want.Description {
		t.Errorf("description is %q but should be %q", got.Description, want.Description)
	}
	assertArgSlice(t, got.args, want.args)
}

/***** Testers ************************************************/

type optionArgsTester struct {
	oArgs []string
}

func (tester optionArgsTester) assertOptionArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{args: tester.oArgs}
		got := sampleOption.Args()
		assertArgSlice(t, got, tester.oArgs)
	}
}

type optionAddArgTester struct {
	iArg string
	oErr error
}

func (tester optionAddArgTester) assertOptionArgs() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		tempArgs := make([]string, 0, iterations)
		sampleOption := Option{args: make([]string, 0)}
		for i := 1; i <= iterations; i++ {
			if gotErr := sampleOption.AddArg(tester.iArg); gotErr != nil {
				t.Errorf("got %q error but should be nil", gotErr)
			}
			tempArgs = append(tempArgs, tester.iArg)
		}
		assertArgSlice(t, sampleOption.args, tempArgs)
	}
}

func (tester optionAddArgTester) assertErrEmptyArgString() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{args: make([]string, 0)}
		got := sampleOption.AddArg(tester.iArg)
		if got == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		assertError(t, got, tester.oErr)
	}
}

type newOptionTester struct {
	iAliases     []string
	iDescription string
	oOption      *Option
	oErr         error
}

func (tester newOptionTester) assertOption() func(*testing.T) {
	return func(t *testing.T) {
		got, gotErr := NewOption(tester.iAliases, tester.iDescription)
		if gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		assertAliasSlice(t, got.aliases, tester.oOption.aliases)
		if got.Description != tester.oOption.Description {
			t.Errorf("description is %q but should be %q", got.Description, tester.oOption.Description)
		}
		if got.args == nil || len(got.args) != 0 {
			t.Error("args not initialized to an empty slice")
		}
	}
}

func (tester newOptionTester) assertErrNoAliasProvided() func(*testing.T) {
	return func(t *testing.T) {
		gotOption, got := NewOption(tester.iAliases, tester.iDescription)
		if gotOption != nil {
			t.Errorf("got %+v option but should be nil", gotOption)
		}
		if got == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		assertError(t, got, tester.oErr)
	}
}

func (tester newOptionTester) assertErrEmptyAliasString() func(*testing.T) {
	return func(t *testing.T) {
		gotOption, got := NewOption(tester.iAliases, tester.iDescription)
		if gotOption != nil {
			t.Errorf("got %+v option but should be nil", gotOption)
		}
		if got == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		assertError(t, got, tester.oErr)
	}
}

type optionSetAliasesTester struct {
	iAliases []string
	oErr     error
}

func (tester optionSetAliasesTester) assertOptionAliases() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: []string{"foo"}}
		if gotErr := sampleOption.SetAliases(tester.iAliases); gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		assertAliasSlice(t, sampleOption.aliases, tester.iAliases)
	}
}

func (tester optionSetAliasesTester) assertErrNoAliasProvided() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: []string{"foo"}}
		got := sampleOption.SetAliases(tester.iAliases)
		if got == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		assertError(t, got, tester.oErr)
	}
}

func (tester optionSetAliasesTester) assertErrEmptyAliasString() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: []string{"foo"}}
		got := sampleOption.SetAliases(tester.iAliases)
		if got == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		assertError(t, got, tester.oErr)
	}
}

type optionStringTester struct {
	oString string
}

func (tester optionStringTester) assertString() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := stringToOption(tester.oString)
		if got := sampleOption.String(); got != tester.oString {
			t.Errorf("string is %q but should be %q", got, tester.oString)
		}
	}
}

/***** Test Cases *********************************************/

func TestOptionArgs(t *testing.T) {
	t.Run("baseline", optionArgsTester{
		oArgs: []string{"foo"},
	}.assertOptionArgs())
	t.Run("multiple args", optionArgsTester{
		oArgs: []string{"foo", "bar", "baz"},
	}.assertOptionArgs())
	t.Run("no args", optionArgsTester{
		oArgs: make([]string, 0),
	}.assertOptionArgs())
}

func TestOptionAddArg(t *testing.T) {
	t.Run("baseline", optionAddArgTester{
		iArg: "foo",
	}.assertOptionArgs())
	t.Run("empty arg string", optionAddArgTester{
		oErr: emptyArgStringErr(),
	}.assertErrEmptyArgString())
}

func TestNewOption(t *testing.T) {
	t.Run("baseline", newOptionTester{
		iAliases:     []string{"foo"},
		iDescription: "foo",
		oOption: &Option{
			aliases:     []string{"foo"},
			Description: "foo",
		},
	}.assertOption())
	t.Run("multiple aliases", newOptionTester{
		iAliases:     []string{"foo", "bar"},
		iDescription: "foo",
		oOption: &Option{
			aliases:     []string{"foo", "bar"},
			Description: "foo",
		},
	}.assertOption())
	t.Run("empty description string", newOptionTester{
		iAliases: []string{"foo"},
		oOption:  &Option{aliases: []string{"foo"}},
	}.assertOption())
	t.Run("nil aliases", newOptionTester{
		iDescription: "foo",
		oErr:         noAliasProvidedErr(),
	}.assertErrNoAliasProvided())
	t.Run("no aliases", newOptionTester{
		iAliases:     make([]string, 0),
		iDescription: "foo",
		oErr:         noAliasProvidedErr(),
	}.assertErrNoAliasProvided())
	t.Run("single empty alias string", newOptionTester{
		iAliases:     []string{""},
		iDescription: "foo",
		oErr:         emptyAliasStringErr(),
	}.assertErrEmptyAliasString())
	t.Run("multiple empty alias strings", newOptionTester{
		iAliases:     []string{"foo", "", "bar", ""},
		iDescription: "foo",
		oErr:         emptyAliasStringErr(),
	}.assertErrEmptyAliasString())
}

func TestOptionString(t *testing.T) {
	longDescription := "some very long description that will definitely push the limits\n" +
		Indent + Indent + "of the screen size (it is very likely that this will cause the\n" +
		Indent + Indent + "line break at 64 characters)\n" +
		Indent + Indent + "\n" +
		Indent + Indent + "here's another paragraph just in case with a very long word\n" +
		Indent + Indent + "between these brackets > < that will not appear in the final\n" +
		Indent + Indent + "output because it is longer than a line"

	t.Run("baseline", optionStringTester{
		oString: Indent + "--base",
	}.assertString())
	t.Run("description", optionStringTester{
		oString: fmt.Sprintf(
			"%s--description\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("long description", optionStringTester{
		oString: fmt.Sprintf(
			"%s--long-description\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("single arg", optionStringTester{
		oString: Indent + "--single-arg <arg1>",
	}.assertString())
	t.Run("multiple args", optionStringTester{
		oString: Indent + "--multiple-args <arg1> <arg2> <arg3>",
	}.assertString())
	t.Run("description single arg", optionStringTester{
		oString: fmt.Sprintf(
			"%s--description-single-arg <arg1>\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("description multiple args", optionStringTester{
		oString: fmt.Sprintf(
			"%s--description-multiple-args <arg1> <arg2> <arg3>\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("long description single arg", optionStringTester{
		oString: fmt.Sprintf(
			"%s--long-description-single-arg <arg1>\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("long description multiple args", optionStringTester{
		oString: fmt.Sprintf(
			"%s--long-description-multiple-arg <arg1> <arg2> <arg3>\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("multiple aliases", optionStringTester{
		oString: Indent + "--multiple-aliases, --another-one",
	}.assertString())
	t.Run("multiple aliases description", optionStringTester{
		oString: fmt.Sprintf(
			"%s--multiple-aliases, --another-one\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("multiple aliases long description", optionStringTester{
		oString: fmt.Sprintf(
			"%s--multiple-aliases, --long-description\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("multiple aliases single arg", optionStringTester{
		oString: Indent + "--multiple-aliases, --another-one <foo>",
	}.assertString())
	t.Run("multiple aliases multiple args", optionStringTester{
		oString: Indent + "--multiple-aliases, --another-one <foo> <bar> <baz>",
	}.assertString())
	t.Run("multiple aliases description single arg", optionStringTester{
		oString: fmt.Sprintf(
			"%s--multiple-aliases, -m <foo>\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("multiple aliases long description single arg", optionStringTester{
		oString: fmt.Sprintf(
			"%s--multiple-aliases, --another-one <foo>\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("multiple aliases description single arg", optionStringTester{
		oString: fmt.Sprintf(
			"%s--multiple-aliases, --another-one <foo> <bar> <baz>\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("multiple aliases long description single arg", optionStringTester{
		oString: fmt.Sprintf(
			"%s--multiple-aliases, --another-one <foo> <bar> <baz>\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
}

func TestOptionSetAliases(t *testing.T) {
	t.Run("baseline", optionSetAliasesTester{
		iAliases: []string{"foo"},
	}.assertOptionAliases())
	t.Run("multiple aliases", optionSetAliasesTester{
		iAliases: []string{"foo", "bar"},
	}.assertOptionAliases())
	t.Run("nil aliases", optionSetAliasesTester{
		oErr: noAliasProvidedErr(),
	}.assertErrNoAliasProvided())
	t.Run("no aliases", optionSetAliasesTester{
		iAliases: make([]string, 0),
		oErr:     noAliasProvidedErr(),
	}.assertErrNoAliasProvided())
	t.Run("single empty alias string", optionSetAliasesTester{
		iAliases: []string{""},
		oErr:     emptyAliasStringErr(),
	}.assertErrEmptyAliasString())
	t.Run("multiple empty alias strings", optionSetAliasesTester{
		iAliases: []string{"foo", "", "bar", ""},
		oErr:     emptyAliasStringErr(),
	}.assertErrEmptyAliasString())
}
