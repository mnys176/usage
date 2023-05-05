package usage

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func splitNameAndArgs(summary string) (string, string) {
	var name, args string
	if argsStart := strings.IndexAny(summary, "<["); argsStart > -1 {
		name = summary[:argsStart-1]
		args = summary[argsStart:]
	} else {
		name = summary
	}
	return name, args
}

func splitUsage(usage string) (string, string, string) {
	summaryStart := strings.Index(usage, "Usage:")
	optionsStart := strings.Index(usage, "Options:")
	commandsStart := strings.Index(usage, "Commands:")

	var summarySection, optionsSection, commandsSection string
	if summaryStart > -1 {
		if optionsStart > -1 {
			summarySection = usage[summaryStart:optionsStart]
		} else if commandsStart > -1 {
			summarySection = usage[summaryStart:commandsStart]
		} else {
			summarySection = usage[summaryStart:]
		}
	}
	if optionsStart > -1 {
		if commandsStart > -1 {
			optionsSection = usage[optionsStart:commandsStart]
		} else {
			optionsSection = usage[optionsStart:]
		}
	}
	if commandsStart > -1 {
		commandsSection = usage[commandsStart:]
	}
	return summarySection, optionsSection, commandsSection
}

func stringToNameAndArgs(summarySection string) (string, argSlice) {
	splitter := regexp.MustCompile(`[\n:]\n` + Indent)
	summary := strings.TrimSpace(splitter.Split(summarySection, 3)[1])
	nameString, argString := splitNameAndArgs(summary)

	var name string
	if parts := strings.Split(nameString, " "); len(parts) > 1 {
		name = parts[len(parts)-1]
	} else {
		name = nameString
	}

	if strings.HasPrefix(argString, "<command>") {
		return name, newArgSlice("")
	}
	if argsStart := strings.IndexRune(argString, '<'); argsStart > -1 {
		return name, newArgSlice(argString[argsStart:])
	}
	return name, newArgSlice("")
}

func assertError(t *testing.T, got, want error) {
	if !errors.Is(got, want) {
		t.Errorf("got %q error but wanted %q", got, want)
	}
}

func assertEntryStruct(t *testing.T, got, want *entry) {
	if got.name != want.name {
		t.Errorf("name is %q but should be %q", got.name, want.name)
	}
	if got.Description != want.Description {
		t.Errorf("description is %q but should be %q", got.Description, want.Description)
	}
	assertArgSlice(t, got.args, want.args)
	assertOptionSlice(t, got.options, want.options)
}

func assertEntrySlice(t *testing.T, got, want []entry) {
	if len(got) != len(want) {
		t.Fatalf("%d entries returned but wanted %d", len(got), len(want))
	}
	for i, gotEntry := range got {
		assertEntryStruct(t, &gotEntry, &want[i])
	}
}

type usageArgsTester struct {
	oArgs []string
}

func (tester usageArgsTester) assertUsageArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{args: tester.oArgs}
		got := sampleUsage.Args()
		assertArgSlice(t, got, tester.oArgs)
	}
}

type usageOptionsTester struct {
	oOptions []option
}

func (tester usageOptionsTester) assertUsageOptions() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{options: tester.oOptions}
		got := sampleUsage.Options()
		assertOptionSlice(t, got, tester.oOptions)
	}
}

type usageEntriesTester struct {
	oEntries []entry
}

func (tester usageEntriesTester) assertUsageEntries() func(*testing.T) {
	return func(t *testing.T) {
		sort.Slice(tester.oEntries, func(i, j int) bool {
			return tester.oEntries[i].name < tester.oEntries[j].name
		})
		sampleUsage := usage{entries: make(map[string]entry)}
		for _, sampleEntry := range tester.oEntries {
			sampleUsage.entries[sampleEntry.name] = sampleEntry
		}
		got := sampleUsage.Entries()
		assertEntrySlice(t, got, tester.oEntries)
	}
}

type usageAddArgTester struct {
	iArg string
	oErr error
}

