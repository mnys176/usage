package usage

import (
	"errors"
	"testing"
)

func makeError(s string) error {
	return errors.New(s)
}

func makeOption(aliases []string, description string) Option {
	return &defaultOption{
		aliases:     aliases,
		description: description,
		args:        make([]string, 0),
	}
}
func makeEntry(name string, description string) Entry {
	return &defaultEntry{
		name:        name,
		description: description,
		options:     make([]Option, 0),
		args:        make([]string, 0),
	}
}

func assertOptionArgs_OptionAddArg(iArg string) func(*testing.T) {
	return func(t *testing.T) {
		option := makeOption([]string{"foo", "bar"}, "foo")
		if err := option.AddArg(iArg); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		optionArgs := option.Args()
		if len(optionArgs) != 1 {
			t.Fatalf("%d args returned but wanted 1", len(optionArgs))
		}
		if optionArgs[0] != iArg {
			t.Errorf("arg is %q but should be %q", optionArgs[0], iArg)
		}
	}
}

func assertRepeatedOptionArgs_OptionAddArg(iArg string) func(*testing.T) {
	return func(t *testing.T) {
		option := makeOption([]string{"foo", "bar"}, "foo")
		err1 := option.AddArg(iArg)
		err2 := option.AddArg(iArg)
		err3 := option.AddArg(iArg)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		optionArgs := option.Args()
		if len(optionArgs) != 3 {
			t.Fatalf("%d args returned but wanted 3", len(optionArgs))
		}
		if optionArgs[0] != iArg {
			t.Errorf("arg is %q but should be %q", optionArgs[0], iArg)
		}
	}
}

func assertEmptyArgStringError_OptionAddArg(iArg string, want error) func(*testing.T) {
	return func(t *testing.T) {
		option := makeOption([]string{"foo", "bar"}, "foo")
		oErr := option.AddArg(iArg)
		if oErr == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertOption_NewOption(iAliases []string, iDescription string, want Option) func(*testing.T) {
	return func(t *testing.T) {
		oOption, _ := NewOption(iAliases, iDescription)
		oOptionAliases := oOption.Aliases()
		wantAliases := want.Aliases()
		if len(oOptionAliases) != len(wantAliases) {
			t.Fatalf("%d aliases returned but wanted %d", len(oOptionAliases), len(wantAliases))
		}
		for i := range oOptionAliases {
			if oOptionAliases[i] != wantAliases[i] {
				t.Errorf("alias is %q but should be %q", oOptionAliases[i], wantAliases[i])
			}
		}
		oOptionDescription := oOption.Description()
		wantDescription := want.Description()
		if oOptionDescription != wantDescription {
			t.Errorf("description is %q but should be %q", oOptionDescription, wantDescription)
		}
		oOptionArgs := oOption.Args()
		wantArgs := want.Args()
		if len(oOptionArgs) != len(wantArgs) {
			t.Fatalf("%d args returned but wanted %d", len(oOptionArgs), len(wantArgs))
		}
		for i := range oOptionArgs {
			if oOptionArgs[i] != wantArgs[i] {
				t.Errorf("arg is %q but should be %q", oOptionArgs[i], wantArgs[i])
			}
		}
	}
}

func assertEmptyOptionAliasStringError_NewOption(iAliases []string, iDescription string, want error) func(*testing.T) {
	return func(t *testing.T) {
		_, oErr := NewOption(iAliases, iDescription)
		if oErr == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertNoOptionAliasProvidedError_NewOption(iAliases []string, iDescription string, want error) func(*testing.T) {
	return func(t *testing.T) {
		_, oErr := NewOption(iAliases, iDescription)
		if oErr == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertEntryArgs_EntryAddArg(iArg string) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		if err := entry.AddArg(iArg); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		entryArgs := entry.Args()
		if len(entryArgs) != 1 {
			t.Fatalf("%d args returned but wanted 1", len(entryArgs))
		}
		if entryArgs[0] != iArg {
			t.Errorf("arg is %q but should be %q", entryArgs[0], iArg)
		}
	}
}

func assertRepeatedEntryArgs_EntryAddArg(iArg string) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err1 := entry.AddArg(iArg)
		err2 := entry.AddArg(iArg)
		err3 := entry.AddArg(iArg)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		entryArgs := entry.Args()
		if len(entryArgs) != 3 {
			t.Fatalf("%d args returned but wanted 3", len(entryArgs))
		}
		if entryArgs[0] != iArg {
			t.Errorf("arg is %q but should be %q", entryArgs[0], iArg)
		}
	}
}

func assertEmptyArgStringError_EntryAddArg(iArg string, want error) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		oErr := entry.AddArg(iArg)
		if oErr == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertEntryOptions_EntryAddOption(iOption Option) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		if err := entry.AddOption(iOption); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		entryOptions := entry.Options()
		if len(entryOptions) != 1 {
			t.Fatalf("%d options returned but wanted 1", len(entryOptions))
		}
		entryOption0Aliases := entryOptions[0].Aliases()
		iOptionAliases := iOption.Aliases()
		if len(entryOption0Aliases) != len(iOptionAliases) {
			t.Fatalf(
				"%d option aliases returned but wanted %d",
				len(entryOption0Aliases),
				len(iOption.Aliases()),
			)
		}
		for i, alias := range entryOption0Aliases {
			if alias != iOptionAliases[i] {
				t.Errorf("option alias is %q but should be %q", alias, iOptionAliases[i])
			}
		}
	}
}

