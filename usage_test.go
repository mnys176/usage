package usage

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"text/template"
)

type initTester struct {
	iName string
	oErr  error
}

func (tester initTester) assertEntry() func(*testing.T) {
	return func(t *testing.T) {
		gotErr := Init(tester.iName)
		assertNilError(t, gotErr)
		assertDefaultEntry(t, global, stringToEntry(tester.iName))
		global = nil
	}
}

func (tester initTester) assertEmptyNameStringError() func(*testing.T) {
	return func(t *testing.T) {
		got := Init(tester.iName)
		assertEmptyNameStringError(t, got, tester.oErr)
		assertNilEntry(t, global)
	}
}

type argsTester struct {
	oArgs  []string
	oPanic error
}

func (tester argsTester) assertArgs() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{args: tester.oArgs}
		got := Args()
		assertArgs(t, got, tester.oArgs)
		global = nil
	}
}

func (tester argsTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		Args()
		assertNilEntry(t, global)
	}
}

type optionsTester struct {
	oOptions []Option
	oPanic   error
}

func (tester optionsTester) assertOptions() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{options: tester.oOptions}
		got := Options()
		assertOptions(t, got, tester.oOptions)
		global = nil
	}
}

func (tester optionsTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		Options()
		assertNilEntry(t, global)
	}
}

type entriesTester struct {
	oEntries []Entry
	oPanic   error
}

func (tester entriesTester) assertEntries() func(*testing.T) {
	return func(t *testing.T) {
		sort.Slice(tester.oEntries, func(i, j int) bool {
			return tester.oEntries[i].name < tester.oEntries[j].name
		})
		global = &Entry{children: make(map[string]*Entry)}
		for i, e := range tester.oEntries {
			global.children[e.name] = &tester.oEntries[i]
		}
		got := Entries()
		assertEntries(t, got, tester.oEntries)
		global = nil
	}
}

func (tester entriesTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		Entries()
		assertNilEntry(t, global)
	}
}

type addArgTester struct {
	iArg   string
	oErr   error
	oPanic error
}

func (tester addArgTester) assertArgs() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		args := make([]string, 0, iterations)
		global = &Entry{args: make([]string, 0)}
		for i := 1; i <= iterations; i++ {
			gotErr := AddArg(tester.iArg)
			assertNilError(t, gotErr)
			args = append(args, tester.iArg)
		}
		assertArgs(t, global.args, args)
		global = nil
	}
}

func (tester addArgTester) assertEmptyArgStringError() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{args: make([]string, 0)}
		got := AddArg(tester.iArg)
		assertEmptyArgStringError(t, got, tester.oErr)
		global = nil
	}
}

func (tester addArgTester) assertExistingEntriesError() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{
			children: map[string]*Entry{
				"foo": {name: "foo"},
			},
			args: make([]string, 0),
		}
		got := AddArg(tester.iArg)
		assertExistingEntriesError(t, got, tester.oErr)
		global = nil
	}
}

func (tester addArgTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		AddArg(tester.iArg)
		assertNilEntry(t, global)
	}
}

type addOptionTester struct {
	iOption *Option
	oErr    error
	oPanic  error
}

func (tester addOptionTester) assertOptions() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		options := make([]Option, 0, iterations)
		global = &Entry{options: make([]Option, 0)}
		for i := 1; i <= iterations; i++ {
			gotErr := AddOption(tester.iOption)
			assertNilError(t, gotErr)
			options = append(options, *tester.iOption)
		}
		assertOptions(t, global.options, options)
		global = nil
	}
}

func (tester addOptionTester) assertNoOptionError() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{options: make([]Option, 0)}
		got := AddOption(tester.iOption)
		assertNoOptionError(t, got, tester.oErr)
		global = nil
	}
}

func (tester addOptionTester) assertNoAliasesError() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{options: make([]Option, 0)}
		got := AddOption(tester.iOption)
		assertNoAliasesError(t, got, tester.oErr)
		global = nil
	}
}

