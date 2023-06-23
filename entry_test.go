package usage

import (
	"errors"
	"fmt"
	"sort"
	"testing"
)

type entryArgsTester struct {
	oArgs []string
}

func (tester entryArgsTester) assertArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{args: tester.oArgs}
		got := sampleEntry.Args()
		assertArgs(t, got, tester.oArgs)
	}
}

type entryOptionsTester struct {
	oOptions []Option
}

func (tester entryOptionsTester) assertEntryOptions() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{options: tester.oOptions}
		got := sampleEntry.Options()
		assertOptions(t, got, tester.oOptions)
	}
}

type entryEntriesTester struct {
	oEntries []Entry
}

func (tester entryEntriesTester) assertEntries() func(*testing.T) {
	return func(t *testing.T) {
		sort.Slice(tester.oEntries, func(i, j int) bool {
			return tester.oEntries[i].name < tester.oEntries[j].name
		})
		sampleEntry := Entry{children: make(map[string]*Entry)}
		for i, e := range tester.oEntries {
			sampleEntry.children[e.name] = &tester.oEntries[i]
		}
		got := sampleEntry.Entries()
		assertEntries(t, got, tester.oEntries)
	}
}

type entryNameTester struct {
	oName string
}

func (tester entryNameTester) assertName() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{name: tester.oName}
		got := sampleEntry.Name()
		assertName(t, got, tester.oName)
	}
}

type entryAddArgTester struct {
	iArg string
	oErr error
}

func (tester entryAddArgTester) assertArgs() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		args := make([]string, 0, iterations)
		sampleEntry := Entry{args: make([]string, 0)}
		for i := 1; i <= iterations; i++ {
			gotErr := sampleEntry.AddArg(tester.iArg)
			assertNilError(t, gotErr)
			args = append(args, tester.iArg)
		}
		assertArgs(t, sampleEntry.args, args)
	}
}

func (tester entryAddArgTester) assertEmptyArgStringError() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{args: make([]string, 0)}
		got := sampleEntry.AddArg(tester.iArg)
		assertEmptyArgStringError(t, got, tester.oErr)
	}
}

func (tester entryAddArgTester) assertExistingEntriesError() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := &Entry{
			children: map[string]*Entry{
				"foo": {name: "foo"},
			},
			args: make([]string, 0),
		}
		got := sampleEntry.AddArg(tester.iArg)
		assertExistingEntriesError(t, got, tester.oErr)
	}
}

type entryAddOptionTester struct {
	iOption *Option
	oErr    error
}

func (tester entryAddOptionTester) assertOptions() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		options := make([]Option, 0, iterations)
		sampleEntry := &Entry{options: make([]Option, 0)}
		for i := 1; i <= iterations; i++ {
			gotErr := sampleEntry.AddOption(tester.iOption)
			assertNilError(t, gotErr)
			options = append(options, *tester.iOption)
		}
		assertOptions(t, sampleEntry.options, options)
	}
}

func (tester entryAddOptionTester) assertNoOptionError() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := &Entry{options: make([]Option, 0)}
		got := sampleEntry.AddOption(tester.iOption)
		assertNoOptionError(t, got, tester.oErr)
	}
}

func (tester entryAddOptionTester) assertNoAliasesError() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := &Entry{options: make([]Option, 0)}
		got := sampleEntry.AddOption(tester.iOption)
		assertNoAliasesError(t, got, tester.oErr)
	}
}

func (tester entryAddOptionTester) assertEmptyAliasStringError() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := &Entry{options: make([]Option, 0)}
		got := sampleEntry.AddOption(tester.iOption)
		assertEmptyAliasStringError(t, got, tester.oErr)
	}
}

type entryAddEntryTester struct {
	iEntry *Entry
	oErr   error
}

func (tester entryAddEntryTester) assertChildren() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		entries := make(map[string]*Entry)
		sampleEntry := &Entry{children: make(map[string]*Entry)}
		for i := 1; i <= iterations; i++ {
			child := *tester.iEntry
			child.name += fmt.Sprintf("-%d", i)
			gotErr := sampleEntry.AddEntry(&child)
			assertNilError(t, gotErr)
			assertParent(t, child.parent, sampleEntry)
			entries[child.name] = &child
		}
		assertChildren(t, sampleEntry.children, entries)
	}
}