func assertRepeatedEntryOptions_EntryAddOption(iOption Option) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		err1 := entry.AddOption(iOption)
		err2 := entry.AddOption(iOption)
		err3 := entry.AddOption(iOption)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		entryOptions := entry.Options()
		if len(entryOptions) != 3 {
			t.Fatalf("%d options returned but wanted 3", len(entryOptions))
		}
		entryOption0Aliases := entryOptions[0].Aliases()
		iOptionAliases := iOption.Aliases()
		if len(entryOptions[0].Aliases()) != len(iOptionAliases) {
			t.Fatalf(
				"%d option aliases returned but wanted %d",
				len(entryOption0Aliases),
				len(iOption.Aliases()),
			)
		}
		for i, alias := range entryOption0Aliases {
			if alias != iOptionAliases[i] {
				t.Errorf("option alias is %q but should be %q", alias, iOptionAliases[i])
			}
		}
	}
}

func assertEmptyOptionAliasStringError_EntryAddOption(iOption Option, want error) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		oErr := entry.AddOption(iOption)
		if oErr == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertNoOptionAliasProvidedError_EntryAddOption(iOption Option, want error) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeEntry("foo", "foo")
		oErr := entry.AddOption(iOption)
		if oErr == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertEntry_NewEntry(iName, iDescription string, want Entry) func(*testing.T) {
	return func(t *testing.T) {
		oEntry, _ := NewEntry(iName, iDescription)
		oEntryName := oEntry.Name()
		wantName := want.Name()
		if oEntryName != wantName {
			t.Errorf("name is %q but should be %q", oEntryName, wantName)
		}
		oEntryDescription := oEntry.Description()
		wantDescription := want.Description()
		if oEntryDescription != wantDescription {
			t.Errorf("description is %q but should be %q", oEntryDescription, wantDescription)
		}
		oEntryArgs := oEntry.Args()
		if oEntryArgs == nil || len(oEntryArgs) != 0 {
			t.Error("args not initialized to an empty slice")
		}
		oEntryOptions := oEntry.Options()
		if oEntryOptions == nil || len(oEntryOptions) != 0 {
			t.Error("options not initialized to an empty slice")
		}
	}
}

func assertEmptyNameStringError_NewEntry(iName, iDescription string, want error) func(*testing.T) {
	return func(t *testing.T) {
		_, oErr := NewEntry(iName, iDescription)
		if oErr == nil {
			t.Fatal("no error returned with an empty name string")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}
