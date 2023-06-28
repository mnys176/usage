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
