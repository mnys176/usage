package usage

import (
	"errors"
	"testing"
)

func makeError(s string) error {
	return errors.New(s)
}

func makeDefaultOption(aliases []string, description string) *defaultOption {
	return &defaultOption{
		aliases:     aliases,
		description: description,
		args:        make([]string, 0),
	}
}
func makeDefaultEntry(name string, description string) *defaultEntry {
	return &defaultEntry{
		name:        name,
		description: description,
		options:     make([]Option, 0),
		args:        make([]string, 0),
	}
}

func assertDefaultOptionAddArgOptionArgs(iArg string) func(*testing.T) {
	return func(t *testing.T) {
		option := makeDefaultOption([]string{"foo", "bar"}, "foo")
		if err := option.AddArg(iArg); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		if len(option.args) != 1 {
			t.Fatalf("%d args returned but wanted 1", len(option.args))
		}
		if option.args[0] != iArg {
			t.Errorf("arg is %q but should be %q", option.args[0], iArg)
		}
	}
}

func assertDefaultOptionAddArgRepeatedOptionArgs(iArg string) func(*testing.T) {
	return func(t *testing.T) {
		option := makeDefaultOption([]string{"foo", "bar"}, "foo")
		err1 := option.AddArg(iArg)
		err2 := option.AddArg(iArg)
		err3 := option.AddArg(iArg)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		if len(option.args) != 3 {
			t.Fatalf("%d args returned but wanted 3", len(option.args))
		}
		if option.args[0] != iArg {
			t.Errorf("arg is %q but should be %q", option.args[0], iArg)
		}
	}
}

func assertDefaultOptionAddArgEmptyArgStringError(iArg string, want error) func(*testing.T) {
	return func(t *testing.T) {
		option := makeDefaultOption([]string{"foo", "bar"}, "foo")
		oErr := option.AddArg(iArg)
		if oErr == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertNewOptionOption(iAliases []string, iDescription string, want *defaultOption) func(*testing.T) {
	return func(t *testing.T) {
		oOption, _ := NewOption(iAliases, iDescription)
		aliases := oOption.Aliases()
		if len(aliases) != len(want.aliases) {
			t.Fatalf("%d aliases returned but wanted %d", len(aliases), len(want.aliases))
		}
		for i := range aliases {
			if aliases[i] != want.aliases[i] {
				t.Errorf("alias is %q but should be %q", aliases[i], want.aliases[i])
			}
		}
		description := oOption.Description()
		if description != want.description {
			t.Errorf("description is %q but should be %q", description, want.description)
		}
		args := oOption.Args()
		if len(args) != len(want.args) {
			t.Fatalf("%d args returned but wanted %d", len(args), len(want.args))
		}
		for i := range args {
			if args[i] != want.args[i] {
				t.Errorf("arg is %q but should be %q", args[i], want.args[i])
			}
		}
	}
}

func assertNewOptionEmptyOptionAliasStringError(iAliases []string, iDescription string, want error) func(*testing.T) {
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

func assertNewOptionNoOptionAliasProvidedError(iAliases []string, iDescription string, want error) func(*testing.T) {
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

func assertDefaultEntryAddArgEntryArgs(iArg string) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeDefaultEntry("foo", "foo")
		if err := entry.AddArg(iArg); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		if len(entry.args) != 1 {
			t.Fatalf("%d args returned but wanted 1", len(entry.args))
		}
		if entry.args[0] != iArg {
			t.Errorf("arg is %q but should be %q", entry.args[0], iArg)
		}
	}
}

func assertDefaultEntryAddArgRepeatedEntryArgs(iArg string) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeDefaultEntry("foo", "foo")
		err1 := entry.AddArg(iArg)
		err2 := entry.AddArg(iArg)
		err3 := entry.AddArg(iArg)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		if len(entry.args) != 3 {
			t.Fatalf("%d args returned but wanted 3", len(entry.args))
		}
		if entry.args[0] != iArg {
			t.Errorf("arg is %q but should be %q", entry.args[0], iArg)
		}
	}
}

func assertDefaultEntryAddArgEmptyArgStringError(iArg string, want error) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeDefaultEntry("foo", "foo")
		oErr := entry.AddArg(iArg)
		if oErr == nil {
			t.Fatal("no error returned with an empty arg string")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertDefaultEntryAddOptionEntryOptions(iOption Option) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeDefaultEntry("foo", "foo")
		if err := entry.AddOption(iOption); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		if len(entry.options) != 1 {
			t.Fatalf("%d options returned but wanted 1", len(entry.options))
		}
		if len(entry.options[0].Aliases()) != len(iOption.Aliases()) {
			t.Fatalf(
				"%d option aliases returned but wanted %d",
				len(entry.options[0].Aliases()),
				len(iOption.Aliases()),
			)
		}
		iOptionAliases := iOption.Aliases()
		for i, alias := range entry.options[0].Aliases() {
			if alias != iOptionAliases[i] {
				t.Errorf("option alias is %q but should be %q", alias, iOptionAliases[i])
			}
		}
	}
}

func assertDefaultEntryAddOptionRepeatedEntryOptions(iOption Option) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeDefaultEntry("foo", "foo")
		err1 := entry.AddOption(iOption)
		err2 := entry.AddOption(iOption)
		err3 := entry.AddOption(iOption)
		if err := errors.Join(err1, err2, err3); err != nil {
			t.Errorf("got %q error but should be nil", err)
		}
		if len(entry.options) != 3 {
			t.Fatalf("%d options returned but wanted 3", len(entry.options))
		}
		if len(entry.options[0].Aliases()) != len(iOption.Aliases()) {
			t.Fatalf(
				"%d option aliases returned but wanted %d",
				len(entry.options[0].Aliases()),
				len(iOption.Aliases()),
			)
		}
		iOptionAliases := iOption.Aliases()
		for i, alias := range entry.options[0].Aliases() {
			if alias != iOptionAliases[i] {
				t.Errorf("option alias is %q but should be %q", alias, iOptionAliases[i])
			}
		}
	}
}

func assertDefaultEntryAddOptionEmptyOptionAliasStringError(iOption Option, want error) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeDefaultEntry("foo", "foo")
		oErr := entry.AddOption(iOption)
		if oErr == nil {
			t.Fatal("no error returned with an empty alias string")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertDefaultEntryAddOptionNoOptionAliasProvidedError(iOption Option, want error) func(*testing.T) {
	return func(t *testing.T) {
		entry := makeDefaultEntry("foo", "foo")
		oErr := entry.AddOption(iOption)
		if oErr == nil {
			t.Fatal("no error returned with no provided aliases")
		}
		if !errors.Is(oErr, want) {
			t.Errorf("got %q error but wanted %q", oErr, want)
		}
	}
}

func assertNewEntryEntry(iName, iDescription string, want *defaultEntry) func(*testing.T) {
	return func(t *testing.T) {
		oEntry, _ := NewEntry(iName, iDescription)
		if name := oEntry.Name(); name != want.name {
			t.Errorf("name is %q but should be %q", name, want.name)
		}
		if description := oEntry.Description(); description != want.description {
			t.Errorf("description is %q but should be %q", description, want.description)
		}
		if args := oEntry.Args(); args == nil || len(args) != 0 {
			t.Error("args not initialized to an empty slice")
		}
		if options := oEntry.Options(); options == nil || len(options) != 0 {
			t.Error("options not initialized to an empty slice")
		}
	}
}

func assertNewEntryEmptyNameStringError(iName, iDescription string, want error) func(*testing.T) {
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
