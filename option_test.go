package usage

import (
	"errors"
	"testing"
)

func TestDefaultOptionAddArg(t *testing.T) {
	makeError := func(s string) error {
		return errors.New(s)
	}

	makeDefaultOption := func(aliases []string, description string) *defaultOption {
		return &defaultOption{
			aliases:     aliases,
			description: description,
			args:        make([]string, 0),
		}
	}

	assertOptionArgs := func(iArg string) func(*testing.T) {
		return func(t *testing.T) {
			option := makeDefaultOption([]string{"foo", "bar"}, "foo")
			option.AddArg(iArg)
			if len(option.args) != 1 {
				t.Fatalf("%d args returned but wanted 1", len(option.args))
			}
			if option.args[0] != iArg {
				t.Errorf("arg is %q but should be %q", option.args[0], iArg)
			}
		}
	}

	assertRepeatedOptionArgs := func(iArg string) func(*testing.T) {
		return func(t *testing.T) {
			option := makeDefaultOption([]string{"foo", "bar"}, "foo")
			option.AddArg(iArg)
			option.AddArg(iArg)
			option.AddArg(iArg)
			if len(option.args) != 1 {
				t.Fatalf("%d args returned but wanted 1", len(option.args))
			}
			if option.args[0] != iArg {
				t.Errorf("arg is %q but should be %q", option.args[0], iArg)
			}
		}
	}

	assertEmptyArgStringError := func(iArg string, want error) func(*testing.T) {
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

	t.Run("baseline", assertOptionArgs("foo"))
	t.Run("repeated arg strings", assertRepeatedOptionArgs("foo"))
	t.Run("empty arg string", assertEmptyArgStringError("", makeError("usage: arg string must not be empty")))
}

func TestNewOption(t *testing.T) {
	makeError := func(s string) error {
		return errors.New(s)
	}

	makeDefaultOption := func(aliases []string, description string) *defaultOption {
		return &defaultOption{
			aliases:     aliases,
			description: description,
			args:        make([]string, 0),
		}
	}

	assertOption := func(iAliases []string, iDescription string, want *defaultOption) func(*testing.T) {
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

	assertEmptyOptionAliasStringError := func(iAliases []string, iDescription string, want error) func(*testing.T) {
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

	assertNoOptionAliasProvidedError := func(iAliases []string, iDescription string, want error) func(*testing.T) {
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

	t.Run("baseline", assertOption(
		[]string{"foo", "bar"},
		"foo",
		makeDefaultOption([]string{"foo", "bar"}, "foo"),
	))

	t.Run("single alias", assertOption(
		[]string{"foo"},
		"foo",
		makeDefaultOption([]string{"foo"}, "foo"),
	))

	t.Run("single repeated alias", assertOption(
		[]string{"foo", "foo"},
		"foo",
		makeDefaultOption([]string{"foo"}, "foo"),
	))

	t.Run("multiple repeated aliases", assertOption(
		[]string{"foo", "bar", "foo", "bar"},
		"foo",
		makeDefaultOption([]string{"foo", "bar"}, "foo"),
	))

	t.Run("empty description string", assertOption(
		[]string{"foo", "bar"},
		"",
		makeDefaultOption([]string{"foo", "bar"}, ""),
	))

	t.Run("nil aliases", assertNoOptionAliasProvidedError(
		nil,
		"foo",
		makeError("usage: option must have at least one alias"),
	))

	t.Run("no aliases", assertNoOptionAliasProvidedError(
		make([]string, 0),
		"foo",
		makeError("usage: option must have at least one alias"),
	))

	t.Run("single empty alias string", assertEmptyOptionAliasStringError(
		[]string{""},
		"foo",
		makeError("usage: alias string must not be empty"),
	))

	t.Run("multiple empty alias strings", assertEmptyOptionAliasStringError(
		[]string{"foo", "", "bar", ""},
		"foo",
		makeError("usage: alias string must not be empty"),
	))
}