func (tester usageAddArgTester) assertUsageArgs() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		tempArgs := make([]string, 0, iterations)
		sampleUsage := usage{args: make([]string, 0)}
		for i := 1; i <= iterations; i++ {
			if gotErr := sampleUsage.AddArg(tester.iArg); gotErr != nil {
				t.Errorf("got %q error but should be nil", gotErr)
			}
			tempArgs = append(tempArgs, tester.iArg)
		}
		assertArgSlice(t, sampleUsage.args, tempArgs)
	}
}

func (tester usageAddArgTester) assertErrEmptyArgString() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{args: make([]string, 0)}
		got := sampleUsage.AddArg(tester.iArg)
		if got == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		assertError(t, got, tester.oErr)
	}
}

func (tester usageAddArgTester) assertErrExistingEntries() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{
			entries: map[string]entry{
				"foo": {
					name:        "foo",
					Description: "foo",
					options: []option{{
						aliases:     []string{"foo"},
						Description: "foo",
						args:        []string{"foo"},
					}},
					args: []string{"foo"},
				},
			},
			args: make([]string, 0),
		}

		got := sampleUsage.AddArg(tester.iArg)
		if got == nil {
			t.Fatal("no error returned with existing entries")
		}
		assertError(t, got, tester.oErr)
	}
}

type usageAddOptionTester struct {
	iOption *option
	oErr    error
}

func (tester usageAddOptionTester) assertUsageOptions() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		tempOptions := make([]option, 0, iterations)
		sampleUsage := usage{options: make([]option, 0)}
		for i := 1; i <= iterations; i++ {
			if gotErr := sampleUsage.AddOption(tester.iOption); gotErr != nil {
				t.Errorf("got %q error but should be nil", gotErr)
			}
			tempOptions = append(tempOptions, *tester.iOption)
		}
		assertOptionSlice(t, sampleUsage.options, tempOptions)
	}
}

func (tester usageAddOptionTester) assertErrNoOptionProvided() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{options: make([]option, 0)}
		got := sampleUsage.AddOption(tester.iOption)
		if got == nil {
			t.Fatal("no error returned with nil option")
		}
		assertError(t, got, tester.oErr)
	}
}

type usageAddEntryTester struct {
	iEntry *entry
	oErr   error
}

func (tester usageAddEntryTester) assertUsageEntries() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{entries: make(map[string]entry)}
		if gotErr := sampleUsage.AddEntry(tester.iEntry); gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		sampleEntries := make([]entry, 0)
		for _, sampleEntry := range sampleUsage.entries {
			sampleEntries = append(sampleEntries, sampleEntry)
		}
		sort.Slice(sampleEntries, func(i, j int) bool {
			return sampleEntries[i].name < sampleEntries[j].name
		})
		assertEntrySlice(t, sampleEntries, []entry{*tester.iEntry})
	}
}

func (tester usageAddEntryTester) assertErrNoEntryProvided() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{entries: make(map[string]entry)}
		got := sampleUsage.AddEntry(tester.iEntry)
		if got == nil {
			t.Fatal("no error returned with nil entry")
		}
		assertError(t, got, tester.oErr)
	}
}

func (tester usageAddEntryTester) assertErrExistingArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{
			entries: make(map[string]entry),
			args:    []string{"foo"},
		}
		got := sampleUsage.AddEntry(tester.iEntry)
		if got == nil {
			t.Fatal("no error returned with existing args")
		}
		assertError(t, got, tester.oErr)
	}
}

type usageSetNameTester struct {
	iName string
	oErr  error
}

func (tester usageSetNameTester) assertUsageName() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{name: "bar"}
		if gotErr := sampleUsage.SetName(tester.iName); gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		if sampleUsage.name != tester.iName {
			t.Errorf("name is %q but should be %q", sampleUsage.name, tester.iName)
		}
	}
}

func (tester usageSetNameTester) assertErrEmptyNameString() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{name: "bar"}
		got := sampleUsage.SetName(tester.iName)
		if got == nil {
			t.Fatal("no error returned with empty name string")
		}
		assertError(t, got, tester.oErr)
	}
}