func (tester entryAddEntryTester) assertNoEntryError() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := &Entry{children: make(map[string]*Entry)}
		got := sampleEntry.AddEntry(tester.iEntry)
		assertNoEntryError(t, got, tester.oErr)
	}
}

func (tester entryAddEntryTester) assertEmptyNameStringError() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := &Entry{children: make(map[string]*Entry)}
		got := sampleEntry.AddEntry(tester.iEntry)
		assertEmptyNameStringError(t, got, tester.oErr)
	}
}

func (tester entryAddEntryTester) assertExistingArgsError() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := &Entry{
			children: make(map[string]*Entry),
			args:     []string{"foo"},
		}
		got := sampleEntry.AddEntry(tester.iEntry)
		assertExistingArgsError(t, got, tester.oErr)
	}
}

type entrySetNameTester struct {
	iName string
	oErr  error
}

func (tester entrySetNameTester) assertName() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{name: tester.iName}
		gotErr := sampleEntry.SetName(tester.iName)
		assertNilError(t, gotErr)
		assertName(t, sampleEntry.name, tester.iName)
	}
}

func (tester entrySetNameTester) assertEmptyNameStringError() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := Entry{name: "foo"}
		got := sampleEntry.SetName(tester.iName)
		assertEmptyNameStringError(t, got, tester.oErr)
	}
}

type newEntryTester struct {
	iName        string
	iDescription string
	oEntry       *Entry
	oErr         error
}

func (tester newEntryTester) assertDefaultEntry() func(*testing.T) {
	return func(t *testing.T) {
		got, gotErr := NewEntry(tester.iName, tester.iDescription)
		assertNilError(t, gotErr)
		assertDefaultEntry(t, got, tester.oEntry)
	}
}

func (tester newEntryTester) assertEmptyNameStringError() func(*testing.T) {
	return func(t *testing.T) {
		gotEntry, got := NewEntry(tester.iName, tester.iDescription)
		assertNilEntry(t, gotEntry)
		assertError(t, got, tester.oErr)
	}
}

func TestEntryArgs(t *testing.T) {
	t.Run("baseline", entryArgsTester{
		oArgs: []string{"foo"},
	}.assertArgs())
	t.Run("multiple args", entryArgsTester{
		oArgs: []string{"foo", "bar", "baz"},
	}.assertArgs())
	t.Run("no args", entryArgsTester{
		oArgs: make([]string, 0),
	}.assertArgs())
}

func TestEntryOptions(t *testing.T) {
	t.Run("baseline", entryOptionsTester{
		oOptions: []Option{{
			Description: "foo",
			Tmpl:        "foo",
			aliases:     []string{"foo"},
			args:        []string{"foo"},
		}},
	}.assertEntryOptions())
	t.Run("multiple options", entryOptionsTester{
		oOptions: []Option{
			{
				Description: "foo",
				Tmpl:        "foo",
				aliases:     []string{"foo"},
				args:        []string{"foo"},
			},
			{
				Description: "bar",
				Tmpl:        "bar",
				aliases:     []string{"bar"},
				args:        []string{"bar"},
			},
			{
				Description: "baz",
				Tmpl:        "baz",
				aliases:     []string{"baz"},
				args:        []string{"baz"},
			},
		},
	}.assertEntryOptions())
	t.Run("no options", entryOptionsTester{
		oOptions: make([]Option, 0),
	}.assertEntryOptions())
}