func (tester addOptionTester) assertEmptyAliasStringError() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{options: make([]Option, 0)}
		got := AddOption(tester.iOption)
		assertEmptyAliasStringError(t, got, tester.oErr)
		global = nil
	}
}

func (tester addOptionTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		AddOption(tester.iOption)
		assertNilEntry(t, global)
	}
}

type addEntryTester struct {
	iEntry *Entry
	oErr   error
	oPanic error
}

func (tester addEntryTester) assertChildren() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		entries := make(map[string]*Entry)
		global = &Entry{children: make(map[string]*Entry)}
		for i := 1; i <= iterations; i++ {
			child := *tester.iEntry
			child.name += fmt.Sprintf("-%d", i)
			gotErr := AddEntry(&child)
			assertNilError(t, gotErr)
			assertParent(t, child.parent, global)
			entries[child.name] = &child
		}
		assertChildren(t, global.children, entries)
		global = nil
	}
}

func (tester addEntryTester) assertNoEntryError() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{children: make(map[string]*Entry)}
		got := AddEntry(tester.iEntry)
		assertNoEntryError(t, got, tester.oErr)
		global = nil
	}
}

func (tester addEntryTester) assertEmptyNameStringError() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{children: make(map[string]*Entry)}
		got := AddEntry(tester.iEntry)
		assertEmptyNameStringError(t, got, tester.oErr)
		global = nil
	}
}

func (tester addEntryTester) assertExistingArgsError() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{
			children: make(map[string]*Entry),
			args:     []string{"foo"},
		}
		got := AddEntry(tester.iEntry)
		assertExistingArgsError(t, got, tester.oErr)
		global = nil
	}
}

func (tester addEntryTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		AddEntry(tester.iEntry)
		assertNilEntry(t, global)
	}
}

type setNameTester struct {
	iName  string
	oErr   error
	oPanic error
}

func (tester setNameTester) assertName() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{name: tester.iName}
		gotErr := SetName(tester.iName)
		assertNilError(t, gotErr)
		assertName(t, global.name, tester.iName)
		global = nil
	}
}

func (tester setNameTester) assertEmptyNameStringError() func(*testing.T) {
	return func(t *testing.T) {
		global = &Entry{name: "foo"}
		got := SetName(tester.iName)
		assertEmptyNameStringError(t, got, tester.oErr)
		global = nil
	}
}

func (tester setNameTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		SetName(tester.iName)
		assertNilEntry(t, global)
	}
}

type usageTester struct {
	oUsage string
	oPanic error
}

func (tester usageTester) assertUsage() func(*testing.T) {
	return func(t *testing.T) {
		global = stringToEntry(tester.oUsage)
		got := Usage()
		assertUsage(t, got, tester.oUsage)
		global = nil
	}
}

func (tester usageTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		Usage()
		assertNilEntry(t, global)
	}
}

type lookupTester struct {
	iLookup string
	oUsage  string
	oPanic  error
}

func (tester lookupTester) assertUsage() func(*testing.T) {
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
		global = &Entry{
			name:     "base",
			children: make(map[string]*Entry),
			tmpl:     template.Must(template.New("").Funcs(fn).Parse(rawTmpl)),
		}
		ptr := global
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
		got := Lookup(tester.iLookup)
		assertUsage(t, got, tester.oUsage)
		global = nil
	}
}

func (tester lookupTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		Lookup(tester.iLookup)
		assertNilEntry(t, global)
	}
}

type setEntryTemplateTester struct {
	iTemplate *template.Template
	oPanic    error
}

func (tester setEntryTemplateTester) assertTemplate() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		global = &Entry{name: "base", children: make(map[string]*Entry)}
		ptr := global
		for i := 1; i <= iterations; i++ {
			entry := Entry{
				name:     fmt.Sprintf("level-%d", i),
				children: make(map[string]*Entry),
				parent:   ptr,
			}
			ptr.children[entry.name] = &entry
			ptr = &entry
		}
		SetEntryTemplate(tester.iTemplate)
		global = nil
	}
}