type newUsageTester struct {
	iName  string
	oUsage *usage
	oErr   error
}

func (tester newUsageTester) assertUsage() func(*testing.T) {
	return func(t *testing.T) {
		got, gotErr := NewUsage(tester.iName)
		if gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		if got.name != tester.oUsage.name {
			t.Errorf("name is %q but should be %q", got.name, tester.oUsage.name)
		}
		if got.args == nil || len(got.args) != 0 {
			t.Error("args not initialized to an empty slice")
		}
		if got.options == nil || len(got.options) != 0 {
			t.Error("options not initialized to an empty slice")
		}
		if got.entries == nil || len(got.entries) != 0 {
			t.Error("entries not initialized to an empty map")
		}
	}
}

func (tester newUsageTester) assertErrEmptyNameString() func(*testing.T) {
	return func(t *testing.T) {
		gotUsage, got := NewUsage(tester.iName)
		if gotUsage != nil {
			t.Errorf("got %+v usage but should be nil", gotUsage)
		}
		if got == nil {
			t.Fatal("no error returned with an empty name string")
		}
		assertError(t, got, tester.oErr)
	}
}

func TestUsageArgs(t *testing.T) {
	t.Run("baseline", usageArgsTester{
		oArgs: []string{"foo"},
	}.assertUsageArgs())
	t.Run("multiple args", usageArgsTester{
		oArgs: []string{"foo", "bar", "baz"},
	}.assertUsageArgs())
	t.Run("no args", usageArgsTester{
		oArgs: make([]string, 0),
	}.assertUsageArgs())
}

func TestUsageOptions(t *testing.T) {
	t.Run("baseline", usageOptionsTester{
		oOptions: []option{{
			aliases:     []string{"foo"},
			Description: "foo",
			args:        []string{"foo"},
		}},
	}.assertUsageOptions())
	t.Run("multiple options", usageOptionsTester{
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
	}.assertUsageOptions())
	t.Run("no options", usageOptionsTester{
		oOptions: make([]option, 0),
	}.assertUsageOptions())
}

func TestUsageEntries(t *testing.T) {
	t.Run("baseline", usageEntriesTester{
		oEntries: []entry{{
			name:        "foo",
			Description: "foo",
			options: []option{{
				aliases:     []string{"foo"},
				Description: "foo",
				args:        []string{"foo"},
			}},
			args: []string{"foo"},
		}},
	}.assertUsageEntries())
	t.Run("multiple entries", usageEntriesTester{
		oEntries: []entry{
			{
				name:        "foo",
				Description: "foo",
				options: []option{{
					aliases:     []string{"foo"},
					Description: "foo",
					args:        []string{"foo"},
				}},
				args: []string{"foo"},
			},
			{
				name:        "bar",
				Description: "bar",
				options: []option{{
					aliases:     []string{"bar"},
					Description: "bar",
					args:        []string{"bar"},
				}},
				args: []string{"bar"},
			},
			{
				name:        "baz",
				Description: "baz",
				options: []option{{
					aliases:     []string{"baz"},
					Description: "baz",
					args:        []string{"baz"},
				}},
				args: []string{"baz"},
			},
		},
	}.assertUsageEntries())
	t.Run("no entries", usageEntriesTester{
		oEntries: make([]entry, 0),
	}.assertUsageEntries())
}

type usageUsageTester struct {
	oUsage string
}

func (tester usageUsageTester) assertString() func(*testing.T) {
	return func(t *testing.T) {
		summarySection, optionsSection, commandsSection := splitUsage(tester.oUsage)

		var name string
		var args argSlice
		if summarySection != "" {
			name, args = stringToNameAndArgs(summarySection)
		}

		var sampleOptions []option
		if optionsSection != "" {
			sampleOptions = stringToMultipleOptions(optionsSection)
		} else {
			sampleOptions = make([]option, 0)
		}

		var sampleEntries []entry
		if commandsSection != "" {
			sampleEntries = stringToMultipleEntries(commandsSection)
		}

		sampleUsage := usage{
			name:    name,
			options: sampleOptions,
			entries: make(map[string]entry),
			args:    args,
		}
		for _, e := range sampleEntries {
			sampleUsage.entries[e.name] = e
		}

		if got := sampleUsage.Usage(); got != tester.oUsage {
			t.Errorf("string is %q but should be %q", got, tester.oUsage)
		}
	}
}

