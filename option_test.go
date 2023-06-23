package usage

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type optionArgsTester struct {
	oArgs []string
}

func (tester optionArgsTester) assertArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{args: tester.oArgs}
		got := sampleOption.Args()
		assertArgs(t, got, tester.oArgs)
	}
}

type optionAliasesTester struct {
	oAliases []string
}

func (tester optionAliasesTester) assertAliases() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: tester.oAliases}
		got := sampleOption.Aliases()
		assertAliases(t, got, tester.oAliases)
	}
}

type optionAddArgTester struct {
	iArg string
	oErr error
}

func (tester optionAddArgTester) assertArgs() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		args := make([]string, 0, iterations)
		sampleOption := Option{args: make([]string, 0)}
		for i := 1; i <= iterations; i++ {
			gotErr := sampleOption.AddArg(tester.iArg)
			assertNilError(t, gotErr)
			args = append(args, tester.iArg)
		}
		assertArgs(t, sampleOption.args, args)
	}
}

func (tester optionAddArgTester) assertEmptyArgStringError() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{args: make([]string, 0)}
		got := sampleOption.AddArg(tester.iArg)
		assertEmptyArgStringError(t, got, tester.oErr)
	}
}

type optionSetAliasesTester struct {
	iAliases []string
	oErr     error
}

func (tester optionSetAliasesTester) assertAliases() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: tester.iAliases}
		gotErr := sampleOption.SetAliases(tester.iAliases)
		assertNilError(t, gotErr)
		assertAliases(t, sampleOption.aliases, tester.iAliases)
	}
}

func (tester optionSetAliasesTester) assertNoAliasesError() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: []string{"foo"}}
		got := sampleOption.SetAliases(tester.iAliases)
		assertNoAliasesError(t, got, tester.oErr)
	}
}

func (tester optionSetAliasesTester) assertEmptyAliasStringError() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: []string{"foo"}}
		got := sampleOption.SetAliases(tester.iAliases)
		assertEmptyAliasStringError(t, got, tester.oErr)
	}
}

type optionUsageTester struct {
	oUsage string
}

func (tester optionUsageTester) assertUsage() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := stringToOption(tester.oUsage)
		got, gotErr := sampleOption.Usage()
		assertNilError(t, gotErr)
		assertUsage(t, got, tester.oUsage)
	}
}

type newOptionTester struct {
	iAliases     []string
	iDescription string
	oOption      *Option
	oErr         error
}

func (tester newOptionTester) assertDefaultOption() func(*testing.T) {
	return func(t *testing.T) {
		got, gotErr := NewOption(tester.iAliases, tester.iDescription)
		assertNilError(t, gotErr)
		assertDefaultOption(t, got, tester.oOption)
	}
}

func (tester newOptionTester) assertNoAliasesError() func(*testing.T) {
	return func(t *testing.T) {
		gotOption, got := NewOption(tester.iAliases, tester.iDescription)
		assertNilOption(t, gotOption)
		assertNoAliasesError(t, got, tester.oErr)
	}
}

func (tester newOptionTester) assertEmptyAliasStringError() func(*testing.T) {
	return func(t *testing.T) {
		gotOption, got := NewOption(tester.iAliases, tester.iDescription)
		assertNilOption(t, gotOption)
		assertEmptyAliasStringError(t, got, tester.oErr)
	}
}

func TestOptionArgs(t *testing.T) {
	t.Run("baseline", optionArgsTester{
		oArgs: []string{"foo"},
	}.assertArgs())
	t.Run("multiple args", optionArgsTester{
		oArgs: []string{"foo", "bar", "baz"},
	}.assertArgs())
	t.Run("no args", optionArgsTester{
		oArgs: make([]string, 0),
	}.assertArgs())
}

func TestOptionAliases(t *testing.T) {
	t.Run("baseline", optionAliasesTester{
		oAliases: []string{"foo"},
	}.assertAliases())
	t.Run("multiple aliases", optionAliasesTester{
		oAliases: []string{"foo", "bar", "baz"},
	}.assertAliases())
	t.Run("no aliases", optionAliasesTester{
		oAliases: make([]string, 0),
	}.assertAliases())
}

func TestOptionAddArg(t *testing.T) {
	t.Run("baseline", optionAddArgTester{
		iArg: "foo",
	}.assertArgs())
	t.Run("empty arg string", optionAddArgTester{
		oErr: errors.New("usage: arg string must not be empty"),
	}.assertEmptyArgStringError())
}

