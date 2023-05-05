package usage

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func stringToEntry(str string) *entry {
	subcommandAndArgs, choppedDescription, _ := strings.Cut(str, "\n"+strings.Repeat(Indent, 2))
	subcommandAndArgs = strings.TrimPrefix(subcommandAndArgs, Indent)
	subcommand, args, _ := strings.Cut(subcommandAndArgs, " ")

	var description strings.Builder
	for _, line := range strings.Split(choppedDescription, "\n"+strings.Repeat(Indent, 2)) {
		if line == "" {
			description.WriteString("\n\n")
		} else {
			description.WriteString(" " + line)
		}
	}

	argSlc := newArgSlice(args)

	return &entry{
		name:        subcommand,
		Description: strings.Replace(description.String(), "> <", fmt.Sprintf("> %s <", strings.Repeat("a", 72)), 1),
		args:        argSlc,
	}
}

func stringToMultipleEntries(str string) []entry {
	splitter := regexp.MustCompile(`[\n:]\n` + Indent)
	entryStrings := splitter.Split(str, -1)
	entries := make([]entry, 0)
	for _, e := range entryStrings[1:] {
		entries = append(entries, *stringToEntry(strings.TrimSpace(e)))
	}
	return entries
}

func assertOptionSlice(t *testing.T, got, want []option) {
	if len(got) != len(want) {
		t.Fatalf("%d options returned but wanted %d", len(got), len(want))
	}
	for i, gotOption := range got {
		assertOptionStruct(t, &gotOption, &want[i])
	}
}

type entryArgsTester struct {
	oArgs []string
}

func (tester entryArgsTester) assertEntryArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := entry{args: tester.oArgs}
		got := sampleEntry.Args()
		assertArgSlice(t, got, tester.oArgs)
	}
}

type entryOptionsTester struct {
	oOptions []option
}

func (tester entryOptionsTester) assertEntryOptions() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := entry{options: tester.oOptions}
		got := sampleEntry.Options()
		assertOptionSlice(t, got, tester.oOptions)
	}
}

type entryAddArgTester struct {
	iArg string
	oErr error
}

func (tester entryAddArgTester) assertEntryArgs() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		tempArgs := make([]string, 0, iterations)
		sampleEntry := entry{args: make([]string, 0)}
		for i := 1; i <= iterations; i++ {
			if gotErr := sampleEntry.AddArg(tester.iArg); gotErr != nil {
				t.Errorf("got %q error but should be nil", gotErr)
			}
			tempArgs = append(tempArgs, tester.iArg)
		}
		assertArgSlice(t, sampleEntry.args, tempArgs)
	}
}

func (tester entryAddArgTester) assertErrEmptyArgString() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := entry{args: make([]string, 0)}
		got := sampleEntry.AddArg(tester.iArg)
		if got == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		assertError(t, got, tester.oErr)
	}
}

type entryAddOptionTester struct {
	iOption *option
	oErr    error
}

func (tester entryAddOptionTester) assertEntryOptions() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		tempOptions := make([]option, 0, iterations)
		sampleEntry := entry{options: make([]option, 0)}
		for i := 1; i <= iterations; i++ {
			if gotErr := sampleEntry.AddOption(tester.iOption); gotErr != nil {
				t.Errorf("got %q error but should be nil", gotErr)
			}
			tempOptions = append(tempOptions, *tester.iOption)
		}
		assertOptionSlice(t, sampleEntry.options, tempOptions)
	}
}

func (tester entryAddOptionTester) assertErrNoOptionProvided() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := entry{options: make([]option, 0)}
		got := sampleEntry.AddOption(tester.iOption)
		if got == nil {
			t.Fatal("no error returned with nil option")
		}
		assertError(t, got, tester.oErr)
	}
}

type entrySetNameTester struct {
	iName string
	oErr  error
}

func (tester entrySetNameTester) assertEntryName() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := entry{name: "bar"}
		if gotErr := sampleEntry.SetName(tester.iName); gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		if sampleEntry.name != tester.iName {
			t.Errorf("name is %q but should be %q", sampleEntry.name, tester.iName)
		}
	}
}

func (tester entrySetNameTester) assertErrEmptyNameString() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := entry{name: "bar"}
		got := sampleEntry.SetName(tester.iName)
		if got == nil {
			t.Fatal("no error returned with empty name string")
		}
		assertError(t, got, tester.oErr)
	}
}

type newEntryTester struct {
	iName        string
	iDescription string
	oEntry       *entry
	oErr         error
}

func (tester newEntryTester) assertEntry() func(*testing.T) {
	return func(t *testing.T) {
		got, gotErr := NewEntry(tester.iName, tester.iDescription)
		if gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		if got.name != tester.oEntry.name {
			t.Errorf("name is %q but should be %q", got.name, tester.oEntry.name)
		}
		if got.Description != tester.oEntry.Description {
			t.Errorf("description is %q but should be %q", got.Description, tester.oEntry.Description)
		}
		if got.args == nil || len(got.args) != 0 {
			t.Error("args not initialized to an empty slice")
		}
		if got.options == nil || len(got.options) != 0 {
			t.Error("options not initialized to an empty slice")
		}
	}
}

func (tester newEntryTester) assertErrEmptyNameString() func(*testing.T) {
	return func(t *testing.T) {
		gotEntry, got := NewEntry(tester.iName, tester.iDescription)
		if gotEntry != nil {
			t.Errorf("got %+v entry but should be nil", gotEntry)
		}
		if got == nil {
			t.Fatal("no error returned with an empty name string")
		}
		assertError(t, got, tester.oErr)
	}
}

type entryStringTester struct {
	oString string
}

func (tester entryStringTester) assertString() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := stringToEntry(tester.oString)
		if got := sampleEntry.String(); got != tester.oString {
			t.Errorf("string is %q but should be %q", got, tester.oString)
		}
	}
}