type usageLookupTester struct {
	oUsage string
}

func (tester usageLookupTester) assertString() func(*testing.T) {
	return func(t *testing.T) {
		summarySection, optionsSection, _ := splitUsage(tester.oUsage)
		splitter := regexp.MustCompile(`[\n:]\n` + Indent)
		summary := strings.TrimSpace(splitter.Split(summarySection, 3)[1])

		nameString, _ := splitNameAndArgs(summary)
		parentName := strings.Split(nameString, " ")[0]

		var name string
		var args argSlice
		if summarySection != "" {
			name, args = stringToNameAndArgs(summarySection)
		}

		var sampleOptions []option
		if optionsSection != "" {
			sampleOptions = stringToMultipleOptions(optionsSection)
		} else {
			sampleOptions = make([]option, 0)
		}

		sampleEntry := entry{
			name:    name,
			options: sampleOptions,
			args:    args,
		}

		sampleUsage := usage{
			name:    parentName,
			entries: map[string]entry{sampleEntry.name: sampleEntry},
			options: sampleOptions,
			args:    args,
		}

		if got := sampleUsage.Lookup(sampleEntry.name); got != tester.oUsage {
			t.Errorf("string is %q but should be %q", got, tester.oUsage)
		}
	}
}

func (tester usageLookupTester) assertEmptyString() func(*testing.T) {
	return func(t *testing.T) {
		sampleUsage := usage{entries: make(map[string]entry)}
		if got := sampleUsage.Lookup("foo"); got != tester.oUsage {
			t.Errorf("string is %q but should be empty", got)
		}
	}
}

func TestUsageAddArg(t *testing.T) {
	t.Run("baseline", usageAddArgTester{
		iArg: "foo",
	}.assertUsageArgs())
	t.Run("empty arg string", usageAddArgTester{
		oErr: emptyArgStringErr(),
	}.assertErrEmptyArgString())
	t.Run("existing entries", usageAddArgTester{
		iArg: "foo",
		oErr: existingEntriesErr(),
	}.assertErrExistingEntries())
}

func TestUsageAddOption(t *testing.T) {
	t.Run("baseline", usageAddOptionTester{
		iOption: &option{
			aliases:     []string{"foo"},
			Description: "foo",
			args:        []string{"foo"},
		},
	}.assertUsageOptions())
	t.Run("nil option", usageAddOptionTester{
		oErr: nilOptionProvidedErr(),
	}.assertErrNoOptionProvided())
}

func TestUsageAddEntry(t *testing.T) {
	t.Run("baseline", usageAddEntryTester{
		iEntry: &entry{
			name:        "foo",
			Description: "foo",
			options: []option{{
				aliases:     []string{"foo"},
				Description: "foo",
				args:        []string{"foo"},
			}},
			args: []string{"foo"},
		},
	}.assertUsageEntries())
	t.Run("existing args", usageAddEntryTester{
		iEntry: &entry{
			name:        "foo",
			Description: "foo",
			options: []option{{
				aliases:     []string{"foo"},
				Description: "foo",
				args:        []string{"foo"},
			}},
			args: []string{"foo"},
		},
		oErr: existingArgsErr(),
	}.assertErrExistingArgs())
	t.Run("nil entry", usageAddEntryTester{
		oErr: nilEntryProvidedErr(),
	}.assertErrNoEntryProvided())
}

func TestNewUsage(t *testing.T) {
	t.Run("baseline", newUsageTester{
		iName:  "foo",
		oUsage: &usage{name: "foo"},
	}.assertUsage())
	t.Run("empty name string", newUsageTester{
		oErr: emptyNameStringErr(),
	}.assertErrEmptyNameString())
}