func TestOptionSetAliases(t *testing.T) {
	t.Run("baseline", optionSetAliasesTester{
		iAliases: []string{"foo"},
	}.assertAliases())
	t.Run("multiple aliases", optionSetAliasesTester{
		iAliases: []string{"foo", "bar"},
	}.assertAliases())
	t.Run("nil aliases", optionSetAliasesTester{
		oErr: errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesError())
	t.Run("no aliases", optionSetAliasesTester{
		iAliases: make([]string, 0),
		oErr:     errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesError())
	t.Run("single empty alias string", optionSetAliasesTester{
		iAliases: []string{""},
		oErr:     errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringError())
	t.Run("multiple empty alias strings", optionSetAliasesTester{
		iAliases: []string{"foo", "", "bar", ""},
		oErr:     errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringError())
}

func TestOptionUsage(t *testing.T) {
	indent := "    "
	longDescription := "some very long description that will definitely push the limits\n" +
		indent + indent + "of the screen size (it is very likely that this will cause the\n" +
		indent + indent + "line break at 64 characters)\n" +
		indent + indent + "\n" +
		indent + indent + "here's another paragraph just in case with a very long word\n" +
		indent + indent + "between these brackets > < that will not appear in the final\n" +
		indent + indent + "output because it is longer than a line"

	t.Run("baseline", optionUsageTester{
		oUsage: indent + "--base",
	}.assertUsage())
	t.Run("description", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--description\n%ssome description",
			indent,
			strings.Repeat(indent, 2),
		),
	}.assertUsage())
	t.Run("long description", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--long-description\n%s%s",
			indent,
			strings.Repeat(indent, 2),
			longDescription,
		),
	}.assertUsage())
	t.Run("single arg", optionUsageTester{
		oUsage: indent + "--single-arg <arg1>",
	}.assertUsage())
	t.Run("multiple args", optionUsageTester{
		oUsage: indent + "--multiple-args <arg1> <arg2> <arg3>",
	}.assertUsage())
	t.Run("description single arg", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--description-single-arg <arg1>\n%ssome description",
			indent,
			strings.Repeat(indent, 2),
		),
	}.assertUsage())
	t.Run("description multiple args", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--description-multiple-args <arg1> <arg2> <arg3>\n%ssome description",
			indent,
			strings.Repeat(indent, 2),
		),
	}.assertUsage())
	t.Run("long description single arg", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--long-description-single-arg <arg1>\n%s%s",
			indent,
			strings.Repeat(indent, 2),
			longDescription,
		),
	}.assertUsage())
	t.Run("long description multiple args", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--long-description-multiple-arg <arg1> <arg2> <arg3>\n%s%s",
			indent,
			strings.Repeat(indent, 2),
			longDescription,
		),
	}.assertUsage())
	t.Run("multiple aliases", optionUsageTester{
		oUsage: indent + "--multiple-aliases, --another-one",
	}.assertUsage())
	t.Run("multiple aliases description", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--multiple-aliases, --another-one\n%ssome description",
			indent,
			strings.Repeat(indent, 2),
		),
	}.assertUsage())
	t.Run("multiple aliases long description", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--multiple-aliases, --long-description\n%s%s",
			indent,
			strings.Repeat(indent, 2),
			longDescription,
		),
	}.assertUsage())
	t.Run("multiple aliases single arg", optionUsageTester{
		oUsage: indent + "--multiple-aliases, --another-one <foo>",
	}.assertUsage())
	t.Run("multiple aliases multiple args", optionUsageTester{
		oUsage: indent + "--multiple-aliases, --another-one <foo> <bar> <baz>",
	}.assertUsage())
	t.Run("multiple aliases description single arg", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--multiple-aliases, -m <foo>\n%ssome description",
			indent,
			strings.Repeat(indent, 2),
		),
	}.assertUsage())
	t.Run("multiple aliases long description single arg", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--multiple-aliases, --another-one <foo>\n%s%s",
			indent,
			strings.Repeat(indent, 2),
			longDescription,
		),
	}.assertUsage())
	t.Run("multiple aliases description single arg", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--multiple-aliases, --another-one <foo> <bar> <baz>\n%ssome description",
			indent,
			strings.Repeat(indent, 2),
		),
	}.assertUsage())
	t.Run("multiple aliases long description single arg", optionUsageTester{
		oUsage: fmt.Sprintf(
			"%s--multiple-aliases, --another-one <foo> <bar> <baz>\n%s%s",
			indent,
			strings.Repeat(indent, 2),
			longDescription,
		),
	}.assertUsage())
}

func TestNewOption(t *testing.T) {
	t.Run("baseline", newOptionTester{
		iAliases:     []string{"foo"},
		iDescription: "foo",
		oOption: &Option{
			aliases:     []string{"foo"},
			Description: "foo",
		},
	}.assertDefaultOption())
	t.Run("multiple aliases", newOptionTester{
		iAliases:     []string{"foo", "bar"},
		iDescription: "foo",
		oOption: &Option{
			aliases:     []string{"foo", "bar"},
			Description: "foo",
		},
	}.assertDefaultOption())
	t.Run("empty description string", newOptionTester{
		iAliases: []string{"foo"},
		oOption:  &Option{aliases: []string{"foo"}},
	}.assertDefaultOption())
	t.Run("nil aliases", newOptionTester{
		iDescription: "foo",
		oErr:         errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesError())
	t.Run("no aliases", newOptionTester{
		iAliases:     make([]string, 0),
		iDescription: "foo",
		oErr:         errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesError())
	t.Run("single empty alias string", newOptionTester{
		iAliases:     []string{""},
		iDescription: "foo",
		oErr:         errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringError())
	t.Run("multiple empty alias strings", newOptionTester{
		iAliases:     []string{"foo", "", "bar", ""},
		iDescription: "foo",
		oErr:         errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringError())
}
