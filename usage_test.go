package usage

import (
	"errors"
	"sort"
	"testing"
)

func assertError(t *testing.T, got, want error) {
	if !errors.Is(got, want) {
		t.Errorf("got %q error but wanted %q", got, want)
	}
}

func assertEntryStruct(t *testing.T, got, want *entry) {
	if got.Name != want.Name {
		t.Errorf("name is %q but should be %q", got.Name, want.Name)
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
			return tester.oEntries[i].Name < tester.oEntries[j].Name
		})
		sampleUsage := usage{entries: make(map[string]entry)}
		for _, sampleEntry := range tester.oEntries {
			sampleUsage.entries[sampleEntry.Name] = sampleEntry
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
					Name:        "foo",
					Description: "foo",
					options: []option{{
						Aliases:     []string{"foo"},
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
			return sampleEntries[i].Name < sampleEntries[j].Name
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
			Aliases:     []string{"foo"},
			Description: "foo",
			args:        []string{"foo"},
		}},
	}.assertUsageOptions())
	t.Run("multiple options", usageOptionsTester{
		oOptions: []option{
			{
				Aliases:     []string{"foo"},
				Description: "foo",
				args:        []string{"foo"},
			},
			{
				Aliases:     []string{"bar"},
				Description: "bar",
				args:        []string{"bar"},
			},
			{
				Aliases:     []string{"baz"},
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
			Name:        "foo",
			Description: "foo",
			options: []option{{
				Aliases:     []string{"foo"},
				Description: "foo",
				args:        []string{"foo"},
			}},
			args: []string{"foo"},
		}},
	}.assertUsageEntries())
	t.Run("multiple entries", usageEntriesTester{
		oEntries: []entry{
			{
				Name:        "foo",
				Description: "foo",
				options: []option{{
					Aliases:     []string{"foo"},
					Description: "foo",
					args:        []string{"foo"},
				}},
				args: []string{"foo"},
			},
			{
				Name:        "bar",
				Description: "bar",
				options: []option{{
					Aliases:     []string{"bar"},
					Description: "bar",
					args:        []string{"bar"},
				}},
				args: []string{"bar"},
			},
			{
				Name:        "baz",
				Description: "baz",
				options: []option{{
					Aliases:     []string{"baz"},
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
			Aliases:     []string{"foo"},
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
			Name:        "foo",
			Description: "foo",
			options: []option{{
				Aliases:     []string{"foo"},
				Description: "foo",
				args:        []string{"foo"},
			}},
			args: []string{"foo"},
		},
	}.assertUsageEntries())
	t.Run("existing args", usageAddEntryTester{
		iEntry: &entry{
			Name:        "foo",
			Description: "foo",
			options: []option{{
				Aliases:     []string{"foo"},
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
