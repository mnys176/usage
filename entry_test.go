package usage

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

/***** Helpers ************************************************/

func stringToEntry(str string) *Entry {
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

	argSlc := NewArgSlice(args)

	return &Entry{
		name:        subcommand,
		Description: strings.Replace(description.String(), "> <", fmt.Sprintf("> %s <", strings.Repeat("a", 72)), 1),
		args:        argSlc,
	}
}

func stringToMultipleEntries(str string) []Entry {
	splitter := regexp.MustCompile(`[\n:]\n` + Indent)
	entryStrings := splitter.Split(str, -1)
	entries := make([]Entry, 0)
	for _, e := range entryStrings[1:] {
		entries = append(entries, *stringToEntry(strings.TrimSpace(e)))
	}
	return entries
}

func assertOptionSlice(t *testing.T, got, want []Option) {
	if len(got) != len(want) {
		t.Fatalf("%d options returned but wanted %d", len(got), len(want))
	}
	for i, gotOption := range got {
		assertOptionStruct(t, &gotOption, &want[i])
	}
}

/***** Testers ************************************************/

type entryArgsTester struct {
	oArgs []string
}

func (tester entryArgsTester) assertEntryArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{args: tester.oArgs}
		got := sampleEntry.Args()
		assertArgSlice(t, got, tester.oArgs)
	}
}

type entryOptionsTester struct {
	oOptions []Option
}

func (tester entryOptionsTester) assertEntryOptions() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{options: tester.oOptions}
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
		sampleEntry := Entry{args: make([]string, 0)}
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
		sampleEntry := Entry{args: make([]string, 0)}
		got := sampleEntry.AddArg(tester.iArg)
		if got == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		assertError(t, got, tester.oErr)
	}
}

type entryAddOptionTester struct {
	iOption *Option
	oErr    error
}

func (tester entryAddOptionTester) assertEntryOptions() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		tempOptions := make([]Option, 0, iterations)
		sampleEntry := Entry{options: make([]Option, 0)}
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
		sampleEntry := Entry{options: make([]Option, 0)}
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
		sampleEntry := Entry{name: "bar"}
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
		sampleEntry := Entry{name: "bar"}
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
	oEntry       *Entry
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

type entryUsageTester struct {
	oUsage string
}

func (tester entryUsageTester) assertString() func(*testing.T) {
	return func(t *testing.T) {
		summarySection, optionsSection, _ := splitUsage(tester.oUsage)

		var name string
		var args ArgSlice
		if summarySection != "" {
			name, args = stringToNameAndArgs(summarySection)
		}

		var sampleOptions []Option
		if optionsSection != "" {
			sampleOptions = stringToMultipleOptions(optionsSection)
		} else {
			sampleOptions = make([]Option, 0)
		}

		sampleEntry := Entry{
			name:    name,
			options: sampleOptions,
			args:    args,
		}

		if got := sampleEntry.Usage(); got != tester.oUsage {
			t.Errorf("string is %q but should be %q", got, tester.oUsage)
		}
	}
}

