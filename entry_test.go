package usage

import (
	"errors"
	"testing"
)

func TestDefaultEntryAddArg(t *testing.T) {
	makeError := errors.New

	makeDefaultEntry := func(name string, description string) *defaultEntry {
		return &defaultEntry{
			name:        name,
			description: description,
			options:     make([]Option, 0),
			args:        make([]string, 0),
		}
	}

	assertEntryArgs := func(iArg string) func(*testing.T) {
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

	assertRepeatedEntryArgs := func(iArg string) func(*testing.T) {
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

	assertEmptyArgStringError := func(iArg string, want error) func(*testing.T) {
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

	t.Run("baseline", assertEntryArgs("foo"))
	t.Run("repeated arg strings", assertRepeatedEntryArgs("foo"))
	t.Run("empty arg string", assertEmptyArgStringError("", makeError("usage: arg string must not be empty")))
}

func TestDefaultEntryAddOption(t *testing.T) {
	makeError := errors.New

	makeDefaultOption := func(aliases []string, description string) *defaultOption {
		return &defaultOption{
			aliases:     aliases,
			description: description,
			args:        make([]string, 0),
		}
	}

	makeDefaultEntry := func(name string, description string) *defaultEntry {
		return &defaultEntry{
			name:        name,
			description: description,
			options:     make([]Option, 0),
			args:        make([]string, 0),
		}
	}

	assertEntryOptions := func(iOption Option) func(*testing.T) {
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

	assertRepeatedEntryOptions := func(iOption Option) func(*testing.T) {
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

	assertEmptyOptionAliasStringError := func(iOption Option, want error) func(*testing.T) {
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

	assertNoOptionAliasProvidedError := func(iOption Option, want error) func(*testing.T) {
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

	t.Run("baseline", assertEntryOptions(makeDefaultOption([]string{"foo", "bar"}, "foo")))
	t.Run("repeated options", assertRepeatedEntryOptions(makeDefaultOption([]string{"foo", "bar"}, "foo")))
	t.Run("nil option aliases", assertNoOptionAliasProvidedError(
		makeDefaultOption(nil, "foo"),
		makeError("usage: option must have at least one alias"),
	))
	t.Run("no option aliases", assertNoOptionAliasProvidedError(
		makeDefaultOption(make([]string, 0), "foo"),
		makeError("usage: option must have at least one alias"),
	))
	t.Run("single empty option alias string", assertEmptyOptionAliasStringError(
		makeDefaultOption([]string{""}, "foo"),
		makeError("usage: alias string must not be empty"),
	))
	t.Run("multiple empty option alias strings", assertEmptyOptionAliasStringError(
		makeDefaultOption([]string{"foo", "", "bar", ""}, "foo"),
		makeError("usage: alias string must not be empty"),
	))
}

func TestNewEntry(t *testing.T) {
	makeError := errors.New

	makeDefaultEntry := func(name string, description string) *defaultEntry {
		return &defaultEntry{
			name:        name,
			description: description,
			options:     make([]Option, 0),
			args:        make([]string, 0),
		}
	}

	assertEntry := func(iName, iDescription string, want *defaultEntry) func(*testing.T) {
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

	assertEmptyNameStringError := func(iName, iDescription string, want error) func(*testing.T) {
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

	t.Run("baseline", assertEntry(
		"foo",
		"foo",
		makeDefaultEntry("foo", "foo"),
	))
	t.Run("empty description string", assertEntry(
		"foo",
		"",
		makeDefaultEntry("foo", ""),
	))
	t.Run("empty name string", assertEmptyNameStringError(
		"",
		"foo",
		makeError("usage: name string must not be empty"),
	))
}
