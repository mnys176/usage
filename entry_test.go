package usage

import (
	"errors"
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
}

func TestEntrySetName(t *testing.T) {
	t.Run("baseline", entrySetNameTester{
		iName: "foo",
	}.assertName())
	t.Run("empty name string", entrySetNameTester{
		oErr: errors.New("usage: name string must not be empty"),
	}.assertEmptyNameStringError())
}