/***** Test Cases *********************************************/

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
		oOptions: []Option{{
			aliases:     []string{"foo"},
			Description: "foo",
			args:        []string{"foo"},
		}},
	}.assertEntryOptions())
	t.Run("multiple options", entryOptionsTester{
		oOptions: []Option{
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
		oOptions: make([]Option, 0),
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
		iOption: &Option{
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
		oEntry:       &Entry{name: "foo", Description: "foo"},
	}.assertEntry())
	t.Run("empty description string", newEntryTester{
		iName:  "foo",
		oEntry: &Entry{name: "foo"},
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

func TestEntryUsage(t *testing.T) {
	const dblIndent string = Indent + Indent
	const longDescription string = "some very long description that will definitely push the limits\n" +
		Indent + Indent + "of the screen size (it is very likely that this will cause the\n" +
		Indent + Indent + "line break at 64 characters)\n" +
		Indent + Indent + "\n" +
		Indent + Indent + "here's another paragraph just in case with a very long word\n" +
		Indent + Indent + "between these brackets > < that will not appear in the final\n" +
		Indent + Indent + "output because it is longer than a line"

	const (
		summaryOptionsString string = "[options]"

		arg1, arg2, arg3 string = "<foo>", "<bar>", "<baz>"

		option1SingleAlias     string = "--foo"
		option2SingleAlias     string = "--bar"
		option3SingleAlias     string = "--baz"
		option1MultipleAliases string = "--foo, -f"
		option2MultipleAliases string = "--bar, -b"
		option3MultipleAliases string = "--baz, -z"
	)

	var (
		singleArgString    string = arg1
		multipleArgsString string = fmt.Sprintf("%s %s %s", arg1, arg2, arg3)
	)

	var (
		singleOptionString                                           string = option1SingleAlias
		singleOptionDescriptionString                                string = fmt.Sprintf("%s\n%ssome description", option1SingleAlias, dblIndent)
		singleOptionLongDescriptionString                            string = fmt.Sprintf("%s\n%s%s", option1SingleAlias, dblIndent, longDescription)
		singleOptionSingleArgString                                  string = fmt.Sprintf("%s %s", option1SingleAlias, singleArgString)
		singleOptionSingleArgDescriptionString                       string = fmt.Sprintf("%s %s\n%ssome description", option1SingleAlias, singleArgString, dblIndent)
		singleOptionSingleArgLongDescriptionString                   string = fmt.Sprintf("%s %s\n%s%s", option1SingleAlias, singleArgString, dblIndent, longDescription)
		singleOptionMultipleArgsString                               string = fmt.Sprintf("%s %s", option1SingleAlias, multipleArgsString)
		singleOptionMultipleArgsDescriptionString                    string = fmt.Sprintf("%s %s\n%ssome description", option1SingleAlias, multipleArgsString, dblIndent)
		singleOptionMultipleArgsLongDescriptionString                string = fmt.Sprintf("%s %s\n%s%s", option1SingleAlias, multipleArgsString, dblIndent, longDescription)
		singleOptionMultipleAliasesString                            string = option1MultipleAliases
		singleOptionMultipleAliasesDescriptionString                 string = fmt.Sprintf("%s\n%ssome description", option1MultipleAliases, dblIndent)
		singleOptionMultipleAliasesLongDescriptionString             string = fmt.Sprintf("%s\n%s%s", option1MultipleAliases, dblIndent, longDescription)
		singleOptionMultipleAliasesSingleArgString                   string = fmt.Sprintf("%s %s", option1MultipleAliases, singleArgString)
		singleOptionMultipleAliasesSingleArgDescriptionString        string = fmt.Sprintf("%s %s\n%ssome description", option1MultipleAliases, singleArgString, dblIndent)
		singleOptionMultipleAliasesSingleArgLongDescriptionString    string = fmt.Sprintf("%s %s\n%s%s", option1MultipleAliases, singleArgString, dblIndent, longDescription)
		singleOptionMultipleAliasesMultipleArgsString                string = fmt.Sprintf("%s %s", option1MultipleAliases, multipleArgsString)
		singleOptionMultipleAliasesMultipleArgsDescriptionString     string = fmt.Sprintf("%s %s\n%ssome description", option1MultipleAliases, multipleArgsString, dblIndent)
		singleOptionMultipleAliasesMultipleArgsLongDescriptionString string = fmt.Sprintf("%s %s\n%s%s", option1MultipleAliases, multipleArgsString, dblIndent, longDescription)
		multipleOptionsString                                        string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			option1SingleAlias,
			Indent+option2SingleAlias,
			Indent+option3SingleAlias,
		)
		multipleOptionsDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s\n%ssome description", option1SingleAlias, dblIndent),
			Indent+fmt.Sprintf("%s\n%ssome description", option2SingleAlias, dblIndent),
			Indent+fmt.Sprintf("%s\n%ssome description", option3SingleAlias, dblIndent),
		)
		multipleOptionsLongDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s\n%s%s", option1SingleAlias, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s\n%s%s", option2SingleAlias, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s\n%s%s", option3SingleAlias, dblIndent, longDescription),
		)
		multipleOptionsSingleArgString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s", option1SingleAlias, singleArgString),
			Indent+fmt.Sprintf("%s %s", option2SingleAlias, singleArgString),
			Indent+fmt.Sprintf("%s %s", option3SingleAlias, singleArgString),
		)
		multipleOptionsSingleArgDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%ssome description", option1SingleAlias, singleArgString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", option2SingleAlias, singleArgString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", option3SingleAlias, singleArgString, dblIndent),
		)
		multipleOptionsSingleArgLongDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%s%s", option1SingleAlias, singleArgString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", option2SingleAlias, singleArgString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", option3SingleAlias, singleArgString, dblIndent, longDescription),
		)
		multipleOptionsMultipleArgsString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s", option1SingleAlias, multipleArgsString),
			Indent+fmt.Sprintf("%s %s", option2SingleAlias, multipleArgsString),
			Indent+fmt.Sprintf("%s %s", option3SingleAlias, multipleArgsString),
		)
		multipleOptionsMultipleArgsDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%ssome description", option1SingleAlias, multipleArgsString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", option2SingleAlias, multipleArgsString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", option3SingleAlias, multipleArgsString, dblIndent),
		)
		multipleOptionsMultipleArgsLongDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%s%s", option1SingleAlias, multipleArgsString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", option2SingleAlias, multipleArgsString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", option3SingleAlias, multipleArgsString, dblIndent, longDescription),
		)
		multipleOptionsMultipleAliasesString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			option1MultipleAliases,
			Indent+option2MultipleAliases,
			Indent+option3MultipleAliases,
		)
		multipleOptionsMultipleAliasesDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s\n%ssome description", option1MultipleAliases, dblIndent),
			Indent+fmt.Sprintf("%s\n%ssome description", option2MultipleAliases, dblIndent),
			Indent+fmt.Sprintf("%s\n%ssome description", option3MultipleAliases, dblIndent),
		)
		multipleOptionsMultipleAliasesLongDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s\n%s%s", option1MultipleAliases, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s\n%s%s", option2MultipleAliases, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s\n%s%s", option3MultipleAliases, dblIndent, longDescription),
		)
		multipleOptionsMultipleAliasesSingleArgString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s", option1MultipleAliases, singleArgString),
			Indent+fmt.Sprintf("%s %s", option2MultipleAliases, singleArgString),
			Indent+fmt.Sprintf("%s %s", option3MultipleAliases, singleArgString),
		)
		multipleOptionsMultipleAliasesSingleArgDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%ssome description", option1MultipleAliases, singleArgString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", option2MultipleAliases, singleArgString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", option3MultipleAliases, singleArgString, dblIndent),
		)
		multipleOptionsMultipleAliasesSingleArgLongDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%s%s", option1MultipleAliases, singleArgString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", option2MultipleAliases, singleArgString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", option3MultipleAliases, singleArgString, dblIndent, longDescription),
		)
		multipleOptionsMultipleAliasesMultipleArgsString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s", option1MultipleAliases, multipleArgsString),
			Indent+fmt.Sprintf("%s %s", option2MultipleAliases, multipleArgsString),
			Indent+fmt.Sprintf("%s %s", option3MultipleAliases, multipleArgsString),
		)
		multipleOptionsMultipleAliasesMultipleArgsDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%ssome description", option1MultipleAliases, multipleArgsString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", option2MultipleAliases, multipleArgsString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", option3MultipleAliases, multipleArgsString, dblIndent),
		)
		multipleOptionsMultipleAliasesMultipleArgsLongDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%s%s", option1MultipleAliases, multipleArgsString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", option2MultipleAliases, multipleArgsString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", option3MultipleAliases, multipleArgsString, dblIndent, longDescription),
		)
	)

	t.Run("baseline", entryUsageTester{
		oUsage: fmt.Sprintf("Usage:\n%s%%s base\n", Indent),
	}.assertString())
	t.Run("single arg", entryUsageTester{
		oUsage: fmt.Sprintf("Usage:\n%s%%s single-arg %s\n", Indent, singleArgString),
	}.assertString())
	t.Run("multiple args", entryUsageTester{
		oUsage: fmt.Sprintf("Usage:\n%s%%s multiple-args %s\n", Indent, multipleArgsString),
	}.assertString())
	t.Run("single option", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionString,
		),
	}.assertString())
	t.Run("single option description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionDescriptionString,
		),
	}.assertString())
	t.Run("single option long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionLongDescriptionString,
		),
	}.assertString())
	t.Run("single option single arg", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionSingleArgString,
		),
	}.assertString())
	t.Run("single option single arg description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("single option single arg long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple args", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleArgsString,
		),
	}.assertString())
	t.Run("single option multiple args description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple args long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-aliases %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesString,
		),
	}.assertString())
	t.Run("single option multiple aliases description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-aliases-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-aliases-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases single arg", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-aliases-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesSingleArgString,
		),
	}.assertString())
	t.Run("single option multiple aliases single arg description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-aliases-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases single arg long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-aliases-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases multiple args", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-aliases-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesMultipleArgsString,
		),
	}.assertString())
	t.Run("single option multiple aliases multiple args description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-aliases-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases multiple args long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-aliases-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsString,
		),
	}.assertString())
	t.Run("multiple options description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsDescriptionString,
		),
	}.assertString())
	t.Run("multiple options long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single arg", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsSingleArgString,
		),
	}.assertString())
	t.Run("multiple options single arg description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single arg long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple args", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleArgsString,
		),
	}.assertString())
	t.Run("multiple options multiple args description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple args long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-aliases %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-aliases-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-aliases-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases single arg", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-aliases-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesSingleArgString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases single arg description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-aliases-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases single arg long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-aliases-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases multiple args", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-aliases-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesMultipleArgsString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases multiple args description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-aliases-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases multiple args long description", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-aliases-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("single option single global arg", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-single-global-arg %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			singleArgString,
			Indent+singleOptionSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple global args", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s single-option-multiple-global-args %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			multipleArgsString,
			Indent+singleOptionMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single global arg", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-single-global-arg %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			singleArgString,
			Indent+multipleOptionsSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple global args", entryUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%s%%s multiple-options-multiple-global-args %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			multipleArgsString,
			Indent+multipleOptionsMultipleArgsLongDescriptionString,
		),
	}.assertString())
}