func (tester setEntryTemplateTester) assertUninitializedErrorPanic() func(*testing.T) {
	return func(t *testing.T) {
		defer assertUninitializedPanic(t, tester.oPanic)
		SetEntryTemplate(tester.iTemplate)
		assertNilEntry(t, global)
	}
}

func TestInit(t *testing.T) {
	t.Run("baseline", initTester{
		iName: "foo",
	}.assertEntry())
	t.Run("empty name string", initTester{
		oErr: errors.New("usage: name string must not be empty"),
	}.assertEmptyNameStringError())
}

func TestArgs(t *testing.T) {
	t.Run("baseline", argsTester{
		oArgs: []string{"foo"},
	}.assertArgs())
	t.Run("uninitialized", argsTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}

func TestOptions(t *testing.T) {
	t.Run("baseline", optionsTester{
		oOptions: []Option{{
			Description: "foo",
			aliases:     []string{"foo"},
			args:        []string{"foo"},
		}},
	}.assertOptions())
	t.Run("uninitialized", optionsTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}

func TestEntries(t *testing.T) {
	t.Run("baseline", entriesTester{
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
	t.Run("uninitialized", entriesTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}

func TestAddArg(t *testing.T) {
	t.Run("baseline", addArgTester{
		iArg: "foo",
	}.assertArgs())
	t.Run("empty arg string", addArgTester{
		oErr: errors.New("usage: arg string must not be empty"),
	}.assertEmptyArgStringError())
	t.Run("existing entries", addArgTester{
		oErr: errors.New("usage: cannot add arg with child entries present"),
	}.assertExistingEntriesError())
	t.Run("uninitialized", addArgTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}

func TestAddOption(t *testing.T) {
	t.Run("baseline", addOptionTester{
		iOption: &Option{
			Description: "foo",
			aliases:     []string{"foo"},
			args:        []string{"foo"},
		},
	}.assertOptions())
	t.Run("nil option", addOptionTester{
		oErr: errors.New("usage: no option provided"),
	}.assertNoOptionError())
	t.Run("nil aliases", addOptionTester{
		iOption: &Option{args: []string{"foo"}},
		oErr:    errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesError())
	t.Run("no aliases", addOptionTester{
		iOption: &Option{aliases: []string{}},
		oErr:    errors.New("usage: option must have at least one alias"),
	}.assertNoAliasesError())
	t.Run("single empty alias string", addOptionTester{
		iOption: &Option{aliases: []string{""}},
		oErr:    errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringError())
	t.Run("multiple empty alias strings", addOptionTester{
		iOption: &Option{aliases: []string{"foo", "", "bar", ""}},
		oErr:    errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringError())
	t.Run("uninitialized", addOptionTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}

func TestAddEntry(t *testing.T) {
	t.Run("baseline", addEntryTester{
		iEntry: &Entry{name: "foo"},
	}.assertChildren())
	t.Run("nil entry", addEntryTester{
		oErr: errors.New("usage: no entry provided"),
	}.assertNoEntryError())
	t.Run("empty name string", addEntryTester{
		iEntry: &Entry{},
		oErr:   errors.New("usage: name string must not be empty"),
	}.assertEmptyNameStringError())
	t.Run("existing args", addEntryTester{
		iEntry: &Entry{name: "foo"},
		oErr:   errors.New("usage: cannot add child entry with args present"),
	}.assertExistingArgsError())
	t.Run("uninitialized", addEntryTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}

func TestSetName(t *testing.T) {
	t.Run("baseline", setNameTester{
		iName: "foo",
	}.assertName())
	t.Run("empty name string", setNameTester{
		oErr: errors.New("usage: name string must not be empty"),
	}.assertEmptyNameStringError())
	t.Run("uninitialized", setNameTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}

func TestUsage(t *testing.T) {
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

	t.Run("baseline", usageTester{
		oUsage: "base",
	}.assertUsage())
	t.Run("ancestry", usageTester{
		oUsage: "parent:base",
	}.assertUsage())
	t.Run("args", usageTester{
		oUsage: "base <args>",
	}.assertUsage())
	t.Run("ancestry args", usageTester{
		oUsage: "parent:base <args>",
	}.assertUsage())
	t.Run("options", usageTester{
		oUsage: "base [options]",
	}.assertUsage())
	t.Run("ancestry options", usageTester{
		oUsage: "parent:base [options]",
	}.assertUsage())
	t.Run("options args", usageTester{
		oUsage: "base [options] <args>",
	}.assertUsage())
	t.Run("ancestry options args", usageTester{
		oUsage: "parent:base [options] <args>",
	}.assertUsage())
	t.Run("description", usageTester{
		oUsage: "base\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry description", usageTester{
		oUsage: "parent:base\n" + indent + description,
	}.assertUsage())
	t.Run("args description", usageTester{
		oUsage: "base <args>\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry args description", usageTester{
		oUsage: "parent:base <args>\n" + indent + description,
	}.assertUsage())
	t.Run("options description", usageTester{
		oUsage: "base [options]\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry options description", usageTester{
		oUsage: "parent:base [options]\n" + indent + description,
	}.assertUsage())
	t.Run("options args description", usageTester{
		oUsage: "base [options] <args>\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry options args description", usageTester{
		oUsage: "parent:base [options] <args>\n" + indent + description,
	}.assertUsage())
	t.Run("entries", usageTester{
		oUsage: "base <command>",
	}.assertUsage())
	t.Run("ancestry entries", usageTester{
		oUsage: "parent:base <command>",
	}.assertUsage())
	t.Run("options entries", usageTester{
		oUsage: "base [options] <command>",
	}.assertUsage())
	t.Run("ancestry options entries", usageTester{
		oUsage: "parent:base [options] <command>",
	}.assertUsage())
	t.Run("entries description", usageTester{
		oUsage: "base <command>\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry entries description", usageTester{
		oUsage: "parent:base <command>\n" + indent + description,
	}.assertUsage())
	t.Run("options entries description", usageTester{
		oUsage: "base [options] <command>\n" + indent + description,
	}.assertUsage())
	t.Run("ancestry options entries description", usageTester{
		oUsage: "parent:base [options] <command>\n" + indent + description,
	}.assertUsage())
	t.Run("uninitialized", usageTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}

func TestLookup(t *testing.T) {
	t.Run("baseline", lookupTester{
		iLookup: "level-1",
		oUsage:  "base:level-1 <command>",
	}.assertUsage())
	t.Run("leaf lookup", lookupTester{
		iLookup: "level-3",
		oUsage:  "base:level-1:level-2:level-3",
	}.assertUsage())
	t.Run("root lookup", lookupTester{
		iLookup: "base",
		oUsage:  "base <command>",
	}.assertUsage())
	t.Run("untracked entry", lookupTester{
		iLookup: "foo",
	}.assertUsage())
	t.Run("empty name string", lookupTester{}.assertUsage())
	t.Run("uninitialized", lookupTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}

func TestSetEntryTemplate(t *testing.T) {
	t.Run("baseline", setEntryTemplateTester{
		iTemplate: template.Must(template.New("").Parse("foo")),
	}.assertTemplate())
	t.Run("baseline", setEntryTemplateTester{
		iTemplate: template.Must(
			template.New("").
				Funcs(template.FuncMap{"fn": strings.ToUpper}).
				Parse(`{{fn "foo"}}`),
		),
	}.assertTemplate())
	t.Run("uninitialized", setEntryTemplateTester{
		oPanic: errors.New("usage: global usage not initialized"),
	}.assertUninitializedErrorPanic())
}
