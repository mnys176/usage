package usage

import (
	"testing"
)

func TestDefaultOptionAddArg(t *testing.T) {
	t.Run("baseline", assertDefaultOptionAddArgOptionArgs("foo"))
	t.Run("repeated arg strings", assertDefaultOptionAddArgRepeatedOptionArgs("foo"))
	t.Run("empty arg string", assertDefaultOptionAddArgEmptyArgStringError("", makeError("usage: arg string must not be empty")))
}

func TestNewOption(t *testing.T) {
	t.Run("baseline", assertNewOptionOption(
		[]string{"foo", "bar"},
		"foo",
		makeDefaultOption([]string{"foo", "bar"}, "foo"),
	))
	t.Run("single alias", assertNewOptionOption(
		[]string{"foo"},
		"foo",
		makeDefaultOption([]string{"foo"}, "foo"),
	))
	t.Run("single repeated alias", assertNewOptionOption(
		[]string{"foo", "foo"},
		"foo",
		makeDefaultOption([]string{"foo", "foo"}, "foo"),
	))
	t.Run("multiple repeated aliases", assertNewOptionOption(
		[]string{"foo", "bar", "foo", "bar"},
		"foo",
		makeDefaultOption([]string{"foo", "bar", "foo", "bar"}, "foo"),
	))
	t.Run("empty description string", assertNewOptionOption(
		[]string{"foo", "bar"},
		"",
		makeDefaultOption([]string{"foo", "bar"}, ""),
	))
	t.Run("nil aliases", assertNewOptionNoOptionAliasProvidedError(
		nil,
		"foo",
		makeError("usage: option must have at least one alias"),
	))
	t.Run("no aliases", assertNewOptionNoOptionAliasProvidedError(
		make([]string, 0),
		"foo",
		makeError("usage: option must have at least one alias"),
	))
	t.Run("single empty alias string", assertNewOptionEmptyOptionAliasStringError(
		[]string{""},
		"foo",
		makeError("usage: alias string must not be empty"),
	))
	t.Run("multiple empty alias strings", assertNewOptionEmptyOptionAliasStringError(
		[]string{"foo", "", "bar", ""},
		"foo",
		makeError("usage: alias string must not be empty"),
	))
}
