package usage

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func assertAliases(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d aliases returned but wanted %d", len(got), len(want))
	}
	for i, gotAlias := range got {
		if gotAlias != want[i] {
			t.Errorf("alias is %q but should be %q", gotAlias, want[i])
		}
	}
}

func assertOption(t *testing.T, got, want *Option) {
	assertAliases(t, got.aliases, want.aliases)
	if got.Description != want.Description {
		t.Errorf("description is %q but should be %q", got.Description, want.Description)
	}
	if got.args == nil || len(got.args) != 0 {
		t.Error("args not initialized to an empty slice")
	}
	if got.Tmpl == "" {
		t.Error("template not initialized to a template")
	}
	got.Usage()
}

func assertError(t *testing.T, got, want error) {
	if !errors.Is(got, want) {
		t.Errorf("got %q error but wanted %q", got, want)
	}
}

func stringToOption(str string) *Option {
	indent := "    "
	aliasesAndArgs, choppedDescription, _ := strings.Cut(str, "\n"+strings.Repeat(indent, 2))
	aliasesAndArgs = strings.TrimPrefix(aliasesAndArgs, indent)

	aliasesString, argsString, _ := strings.Cut(aliasesAndArgs, " | ")
	aliases := strings.Split(aliasesString, " ")
	args := make([]string, 0)
	if len(argsString) > 0 {
		args = strings.Split(argsString, " ")
	}

	var description strings.Builder
	for _, line := range strings.Split(choppedDescription, "\n"+strings.Repeat(indent, 2)) {
		if line == "" {
			description.WriteString("\n\n")
		} else {
			description.WriteString(" " + line)
		}
	}

	return &Option{
		Description: strings.Replace(strings.TrimSpace(description.String()), "> <", fmt.Sprintf("> %s <", strings.Repeat("a", 72)), 1),
		Tmpl: `    {{join .Aliases " "}}{{if .Args}} | {{join .Args " "}}{{end}}{{if .Description}}
        {{with chop .Description 64}}{{join . "\n        "}}{{end}}{{end}}`,
		aliases: aliases,
		args:    args,
	}
}

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
			if gotErr := sampleOption.AddArg(tester.iArg); gotErr != nil {
				t.Errorf("got %q error but should be nil", gotErr)
			}
			args = append(args, tester.iArg)
		}
		assertArgs(t, sampleOption.args, args)
	}
}

func (tester optionAddArgTester) assertEmptyArgStringErr() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{args: make([]string, 0)}
		got := sampleOption.AddArg(tester.iArg)
		if got == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		assertError(t, got, tester.oErr)
	}
}

type optionSetAliasesTester struct {
	iAliases []string
	oErr     error
}

func (tester optionSetAliasesTester) assertAliases() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: tester.iAliases}
		if gotErr := sampleOption.SetAliases(tester.iAliases); gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		assertAliases(t, sampleOption.aliases, tester.iAliases)
	}
}

func (tester optionSetAliasesTester) assertNoAliasesErr() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: []string{"foo"}}
		got := sampleOption.SetAliases(tester.iAliases)
		if got == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		assertError(t, got, tester.oErr)
	}
}

func (tester optionSetAliasesTester) assertEmptyAliasStringErr() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{aliases: []string{"foo"}}
		got := sampleOption.SetAliases(tester.iAliases)
		if got == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		assertError(t, got, tester.oErr)
	}
}

type optionUsageTester struct {
	oUsage string
}

func (tester optionUsageTester) assertUsage() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := stringToOption(tester.oUsage)
		got, gotErr := sampleOption.Usage()
		if gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		if got != tester.oUsage {
			t.Errorf("usage is %q but should be %q", got, tester.oUsage)
		}
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
		assertOption(t, got, tester.oOption)
	}
}

func (tester newOptionTester) assertNoAliasesErr() func(*testing.T) {
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

func (tester newOptionTester) assertEmptyAliasStringErr() func(*testing.T) {
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
	}.assertEmptyArgStringErr())
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
	}.assertNoAliasesErr())
	t.Run("no aliases", optionSetAliasesTester{
		iAliases: make([]string, 0),
		oErr:     errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesErr())
	t.Run("single empty alias string", optionSetAliasesTester{
		iAliases: []string{""},
		oErr:     errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringErr())
	t.Run("multiple empty alias strings", optionSetAliasesTester{
		iAliases: []string{"foo", "", "bar", ""},
		oErr:     errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringErr())
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
		oErr:         errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesErr())
	t.Run("no aliases", newOptionTester{
		iAliases:     make([]string, 0),
		iDescription: "foo",
		oErr:         errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesErr())
	t.Run("single empty alias string", newOptionTester{
		iAliases:     []string{""},
		iDescription: "foo",
		oErr:         errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringErr())
	t.Run("multiple empty alias strings", newOptionTester{
		iAliases:     []string{"foo", "", "bar", ""},
		iDescription: "foo",
		oErr:         errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringErr())
}