func TestUsageSetName(t *testing.T) {
	t.Run("baseline", usageSetNameTester{
		iName: "foo",
	}.assertUsageName())
	t.Run("empty name string", usageSetNameTester{
		oErr: emptyNameStringErr(),
	}.assertErrEmptyNameString())
}

func TestUsageUsage(t *testing.T) {
	const dblIndent string = Indent + Indent
	const longDescription string = "some very long description that will definitely push the limits\n" +
		Indent + Indent + "of the screen size (it is very likely that this will cause the\n" +
		Indent + Indent + "line break at 64 characters)\n" +
		Indent + Indent + "\n" +
		Indent + Indent + "here's another paragraph just in case with a very long word\n" +
		Indent + Indent + "between these brackets > < that will not appear in the final\n" +
		Indent + Indent + "output because it is longer than a line"

	const (
		summaryCommandsString  string = "<command>"
		summaryOptionsString   string = "[options]"
		summaryArgsString      string = "<args>"
		summaryExtensionFormat string = Indent + "To learn more about the available" +
			" options for each command, use the\n" + Indent + "--help flag" +
			" like so:\n\n" + Indent + "%s " + summaryCommandsString + " --help"

		arg1, arg2, arg3 string = "<foo>", "<bar>", "<baz>"

		option1SingleAlias     string = "--foo"
		option2SingleAlias     string = "--bar"
		option3SingleAlias     string = "--baz"
		option1MultipleAliases string = "--foo, -f"
		option2MultipleAliases string = "--bar, -b"
		option3MultipleAliases string = "--baz, -z"

		command1Name string = "bar"
		command2Name string = "baz"
		command3Name string = "foo"
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

	var (
		singleCommandString                            string = command3Name
		singleCommandDescriptionString                 string = fmt.Sprintf("%s\n%ssome description", command3Name, dblIndent)
		singleCommandLongDescriptionString             string = fmt.Sprintf("%s\n%s%s", command3Name, dblIndent, longDescription)
		singleCommandSingleArgString                   string = fmt.Sprintf("%s %s", command3Name, singleArgString)
		singleCommandSingleArgDescriptionString        string = fmt.Sprintf("%s %s\n%ssome description", command3Name, singleArgString, dblIndent)
		singleCommandSingleArgLongDescriptionString    string = fmt.Sprintf("%s %s\n%s%s", command3Name, singleArgString, dblIndent, longDescription)
		singleCommandMultipleArgsString                string = fmt.Sprintf("%s %s", command3Name, multipleArgsString)
		singleCommandMultipleArgsDescriptionString     string = fmt.Sprintf("%s %s\n%ssome description", command3Name, multipleArgsString, dblIndent)
		singleCommandMultipleArgsLongDescriptionString string = fmt.Sprintf("%s %s\n%s%s", command3Name, multipleArgsString, dblIndent, longDescription)
		multipleCommandsString                         string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			command1Name,
			Indent+command2Name,
			Indent+command3Name,
		)
		multipleCommandsDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s\n%ssome description", command1Name, dblIndent),
			Indent+fmt.Sprintf("%s\n%ssome description", command2Name, dblIndent),
			Indent+fmt.Sprintf("%s\n%ssome description", command3Name, dblIndent),
		)
		multipleCommandsLongDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s\n%s%s", command1Name, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s\n%s%s", command2Name, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s\n%s%s", command3Name, dblIndent, longDescription),
		)
		multipleCommandsSingleArgString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s", command1Name, singleArgString),
			Indent+fmt.Sprintf("%s %s", command2Name, singleArgString),
			Indent+fmt.Sprintf("%s %s", command3Name, singleArgString),
		)
		multipleCommandsSingleArgDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%ssome description", command1Name, singleArgString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", command2Name, singleArgString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", command3Name, singleArgString, dblIndent),
		)
		multipleCommandsSingleArgLongDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%s%s", command1Name, singleArgString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", command2Name, singleArgString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", command3Name, singleArgString, dblIndent, longDescription),
		)
		multipleCommandsMultipleArgsString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s", command1Name, multipleArgsString),
			Indent+fmt.Sprintf("%s %s", command2Name, multipleArgsString),
			Indent+fmt.Sprintf("%s %s", command3Name, multipleArgsString),
		)
		multipleCommandsMultipleArgsDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%ssome description", command1Name, multipleArgsString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", command2Name, multipleArgsString, dblIndent),
			Indent+fmt.Sprintf("%s %s\n%ssome description", command3Name, multipleArgsString, dblIndent),
		)
		multipleCommandsMultipleArgsLongDescriptionString string = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			fmt.Sprintf("%s %s\n%s%s", command1Name, multipleArgsString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", command2Name, multipleArgsString, dblIndent, longDescription),
			Indent+fmt.Sprintf("%s %s\n%s%s", command3Name, multipleArgsString, dblIndent, longDescription),
		)
	)

	t.Run("baseline", usageUsageTester{
		oUsage: fmt.Sprintf("Usage:\n%sbase\n", Indent),
	}.assertString())
	t.Run("single global arg", usageUsageTester{
		oUsage: fmt.Sprintf("Usage:\n%ssingle-global-arg %s\n", Indent, singleArgString),
	}.assertString())
	t.Run("multiple global args", usageUsageTester{
		oUsage: fmt.Sprintf("Usage:\n%smultiple-global-args %s\n", Indent, multipleArgsString),
	}.assertString())
	t.Run("single option", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionString,
		),
	}.assertString())
	t.Run("single option description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionDescriptionString,
		),
	}.assertString())
	t.Run("single option long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionLongDescriptionString,
		),
	}.assertString())
	t.Run("single option single arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionSingleArgString,
		),
	}.assertString())
	t.Run("single option single arg description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("single option single arg long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleArgsString,
		),
	}.assertString())
	t.Run("single option multiple args description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple args long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-aliases %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesString,
		),
	}.assertString())
	t.Run("single option multiple aliases description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-aliases-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-aliases-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases single arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-aliases-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesSingleArgString,
		),
	}.assertString())
	t.Run("single option multiple aliases single arg description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-aliases-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases single arg long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-aliases-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases multiple args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-aliases-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesMultipleArgsString,
		),
	}.assertString())
	t.Run("single option multiple aliases multiple args description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-aliases-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases multiple args long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-aliases-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsString,
		),
	}.assertString())
	t.Run("multiple options description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsDescriptionString,
		),
	}.assertString())
	t.Run("multiple options long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsSingleArgString,
		),
	}.assertString())
	t.Run("multiple options single arg description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single arg long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleArgsString,
		),
	}.assertString())
	t.Run("multiple options multiple args description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple args long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-aliases %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-aliases-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-aliases-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases single arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-aliases-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesSingleArgString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases single arg description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-aliases-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases single arg long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-aliases-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases multiple args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-aliases-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesMultipleArgsString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases multiple args description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-aliases-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases multiple args long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-aliases-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("single command", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-command %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			fmt.Sprintf(summaryExtensionFormat, "single-command"),
			Indent+singleCommandString,
		),
	}.assertString())
	t.Run("single command description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-command-description %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			fmt.Sprintf(summaryExtensionFormat, "single-command-description"),
			Indent+singleCommandDescriptionString,
		),
	}.assertString())
	t.Run("single command long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-command-long-description %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			fmt.Sprintf(summaryExtensionFormat, "single-command-long-description"),
			Indent+singleCommandLongDescriptionString,
		),
	}.assertString())
	t.Run("single command single arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-command-single-arg %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-command-single-arg"),
			Indent+singleCommandSingleArgString,
		),
	}.assertString())
	t.Run("single command single arg description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-command-single-arg-description %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-command-single-arg-description"),
			Indent+singleCommandSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("single command single arg long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-command-single-arg-long-description %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-command-single-arg-long-description"),
			Indent+singleCommandSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single command multiple args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-command-multiple-args %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-command-multiple-args"),
			Indent+singleCommandMultipleArgsString,
		),
	}.assertString())
	t.Run("single command multiple args description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-command-multiple-args-description %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-command-multiple-args-description"),
			Indent+singleCommandMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("single command multiple args long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-command-multiple-args-long-description %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-command-multiple-args-long-description"),
			Indent+singleCommandMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple commands", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-commands %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-commands"),
			Indent+multipleCommandsString,
		),
	}.assertString())
	t.Run("multiple commands description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-commands-description %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-commands-description"),
			Indent+multipleCommandsDescriptionString,
		),
	}.assertString())
	t.Run("multiple commands long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-commands-long-description %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-commands-long-description"),
			Indent+multipleCommandsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple commands single arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-commands-single-arg %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-commands-single-arg"),
			Indent+multipleCommandsSingleArgString,
		),
	}.assertString())
	t.Run("multiple commands single arg description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-commands-single-arg-description %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-commands-single-arg-description"),
			Indent+multipleCommandsSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("multiple commands single arg long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-commands-single-arg-long-description %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-commands-single-arg-long-description"),
			Indent+multipleCommandsSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple commands multiple args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-commands-multiple-args %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-commands-multiple-args"),
			Indent+multipleCommandsMultipleArgsString,
		),
	}.assertString())
	t.Run("multiple commands multiple args description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-commands-multiple-args-description %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-commands-multiple-args-description"),
			Indent+multipleCommandsMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("multiple commands multiple args long description", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-commands-multiple-args-long-description %s %s\n\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-commands-multiple-args-long-description"),
			Indent+multipleCommandsMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("single option single global arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-single-global-arg %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			singleArgString,
			Indent+singleOptionLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple global args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-global-args %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			multipleArgsString,
			Indent+singleOptionLongDescriptionString,
		),
	}.assertString())
	t.Run("single option single command", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-single-command %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			fmt.Sprintf(summaryExtensionFormat, "single-option-single-command"),
			Indent+singleOptionLongDescriptionString,
			Indent+singleCommandLongDescriptionString,
		),
	}.assertString())
	t.Run("single option single command single command arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-single-command-single-command-arg %s %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-option-single-command-single-command-arg"),
			Indent+singleOptionLongDescriptionString,
			Indent+singleCommandSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option single command multiple command args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-single-command-multiple-command-args %s %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-option-single-command-multiple-command-args"),
			Indent+singleOptionLongDescriptionString,
			Indent+singleCommandMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple commands", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-commands %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			fmt.Sprintf(summaryExtensionFormat, "single-option-multiple-commands"),
			Indent+singleOptionLongDescriptionString,
			Indent+multipleCommandsLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple commands single command arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-commands-single-command-arg %s %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-option-multiple-commands-single-command-arg"),
			Indent+singleOptionLongDescriptionString,
			Indent+multipleCommandsSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple commands multiple command args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%ssingle-option-multiple-commands-multiple-command-args %s %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "single-option-multiple-commands-multiple-command-args"),
			Indent+singleOptionLongDescriptionString,
			Indent+multipleCommandsMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single global arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-single-global-arg %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			singleArgString,
			Indent+multipleOptionsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple global args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-global-args %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			multipleArgsString,
			Indent+multipleOptionsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single command", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-single-command %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-options-single-command"),
			Indent+multipleOptionsLongDescriptionString,
			Indent+singleCommandLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single command single command arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-single-command-single-command-arg %s %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-options-single-command-single-command-arg"),
			Indent+multipleOptionsLongDescriptionString,
			Indent+singleCommandSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single command multiple command args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-single-command-multiple-command-args %s %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-options-single-command-multiple-command-args"),
			Indent+multipleOptionsLongDescriptionString,
			Indent+singleCommandMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple commands", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-commands %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-options-multiple-commands"),
			Indent+multipleOptionsLongDescriptionString,
			Indent+multipleCommandsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple commands single command arg", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-commands-single-command-arg %s %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-options-multiple-commands-single-command-arg"),
			Indent+multipleOptionsLongDescriptionString,
			Indent+multipleCommandsSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple commands multiple command args", usageUsageTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%smultiple-options-multiple-commands-multiple-command-args %s %s %s\n\n%s\n\nOptions:\n%s\n\nCommands:\n%s\n",
			Indent,
			summaryCommandsString,
			summaryOptionsString,
			summaryArgsString,
			fmt.Sprintf(summaryExtensionFormat, "multiple-options-multiple-commands-multiple-command-args"),
			Indent+multipleOptionsLongDescriptionString,
			Indent+multipleCommandsMultipleArgsLongDescriptionString,
		),
	}.assertString())
}

