package usage

import (
	"errors"
	"testing"
)

type entryAddArgTester struct {
	iArg string
	oErr error
}

func (eaat entryAddArgTester) assertEntryArgs() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		if err := entry.AddArg(eaat.iArg); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		got := entry.Args()
		if len(got) != 1 {
			t.Fatalf("%d args returned but wanted 1", len(got))
		}
		if got[0] != eaat.iArg {
			t.Errorf("arg is %q but should be %q", got[0], eaat.iArg)
		}
	}
}

func (eaat entryAddArgTester) assertEntryArgsRepeated() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err1 := entry.AddArg(eaat.iArg)
		err2 := entry.AddArg(eaat.iArg)
		err3 := entry.AddArg(eaat.iArg)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		got := entry.Args()
		if len(got) != 3 {
			t.Fatalf("%d args returned but wanted 3", len(got))
		}
		if got[0] != eaat.iArg {
			t.Errorf("arg is %q but should be %q", got[0], eaat.iArg)
		}
	}
}

func (eaat entryAddArgTester) assertErrEmptyArgString() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		got := entry.AddArg(eaat.iArg)
		if got == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		if !errors.Is(got, eaat.oErr) {
			t.Errorf("got %q error but wanted %q", got, eaat.oErr)
		}
	}
}

type entryAddOptionTester struct {
	iOption Option
	oErr    error
}

func (eaot entryAddOptionTester) assertEntryOptions() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		if err := entry.AddOption(eaot.iOption); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		got := entry.Options()
		if len(got) != 1 {
			t.Fatalf("%d options returned but wanted 1", len(got))
		}
		got0Aliases := got[0].Aliases()
		iOptionAliases := eaot.iOption.Aliases()
		if len(got0Aliases) != len(iOptionAliases) {
			t.Fatalf(
				"%d option aliases returned but wanted %d",
				len(got0Aliases),
				len(iOptionAliases),
			)
		}
		for i, alias := range got0Aliases {
			if alias != iOptionAliases[i] {
				t.Errorf("option alias is %q but should be %q", alias, iOptionAliases[i])
			}
		}
	}
}

func (eaot entryAddOptionTester) assertEntryOptionsRepeated() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err1 := entry.AddOption(eaot.iOption)
		err2 := entry.AddOption(eaot.iOption)
		err3 := entry.AddOption(eaot.iOption)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		got := entry.Options()
		if len(got) != 3 {
			t.Fatalf("%d options returned but wanted 3", len(got))
		}
		got0Aliases := got[0].Aliases()
		iOptionAliases := eaot.iOption.Aliases()
		if len(got0Aliases) != len(iOptionAliases) {
			t.Fatalf(
				"%d option aliases returned but wanted %d",
				len(got0Aliases),
				len(iOptionAliases),
			)
		}
		for i, alias := range got0Aliases {
			if alias != iOptionAliases[i] {
				t.Errorf("option alias is %q but should be %q", alias, iOptionAliases[i])
			}
		}
	}
}

func (eaot entryAddOptionTester) assertErrorEmptyOptionAliasString() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		got := entry.AddOption(eaot.iOption)
		if got == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		if !errors.Is(got, eaot.oErr) {
			t.Errorf("got %q error but wanted %q", got, eaot.oErr)
		}
	}
}

func (eaot entryAddOptionTester) assertErrorNoOptionAliasProvided() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		got := entry.AddOption(eaot.iOption)
		if got == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		if !errors.Is(got, eaot.oErr) {
			t.Errorf("got %q error but wanted %q", got, eaot.oErr)
		}
	}
}

func (eaot entryAddOptionTester) assertErrorNoOptionProvided() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		got := entry.AddOption(eaot.iOption)
		if got == nil {
			t.Fatal("no error returned with nil option")
		}
		if !errors.Is(got, eaot.oErr) {
			t.Errorf("got %q error but wanted %q", got, eaot.oErr)
		}
	}
}

