package usage

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"text/template"
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

type entryAncestryTester struct {
	oAncestry []string
}

func (tester entryAncestryTester) assertAncestry() func(*testing.T) {
	return func(t *testing.T) {
		entries := make([]Entry, len(tester.oAncestry))
		for i, name := range tester.oAncestry {
			entries[i].name = name
			if i < len(tester.oAncestry)-1 {
				entries[i].parent = &entries[i+1]
			}
		}
		sampleEntry := &entries[0]
		got := sampleEntry.Ancestry()
		assertAncestry(t, got, tester.oAncestry)
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

type entryUsageTester struct {
	oUsage string
}

func (tester entryUsageTester) assertUsage() func(*testing.T) {
	return func(t *testing.T) {
		sampleEntry := stringToEntry(tester.oUsage)
		got := sampleEntry.Usage()
		assertUsage(t, got, tester.oUsage)
	}
}

type entryLookupTester struct {
	iLookup string
	oUsage  string
}

func (tester entryLookupTester) assertUsage() func(*testing.T) {
	return func(t *testing.T) {
		const indent = "    "
		rawTmpl := fmt.Sprintf(`{{join (reverse .Ancestry) ":"}}{{if .Options}} [options]{{end}}{{if .Entries}} <command>{{end}}{{if .Args}} <args>{{end}}{{if .Description}}
%s{{with chop .Description 64}}{{join . "\n%s"}}{{end}}{{end}}`, indent, indent)
		fn := template.FuncMap{
			"join":    strings.Join,
			"reverse": reverseAncestryChain,
			"chop":    chopEssay,
		}
		iterations := 3
		sampleEntry := Entry{
			name:     "base",
			children: make(map[string]*Entry),
			tmpl:     template.Must(template.New("").Funcs(fn).Parse(rawTmpl)),
		}
		ptr := &sampleEntry
		for i := 1; i <= iterations; i++ {
			entry := Entry{
				name:     fmt.Sprintf("level-%d", i),
				children: make(map[string]*Entry),
				parent:   ptr,
				tmpl:     template.Must(template.New("").Funcs(fn).Parse(rawTmpl)),
			}
			ptr.children[entry.name] = &entry
			ptr = &entry
		}
		got := sampleEntry.Lookup(tester.iLookup)
		assertUsage(t, got, tester.oUsage)
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
			aliases:     []string{"foo"},
			args:        []string{"foo"},
		}},
	}.assertEntryOptions())
	t.Run("multiple options", entryOptionsTester{
		oOptions: []Option{
			{
				Description: "foo",
				aliases:     []string{"foo"},
				args:        []string{"foo"},
			},
			{
				Description: "bar",
				aliases:     []string{"bar"},
				args:        []string{"bar"},
			},
			{
				Description: "baz",
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

func TestEntryAncestry(t *testing.T) {
	t.Run("baseline", entryAncestryTester{
		oAncestry: []string{"foo", "bar", "baz"},
	}.assertAncestry())
	t.Run("root", entryAncestryTester{
		oAncestry: []string{"foo"},
	}.assertAncestry())
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

func TestEntryUsage(t *testing.T) {
	const (
		indent      = "    "
		description = "some very long description that will definitely push the limits\n" +
			indent + "of the screen size (it is very likely that this will cause the\n" +
			indent + "line break at 64 characters)\n" +
			indent + "\n" +
			indent + "here's another paragraph just in case with a very long word\n" +
			indent + "between these brackets > < that will not appear in the final\n" +
			indent + "output because it is longer than a line"
	)

	t.Run("baseline", entryUsageTester{
		oUsage: "base",
	}.assertUsage())
	t.Run("ancestry", entryUsageTester{
		oUsage: "parent:base",
	}.assertUsage())
	t.Run("args", entryUsageTester{
		oUsage: "base <args>",
	}.assertUsage())
	t.Run("ancestry args", entryUsageTester{
		oUsage: "parent:base <args>",
	}.assertUsage())
	t.Run("options", entryUsageTester{
		oUsage: "base [options]",
	}.assertUsage())
	t.Run("ancestry options", entryUsageTester{
		oUsage: "parent:base [options]",
	}.assertUsage())
	t.Run("options args", entryUsageTester{
		oUsage: "base [options] <args>",
	}.assertUsage())
	t.Run("ancestry options args", entryUsageTester{
		oUsage: "parent:base [options] <args>",
	}.assertUsage())
	t.Run("description", entryUsageTester{
		oUsage: "base\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry description", entryUsageTester{
		oUsage: "parent:base\n" + indent + description,
	}.assertUsage())
	t.Run("args description", entryUsageTester{
		oUsage: "base <args>\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry args description", entryUsageTester{
		oUsage: "parent:base <args>\n" + indent + description,
	}.assertUsage())
	t.Run("options description", entryUsageTester{
		oUsage: "base [options]\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry options description", entryUsageTester{
		oUsage: "parent:base [options]\n" + indent + description,
	}.assertUsage())
	t.Run("options args description", entryUsageTester{
		oUsage: "base [options] <args>\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry options args description", entryUsageTester{
		oUsage: "parent:base [options] <args>\n" + indent + description,
	}.assertUsage())
	t.Run("entries", entryUsageTester{
		oUsage: "base <command>",
	}.assertUsage())
	t.Run("ancestry entries", entryUsageTester{
		oUsage: "parent:base <command>",
	}.assertUsage())
	t.Run("options entries", entryUsageTester{
		oUsage: "base [options] <command>",
	}.assertUsage())
	t.Run("ancestry options entries", entryUsageTester{
		oUsage: "parent:base [options] <command>",
	}.assertUsage())
	t.Run("entries description", entryUsageTester{
		oUsage: "base <command>\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry entries description", entryUsageTester{
		oUsage: "parent:base <command>\n" + indent + description,
	}.assertUsage())
	t.Run("options entries description", entryUsageTester{
		oUsage: "base [options] <command>\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry options entries description", entryUsageTester{
		oUsage: "parent:base [options] <command>\n" + indent + description,
	}.assertUsage())
}

func TestEntryLookup(t *testing.T) {
	t.Run("baseline", entryLookupTester{
		iLookup: "level-1",
		oUsage:  "base:level-1 <command>",
	}.assertUsage())
	t.Run("leaf lookup", entryLookupTester{
		iLookup: "level-3",
		oUsage:  "base:level-1:level-2:level-3",
	}.assertUsage())
	t.Run("root lookup", entryLookupTester{
		iLookup: "base",
		oUsage:  "base <command>",
	}.assertUsage())
	t.Run("untracked entry", entryLookupTester{
		iLookup: "foo",
	}.assertUsage())
	t.Run("empty name string", entryLookupTester{}.assertUsage())
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
