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
		entryArgs := entry.Args()
		if len(entryArgs) != 1 {
			t.Fatalf("%d args returned but wanted 1", len(entryArgs))
		}
		if entryArgs[0] != eaat.iArg {
			t.Errorf("arg is %q but should be %q", entryArgs[0], eaat.iArg)
		}
	}
}

func (eaat entryAddArgTester) assertRepeatedEntryArgs() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err1 := entry.AddArg(eaat.iArg)
		err2 := entry.AddArg(eaat.iArg)
		err3 := entry.AddArg(eaat.iArg)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		entryArgs := entry.Args()
		if len(entryArgs) != 3 {
			t.Fatalf("%d args returned but wanted 3", len(entryArgs))
		}
		if entryArgs[0] != eaat.iArg {
			t.Errorf("arg is %q but should be %q", entryArgs[0], eaat.iArg)
		}
	}
}

func (eaat entryAddArgTester) assertEmptyArgStringError() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err := entry.AddArg(eaat.iArg)
		if err == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		if !errors.Is(err, eaat.oErr) {
			t.Errorf("got %q error but wanted %q", err, eaat.oErr)
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
		entryOptions := entry.Options()
		if len(entryOptions) != 1 {
			t.Fatalf("%d options returned but wanted 1", len(entryOptions))
		}
		errorOption0Aliases := entryOptions[0].Aliases()
		iOptionAliases := eaot.iOption.Aliases()
		if len(errorOption0Aliases) != len(iOptionAliases) {
			t.Fatalf(
				"%d option aliases returned but wanted %d",
				len(errorOption0Aliases),
				len(iOptionAliases),
			)
		}
		for i, alias := range errorOption0Aliases {
			if alias != iOptionAliases[i] {
				t.Errorf("option alias is %q but should be %q", alias, iOptionAliases[i])
			}
		}
	}
}

func (eaot entryAddOptionTester) assertRepeatedEntryOptions() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err1 := entry.AddOption(eaot.iOption)
		err2 := entry.AddOption(eaot.iOption)
		err3 := entry.AddOption(eaot.iOption)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		entryOptions := entry.Options()
		if len(entryOptions) != 3 {
			t.Fatalf("%d options returned but wanted 3", len(entryOptions))
		}
		entryOptionAliases := entryOptions[0].Aliases()
		iOptionAliases := eaot.iOption.Aliases()
		if len(entryOptions[0].Aliases()) != len(iOptionAliases) {
			t.Fatalf(
				"%d option aliases returned but wanted %d",
				len(entryOptionAliases),
				len(iOptionAliases),
			)
		}
		for i, alias := range entryOptionAliases {
			if alias != iOptionAliases[i] {
				t.Errorf("option alias is %q but should be %q", alias, iOptionAliases[i])
			}
		}
	}
}

func (eaot entryAddOptionTester) assertEmptyOptionAliasStringError() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err := entry.AddOption(eaot.iOption)
		if err == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		if !errors.Is(err, eaot.oErr) {
			t.Errorf("got %q error but wanted %q", err, eaot.oErr)
		}
	}
}

func (eaot entryAddOptionTester) assertNoOptionAliasProvidedError() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err := entry.AddOption(eaot.iOption)
		if err == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		if !errors.Is(err, eaot.oErr) {
			t.Errorf("got %q error but wanted %q", err, eaot.oErr)
		}
	}
}

func (eaot entryAddOptionTester) assertNoOptionProvidedError() func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err := entry.AddOption(eaot.iOption)
		if err == nil {
			t.Fatal("no error returned with nil option")
		}
		if !errors.Is(err, eaot.oErr) {
			t.Errorf("got %q error but wanted %q", err, eaot.oErr)
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
		entry, _ := NewEntry(net.iName, net.iDescription)
		entryName := entry.Name()
		oEntryName := net.oEntry.Name()
		if entryName != oEntryName {
			t.Errorf("name is %q but should be %q", entryName, oEntryName)
		}
		entryDescription := entry.Description()
		oEntryDescription := net.oEntry.Description()
		if entryDescription != oEntryDescription {
			t.Errorf("description is %q but should be %q", entryDescription, oEntryDescription)
		}
		entryArgs := entry.Args()
		if entryArgs == nil || len(entryArgs) != 0 {
			t.Error("args not initialized to an empty slice")
		}
		entryOptions := entry.Options()
		if entryOptions == nil || len(entryOptions) != 0 {
			t.Error("options not initialized to an empty slice")
		}
	}
}

func (net newEntryTester) assertEmptyNameStringError() func(*testing.T) {
	return func(t *testing.T) {
		_, err := NewEntry(net.iName, net.iDescription)
		if err == nil {
			t.Fatal("no error returned with an empty name string")
		}
		if !errors.Is(err, net.oErr) {
			t.Errorf("got %q error but wanted %q", err, net.oErr)
		}
	}
}

func TestEntryAddArg(t *testing.T) {
	t.Run("baseline", entryAddArgTester{
		iArg: "foo",
	}.assertEntryArgs())
	t.Run("repeated arg strings", entryAddArgTester{
		iArg: "foo",
	}.assertRepeatedEntryArgs())
	t.Run("empty arg string", entryAddArgTester{
		oErr: makeError("usage: arg string must not be empty"),
	}.assertEmptyArgStringError())
}

func TestEntryAddOption(t *testing.T) {
	t.Run("baseline", entryAddOptionTester{
		iOption: makeOption([]string{"foo", "bar"}, "foo"),
	}.assertEntryOptions())
	t.Run("repeated options", entryAddOptionTester{
		iOption: makeOption([]string{"foo", "bar"}, "foo"),
	}.assertRepeatedEntryOptions())
	t.Run("nil option", entryAddOptionTester{
		oErr: makeError("usage: no option provided"),
	}.assertNoOptionProvidedError())
	t.Run("nil option aliases", entryAddOptionTester{
		iOption: makeOption(nil, "foo"),
		oErr:    makeError("usage: option must have at least one alias"),
	}.assertNoOptionAliasProvidedError())
	t.Run("no option aliases", entryAddOptionTester{
		iOption: makeOption(make([]string, 0), "foo"),
		oErr:    makeError("usage: option must have at least one alias"),
	}.assertNoOptionAliasProvidedError())
	t.Run("single empty option alias string", entryAddOptionTester{
		iOption: makeOption([]string{""}, "foo"),
		oErr:    makeError("usage: alias string must not be empty"),
	}.assertEmptyOptionAliasStringError())
	t.Run("multiple empty option alias strings", entryAddOptionTester{
		iOption: makeOption([]string{"foo", "", "bar", ""}, "foo"),
		oErr:    makeError("usage: alias string must not be empty"),
	}.assertEmptyOptionAliasStringError())
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
	}.assertEmptyNameStringError())
}
