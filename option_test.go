package main

import (
	"errors"
	"testing"
)

func assertArgs(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d args returned but wanted %d", len(got), len(want))
	}
	for i, gotArg := range got {
		if gotArg != want[i] {
			t.Errorf("arg is %q but should be %q", gotArg, want[i])
		}
	}
}

func assertAliases(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d aliases returned but wanted %d", len(got), len(want))
	}
	for i, gotAlias := range got {
		if gotAlias != want[i] {
			t.Errorf("alias is %q but should be %q", gotAlias, want[i])
		}
	}
}

func assertError(t *testing.T, got, want error) {
	if !errors.Is(got, want) {
		t.Errorf("got %q error but wanted %q", got, want)
	}
}

type optionArgsTester struct {
	oArgs []string
}

func (tester optionArgsTester) assertOptionArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{args: tester.oArgs}
		got := sampleOption.Args()
		assertArgs(t, got, tester.oArgs)
	}
}

type optionAddArgTester struct {
	iArg string
	oErr error
}

func (tester optionAddArgTester) assertOptionArgs() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		args := make([]string, 0, iterations)
		sampleOption := Option{args: make([]string, 0)}
		for i := 1; i <= iterations; i++ {
			if gotErr := sampleOption.AddArg(tester.iArg); gotErr != nil {
				t.Errorf("got %q error but should be nil", gotErr)
			}
			args = append(args, tester.iArg)
		}
		assertArgs(t, sampleOption.args, args)
	}
}

func (tester optionAddArgTester) assertEmptyArgStringErr() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := Option{args: make([]string, 0)}
		got := sampleOption.AddArg(tester.iArg)
		if got == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		assertError(t, got, tester.oErr)
	}
}

type newOptionTester struct {
	iAliases     []string
	iDescription string
	oOption      *Option
	oErr         error
}

func (tester newOptionTester) assertOption() func(*testing.T) {
	return func(t *testing.T) {
		got, gotErr := NewOption(tester.iAliases, tester.iDescription)
		if gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		assertAliases(t, got.aliases, tester.oOption.aliases)
		if got.Description != tester.oOption.Description {
			t.Errorf("description is %q but should be %q", got.Description, tester.oOption.Description)
		}
		if got.args == nil || len(got.args) != 0 {
			t.Error("args not initialized to an empty slice")
		}
	}
}

func (tester newOptionTester) assertNoAliasErr() func(*testing.T) {
	return func(t *testing.T) {
		gotOption, got := NewOption(tester.iAliases, tester.iDescription)
		if gotOption != nil {
			t.Errorf("got %+v option but should be nil", gotOption)
		}
		if got == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		assertError(t, got, tester.oErr)
	}
}

func (tester newOptionTester) assertEmptyAliasStringErr() func(*testing.T) {
	return func(t *testing.T) {
		gotOption, got := NewOption(tester.iAliases, tester.iDescription)
		if gotOption != nil {
			t.Errorf("got %+v option but should be nil", gotOption)
		}
		if got == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		assertError(t, got, tester.oErr)
	}
}

func TestOptionArgs(t *testing.T) {
	t.Run("baseline", optionArgsTester{
		oArgs: []string{"foo"},
	}.assertOptionArgs())
	t.Run("multiple args", optionArgsTester{
		oArgs: []string{"foo", "bar", "baz"},
	}.assertOptionArgs())
	t.Run("no args", optionArgsTester{
		oArgs: make([]string, 0),
	}.assertOptionArgs())
}

func TestOptionAddArg(t *testing.T) {
	t.Run("baseline", optionAddArgTester{
		iArg: "foo",
	}.assertOptionArgs())
	t.Run("empty arg string", optionAddArgTester{
		oErr: errors.New("usage: arg string must not be empty"),
	}.assertEmptyArgStringErr())
}

func TestNewOption(t *testing.T) {
	t.Run("baseline", newOptionTester{
		iAliases:     []string{"foo"},
		iDescription: "foo",
		oOption: &Option{
			aliases:     []string{"foo"},
			Description: "foo",
		},
	}.assertOption())
	t.Run("multiple aliases", newOptionTester{
		iAliases:     []string{"foo", "bar"},
		iDescription: "foo",
		oOption: &Option{
			aliases:     []string{"foo", "bar"},
			Description: "foo",
		},
	}.assertOption())
	t.Run("empty description string", newOptionTester{
		iAliases: []string{"foo"},
		oOption:  &Option{aliases: []string{"foo"}},
	}.assertOption())
	t.Run("nil aliases", newOptionTester{
		iDescription: "foo",
		oErr:         errors.New("usage: option must have at least one alias"),
	}.assertNoAliasErr())
	t.Run("no aliases", newOptionTester{
		iAliases:     make([]string, 0),
		iDescription: "foo",
		oErr:         errors.New("usage: option must have at least one alias"),
	}.assertNoAliasErr())
	t.Run("single empty alias string", newOptionTester{
		iAliases:     []string{""},
		iDescription: "foo",
		oErr:         errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringErr())
	t.Run("multiple empty alias strings", newOptionTester{
		iAliases:     []string{"foo", "", "bar", ""},
		iDescription: "foo",
		oErr:         errors.New("usage: alias string must not be empty"),
	}.assertEmptyAliasStringErr())
}