func TestEntryEntries(t *testing.T) {
	t.Run("baseline", entryEntriesTester{
		oEntries: []Entry{{
			Description: "foo",
			Tmpl:        "foo",
			name:        "foo",
			options: []Option{{
				aliases:     []string{"foo"},
				Description: "foo",
				args:        []string{"foo"},
			}},
			args: []string{"foo"},
		}},
	}.assertEntries())
	t.Run("multiple entries", entryEntriesTester{
		oEntries: []Entry{
			{
				Description: "foo",
				Tmpl:        "foo",
				name:        "foo",
				options: []Option{{
					aliases:     []string{"foo"},
					Description: "foo",
					args:        []string{"foo"},
				}},
				args: []string{"foo"},
			},
			{
				Description: "bar",
				Tmpl:        "bar",
				name:        "bar",
				options: []Option{{
					aliases:     []string{"bar"},
					Description: "bar",
					args:        []string{"bar"},
				}},
				args: []string{"bar"},
			},
			{
				Description: "baz",
				Tmpl:        "baz",
				name:        "baz",
				options: []Option{{
					aliases:     []string{"baz"},
					Description: "baz",
					args:        []string{"baz"},
				}},
				args: []string{"baz"},
			},
		},
	}.assertEntries())
	t.Run("no entries", entryEntriesTester{
		oEntries: make([]Entry, 0),
	}.assertEntries())
}

func TestEntryName(t *testing.T) {
	t.Run("baseline", entryNameTester{
		oName: "foo",
	}.assertName())
}

func TestEntryAddArg(t *testing.T) {
	t.Run("baseline", entryAddArgTester{
		iArg: "foo",
	}.assertArgs())
	t.Run("empty arg string", entryAddArgTester{
		oErr: errors.New("usage: arg string must not be empty"),
	}.assertEmptyArgStringError())
	t.Run("existing entries", entryAddArgTester{
		oErr: errors.New("usage: cannot add arg with child entries present"),
	}.assertExistingEntriesError())
}

func TestEntryAddOption(t *testing.T) {
	t.Run("baseline", entryAddOptionTester{
		iOption: &Option{
			Description: "foo",
			Tmpl:        "foo",
			aliases:     []string{"foo"},
			args:        []string{"foo"},
		},
	}.assertOptions())
	t.Run("nil option", entryAddOptionTester{
		oErr: errors.New("usage: no option provided"),
	}.assertNoOptionError())
	t.Run("nil aliases", entryAddOptionTester{
		iOption: &Option{args: []string{"foo"}},
		oErr:    errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesError())
	t.Run("no aliases", entryAddOptionTester{
		iOption: &Option{aliases: []string{}},
		oErr:    errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesError())
	t.Run("single empty alias string", entryAddOptionTester{
		iOption: &Option{aliases: []string{""}},
		oErr:    errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringError())
	t.Run("multiple empty alias strings", entryAddOptionTester{
		iOption: &Option{aliases: []string{"foo", "", "bar", ""}},
		oErr:    errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringError())
}

func TestAddEntry(t *testing.T) {
	t.Run("baseline", entryAddEntryTester{
		iEntry: &Entry{name: "foo"},
	}.assertChildren())
	t.Run("nil entry", entryAddEntryTester{
		oErr: errors.New("usage: no entry provided"),
	}.assertNoEntryError())
	t.Run("empty name string", entryAddEntryTester{
		iEntry: &Entry{},
		oErr:   errors.New("usage: name string must not be empty"),
	}.assertEmptyNameStringError())
	t.Run("existing args", entryAddEntryTester{
		iEntry: &Entry{name: "foo"},
		oErr:   errors.New("usage: cannot add child entry with args present"),
	}.assertExistingArgsError())
}

func TestEntrySetName(t *testing.T) {
	t.Run("baseline", entrySetNameTester{
		iName: "foo",
	}.assertName())
	t.Run("empty name string", entrySetNameTester{
		oErr: errors.New("usage: name string must not be empty"),
	}.assertEmptyNameStringError())
}

func TestNewEntry(t *testing.T) {
	t.Run("baseline", newEntryTester{
		iName:        "foo",
		iDescription: "foo",
		oEntry:       &Entry{name: "foo", Description: "foo"},
	}.assertDefaultEntry())
	t.Run("empty description string", newEntryTester{
		iName:  "foo",
		oEntry: &Entry{name: "foo"},
	}.assertDefaultEntry())
	t.Run("empty name string", newEntryTester{
		iDescription: "foo",
		oErr:         errors.New("usage: name string must not be empty"),
	}.assertEmptyNameStringError())
}