func TestUsageLookup(t *testing.T) {
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

	t.Run("baseline", usageLookupTester{
		oUsage: fmt.Sprintf("Usage:\n%scalling-parent base\n", Indent),
	}.assertString())
	t.Run("single arg", usageLookupTester{
		oUsage: fmt.Sprintf("Usage:\n%scalling-parent single-arg %s\n", Indent, singleArgString),
	}.assertString())
	t.Run("multiple args", usageLookupTester{
		oUsage: fmt.Sprintf("Usage:\n%scalling-parent multiple-args %s\n", Indent, multipleArgsString),
	}.assertString())
	t.Run("single option", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionString,
		),
	}.assertString())
	t.Run("single option description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionDescriptionString,
		),
	}.assertString())
	t.Run("single option long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionLongDescriptionString,
		),
	}.assertString())
	t.Run("single option single arg", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionSingleArgString,
		),
	}.assertString())
	t.Run("single option single arg description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("single option single arg long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple args", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleArgsString,
		),
	}.assertString())
	t.Run("single option multiple args description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple args long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-aliases %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesString,
		),
	}.assertString())
	t.Run("single option multiple aliases description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-aliases-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-aliases-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases single arg", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-aliases-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesSingleArgString,
		),
	}.assertString())
	t.Run("single option multiple aliases single arg description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-aliases-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases single arg long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-aliases-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases multiple args", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-aliases-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesMultipleArgsString,
		),
	}.assertString())
	t.Run("single option multiple aliases multiple args description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-aliases-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple aliases multiple args long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-aliases-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+singleOptionMultipleAliasesMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsString,
		),
	}.assertString())
	t.Run("multiple options description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsDescriptionString,
		),
	}.assertString())
	t.Run("multiple options long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single arg", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsSingleArgString,
		),
	}.assertString())
	t.Run("multiple options single arg description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single arg long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple args", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleArgsString,
		),
	}.assertString())
	t.Run("multiple options multiple args description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple args long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-aliases %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-aliases-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-aliases-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases single arg", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-aliases-single-arg %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesSingleArgString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases single arg description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-aliases-single-arg-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesSingleArgDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases single arg long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-aliases-single-arg-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases multiple args", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-aliases-multiple-args %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesMultipleArgsString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases multiple args description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-aliases-multiple-args-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesMultipleArgsDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple aliases multiple args long description", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-aliases-multiple-args-long-description %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			Indent+multipleOptionsMultipleAliasesMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("single option single global arg", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-single-global-arg %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			singleArgString,
			Indent+singleOptionSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("single option multiple global args", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent single-option-multiple-global-args %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			multipleArgsString,
			Indent+singleOptionMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options single global arg", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-single-global-arg %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			singleArgString,
			Indent+multipleOptionsSingleArgLongDescriptionString,
		),
	}.assertString())
	t.Run("multiple options multiple global args", usageLookupTester{
		oUsage: fmt.Sprintf(
			"Usage:\n%scalling-parent multiple-options-multiple-global-args %s %s\n\nOptions:\n%s\n",
			Indent,
			summaryOptionsString,
			multipleArgsString,
			Indent+multipleOptionsMultipleArgsLongDescriptionString,
		),
	}.assertString())
	t.Run("entry does not exist", usageLookupTester{}.assertEmptyString())
}
