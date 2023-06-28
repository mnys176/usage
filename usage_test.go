package usage

import (
	"errors"
	"sort"
	"testing"
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