type newEntryTester struct {
	iName        string
	iDescription string
	oEntry       Entry
	oErr         error
}

func (net newEntryTester) assertEntry() func(*testing.T) {
	return func(t *testing.T) {
		got, _ := NewEntry(net.iName, net.iDescription)
		gotName := got.Name()
		oEntryName := net.oEntry.Name()
		if gotName != oEntryName {
			t.Errorf("name is %q but should be %q", gotName, oEntryName)
		}
		gotDescription := got.Description()
		oEntryDescription := net.oEntry.Description()
		if gotDescription != oEntryDescription {
			t.Errorf("description is %q but should be %q", gotDescription, oEntryDescription)
		}
		gotArgs := got.Args()
		if gotArgs == nil || len(gotArgs) != 0 {
			t.Error("args not initialized to an empty slice")
		}
		gotOptions := got.Options()
		if gotOptions == nil || len(gotOptions) != 0 {
			t.Error("options not initialized to an empty slice")
		}
	}
}

func (net newEntryTester) assertErrorEmptyNameString() func(*testing.T) {
	return func(t *testing.T) {
		_, got := NewEntry(net.iName, net.iDescription)
		if got == nil {
			t.Fatal("no error returned with an empty name string")
		}
		if !errors.Is(got, net.oErr) {
			t.Errorf("got %q error but wanted %q", got, net.oErr)
		}
	}
}

func TestEntryAddArg(t *testing.T) {
	t.Run("baseline", entryAddArgTester{
		iArg: "foo",
	}.assertEntryArgs())
	t.Run("repeated arg strings", entryAddArgTester{
		iArg: "foo",
	}.assertEntryArgsRepeated())
	t.Run("empty arg string", entryAddArgTester{
		oErr: makeError("usage: arg string must not be empty"),
	}.assertErrEmptyArgString())
}

func TestEntryAddOption(t *testing.T) {
	t.Run("baseline", entryAddOptionTester{
		iOption: makeOption([]string{"foo", "bar"}, "foo"),
	}.assertEntryOptions())
	t.Run("repeated options", entryAddOptionTester{
		iOption: makeOption([]string{"foo", "bar"}, "foo"),
	}.assertEntryOptionsRepeated())
	t.Run("nil option", entryAddOptionTester{
		oErr: makeError("usage: no option provided"),
	}.assertErrorNoOptionProvided())
	t.Run("nil option aliases", entryAddOptionTester{
		iOption: makeOption(nil, "foo"),
		oErr:    makeError("usage: option must have at least one alias"),
	}.assertErrorNoOptionAliasProvided())
	t.Run("no option aliases", entryAddOptionTester{
		iOption: makeOption(make([]string, 0), "foo"),
		oErr:    makeError("usage: option must have at least one alias"),
	}.assertErrorNoOptionAliasProvided())
	t.Run("single empty option alias string", entryAddOptionTester{
		iOption: makeOption([]string{""}, "foo"),
		oErr:    makeError("usage: alias string must not be empty"),
	}.assertErrorEmptyOptionAliasString())
	t.Run("multiple empty option alias strings", entryAddOptionTester{
		iOption: makeOption([]string{"foo", "", "bar", ""}, "foo"),
		oErr:    makeError("usage: alias string must not be empty"),
	}.assertErrorEmptyOptionAliasString())
}

func TestNewEntry(t *testing.T) {
	t.Run("baseline", newEntryTester{
		iName:        "foo",
		iDescription: "foo",
		oEntry:       makeEntry("foo", "foo"),
	}.assertEntry())
	t.Run("empty description string", newEntryTester{
		iName:  "foo",
		oEntry: makeEntry("foo", ""),
	}.assertEntry())
	t.Run("empty name string", newEntryTester{
		iDescription: "foo",
		oErr:         makeError("usage: name string must not be empty"),
	}.assertErrorEmptyNameString())
}