func TestEntryArgs(t *testing.T) {
	t.Run("baseline", entryArgsTester{
		oArgs: []string{"foo"},
	}.assertEntryArgs())
	t.Run("multiple args", entryArgsTester{
		oArgs: []string{"foo", "bar", "baz"},
	}.assertEntryArgs())
	t.Run("no args", entryArgsTester{
		oArgs: make([]string, 0),
	}.assertEntryArgs())
}

func TestEntryOptions(t *testing.T) {
	t.Run("baseline", entryOptionsTester{
		oOptions: []option{{
			aliases:     []string{"foo"},
			Description: "foo",
			args:        []string{"foo"},
		}},
	}.assertEntryOptions())
	t.Run("multiple options", entryOptionsTester{
		oOptions: []option{
			{
				aliases:     []string{"foo"},
				Description: "foo",
				args:        []string{"foo"},
			},
			{
				aliases:     []string{"bar"},
				Description: "bar",
				args:        []string{"bar"},
			},
			{
				aliases:     []string{"baz"},
				Description: "baz",
				args:        []string{"baz"},
			},
		},
	}.assertEntryOptions())
	t.Run("no options", entryOptionsTester{
		oOptions: make([]option, 0),
	}.assertEntryOptions())
}

func TestEntryAddArg(t *testing.T) {
	t.Run("baseline", entryAddArgTester{
		iArg: "foo",
	}.assertEntryArgs())
	t.Run("empty arg string", entryAddArgTester{
		oErr: emptyArgStringErr(),
	}.assertErrEmptyArgString())
}

func TestEntryAddOption(t *testing.T) {
	t.Run("baseline", entryAddOptionTester{
		iOption: &option{
			aliases:     []string{"foo"},
			Description: "foo",
			args:        []string{"foo"},
		},
	}.assertEntryOptions())
	t.Run("nil option", entryAddOptionTester{
		oErr: nilOptionProvidedErr(),
	}.assertErrNoOptionProvided())
}

func TestNewEntry(t *testing.T) {
	t.Run("baseline", newEntryTester{
		iName:        "foo",
		iDescription: "foo",
		oEntry:       &entry{name: "foo", Description: "foo"},
	}.assertEntry())
	t.Run("empty description string", newEntryTester{
		iName:  "foo",
		oEntry: &entry{name: "foo"},
	}.assertEntry())
	t.Run("empty name string", newEntryTester{
		iDescription: "foo",
		oErr:         emptyNameStringErr(),
	}.assertErrEmptyNameString())
}

func TestEntrySetName(t *testing.T) {
	t.Run("baseline", entrySetNameTester{
		iName: "foo",
	}.assertEntryName())
	t.Run("empty name string", entrySetNameTester{
		oErr: emptyNameStringErr(),
	}.assertErrEmptyNameString())
}

func TestEntryString(t *testing.T) {
	longDescription := "some very long description that will definitely push the limits\n" +
		Indent + Indent + "of the screen size (it is very likely that this will cause the\n" +
		Indent + Indent + "line break at 64 characters)\n" +
		Indent + Indent + "\n" +
		Indent + Indent + "here's another paragraph just in case with a very long word\n" +
		Indent + Indent + "between these brackets > < that will not appear in the final\n" +
		Indent + Indent + "output because it is longer than a line"

	t.Run("baseline", entryStringTester{
		oString: Indent + "base",
	}.assertString())
	t.Run("description", entryStringTester{
		oString: fmt.Sprintf(
			"%sdescription\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("long description", entryStringTester{
		oString: fmt.Sprintf(
			"%slong-description\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("single arg", entryStringTester{
		oString: Indent + "single-arg <foo>",
	}.assertString())
	t.Run("multiple args", entryStringTester{
		oString: Indent + "multiple-args <foo> <bar> <baz>",
	}.assertString())
	t.Run("description single arg", entryStringTester{
		oString: fmt.Sprintf(
			"%sdescription-single-arg <foo>\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("description multiple args", entryStringTester{
		oString: fmt.Sprintf(
			"%sdescription-multiple-args <foo> <bar> <baz>\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("long description single arg", entryStringTester{
		oString: fmt.Sprintf(
			"%slong-description-single-arg <foo>\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("long description multiple args", entryStringTester{
		oString: fmt.Sprintf(
			"%slong-description-multiple-args <foo> <bar> <baz>\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("option", entryStringTester{
		oString: Indent + "option",
	}.assertString())
	t.Run("option description", entryStringTester{
		oString: fmt.Sprintf(
			"%soption-description\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("option long description", entryStringTester{
		oString: fmt.Sprintf(
			"%soption-long-description\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("option single arg", entryStringTester{
		oString: Indent + "option-single-arg <foo>",
	}.assertString())
	t.Run("option multiple args", entryStringTester{
		oString: Indent + "option-multiple-args <foo> <bar> <baz>",
	}.assertString())
	t.Run("option description single arg", entryStringTester{
		oString: fmt.Sprintf(
			"%soption-description-single-arg <foo>\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("option description multiple args", entryStringTester{
		oString: fmt.Sprintf(
			"%soption-description-multiple-args <foo> <bar> <baz>\n%ssome description",
			Indent,
			strings.Repeat(Indent, 2),
		),
	}.assertString())
	t.Run("option long description single arg", entryStringTester{
		oString: fmt.Sprintf(
			"%soption-long-description-single-arg <foo>\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
	t.Run("option long description multiple args", entryStringTester{
		oString: fmt.Sprintf(
			"%soption-long-description-multiple-args <foo> <bar> <baz>\n%s%s",
			Indent,
			strings.Repeat(Indent, 2),
			longDescription,
		),
	}.assertString())
}
