package usage

import "testing"

func assertAliasSlice(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d aliases returned but wanted %d", len(got), len(want))
	}
	for i, gotAlias := range got {
		if gotAlias != want[i] {
			t.Errorf("alias is %q but should be %q", gotAlias, want[i])
		}
	}
}

func assertArgSlice(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d args returned but wanted %d", len(got), len(want))
	}
	for i, gotArg := range got {
		if gotArg != want[i] {
			t.Errorf("arg is %q but should be %q", gotArg, want[i])
		}
	}
}

func assertOptionStruct(t *testing.T, got, want *option) {
	assertAliasSlice(t, got.aliases, want.aliases)
	if got.Description != want.Description {
		t.Errorf("description is %q but should be %q", got.Description, want.Description)
	}
	assertArgSlice(t, got.args, want.args)
}

type optionArgsTester struct {
	oArgs []string
}

func (tester optionArgsTester) assertOptionArgs() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := option{args: tester.oArgs}
		got := sampleOption.Args()
		assertArgSlice(t, got, tester.oArgs)
	}
}

type optionAddArgTester struct {
	iArg string
	oErr error
}

func (tester optionAddArgTester) assertOptionArgs() func(*testing.T) {
	return func(t *testing.T) {
		iterations := 3
		tempArgs := make([]string, 0, iterations)
		sampleOption := option{args: make([]string, 0)}
		for i := 1; i <= iterations; i++ {
			if gotErr := sampleOption.AddArg(tester.iArg); gotErr != nil {
				t.Errorf("got %q error but should be nil", gotErr)
			}
			tempArgs = append(tempArgs, tester.iArg)
		}
		assertArgSlice(t, sampleOption.args, tempArgs)
	}
}

func (tester optionAddArgTester) assertErrEmptyArgString() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := option{args: make([]string, 0)}
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
	oOption      *option
	oErr         error
}

func (tester newOptionTester) assertOption() func(*testing.T) {
	return func(t *testing.T) {
		got, gotErr := NewOption(tester.iAliases, tester.iDescription)
		if gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		assertAliasSlice(t, got.aliases, tester.oOption.aliases)
		if got.Description != tester.oOption.Description {
			t.Errorf("description is %q but should be %q", got.Description, tester.oOption.Description)
		}
		if got.args == nil || len(got.args) != 0 {
			t.Error("args not initialized to an empty slice")
		}
	}
}

func (tester newOptionTester) assertErrNoAliasProvided() func(*testing.T) {
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

func (tester newOptionTester) assertErrEmptyAliasString() func(*testing.T) {
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

type optionSetAliasesTester struct {
	iAliases []string
	oErr     error
}

func (tester optionSetAliasesTester) assertOptionAliases() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := option{aliases: []string{"foo"}}
		if gotErr := sampleOption.SetAliases(tester.iAliases); gotErr != nil {
			t.Errorf("got %q error but should be nil", gotErr)
		}
		assertAliasSlice(t, sampleOption.aliases, tester.iAliases)
	}
}

func (tester optionSetAliasesTester) assertErrNoAliasProvided() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := option{aliases: []string{"foo"}}
		got := sampleOption.SetAliases(tester.iAliases)
		if got == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		assertError(t, got, tester.oErr)
	}
}

func (tester optionSetAliasesTester) assertErrEmptyAliasString() func(*testing.T) {
	return func(t *testing.T) {
		sampleOption := option{aliases: []string{"foo"}}
		got := sampleOption.SetAliases(tester.iAliases)
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
		oErr: emptyArgStringErr(),
	}.assertErrEmptyArgString())
}

func TestNewOption(t *testing.T) {
	t.Run("baseline", newOptionTester{
		iAliases:     []string{"foo"},
		iDescription: "foo",
		oOption: &option{
			aliases:     []string{"foo"},
			Description: "foo",
		},
	}.assertOption())
	t.Run("multiple aliases", newOptionTester{
		iAliases:     []string{"foo", "bar"},
		iDescription: "foo",
		oOption: &option{
			aliases:     []string{"foo", "bar"},
			Description: "foo",
		},
	}.assertOption())
	t.Run("empty description string", newOptionTester{
		iAliases: []string{"foo"},
		oOption:  &option{aliases: []string{"foo"}},
	}.assertOption())
	t.Run("nil aliases", newOptionTester{
		iDescription: "foo",
		oErr:         noAliasProvidedErr(),
	}.assertErrNoAliasProvided())
	t.Run("no aliases", newOptionTester{
		iAliases:     make([]string, 0),
		iDescription: "foo",
		oErr:         noAliasProvidedErr(),
	}.assertErrNoAliasProvided())
	t.Run("single empty alias string", newOptionTester{
		iAliases:     []string{""},
		iDescription: "foo",
		oErr:         emptyAliasStringErr(),
	}.assertErrEmptyAliasString())
	t.Run("multiple empty alias strings", newOptionTester{
		iAliases:     []string{"foo", "", "bar", ""},
		iDescription: "foo",
		oErr:         emptyAliasStringErr(),
	}.assertErrEmptyAliasString())
}

func TestOptionSetAliases(t *testing.T) {
	t.Run("baseline", optionSetAliasesTester{
		iAliases: []string{"foo"},
	}.assertOptionAliases())
	t.Run("multiple aliases", optionSetAliasesTester{
		iAliases: []string{"foo", "bar"},
	}.assertOptionAliases())
	t.Run("nil aliases", optionSetAliasesTester{
		oErr: noAliasProvidedErr(),
	}.assertErrNoAliasProvided())
	t.Run("no aliases", optionSetAliasesTester{
		iAliases: make([]string, 0),
		oErr:     noAliasProvidedErr(),
	}.assertErrNoAliasProvided())
	t.Run("single empty alias string", optionSetAliasesTester{
		iAliases: []string{""},
		oErr:     emptyAliasStringErr(),
	}.assertErrEmptyAliasString())
	t.Run("multiple empty alias strings", optionSetAliasesTester{
		iAliases: []string{"foo", "", "bar", ""},
		oErr:     emptyAliasStringErr(),
	}.assertErrEmptyAliasString())
}
