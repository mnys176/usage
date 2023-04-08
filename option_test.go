package usage

import (
	"testing"
)

func TestOptionAddArg(t *testing.T) {
	t.Run("baseline", assertOptionArgs_OptionAddArg("foo"))
	t.Run("repeated arg strings", assertRepeatedOptionArgs_OptionAddArg("foo"))
	t.Run("empty arg string", assertEmptyArgStringError_OptionAddArg("", makeError("usage: arg string must not be empty")))
}

func TestNewOption(t *testing.T) {
	t.Run("baseline", assertOption_NewOption(
		[]string{"foo", "bar"},
		"foo",
		makeOption([]string{"foo", "bar"}, "foo"),
	))
	t.Run("single alias", assertOption_NewOption(
		[]string{"foo"},
		"foo",
		makeOption([]string{"foo"}, "foo"),
	))
	t.Run("single repeated alias", assertOption_NewOption(
		[]string{"foo", "foo"},
		"foo",
		makeOption([]string{"foo", "foo"}, "foo"),
	))
	t.Run("multiple repeated aliases", assertOption_NewOption(
		[]string{"foo", "bar", "foo", "bar"},
		"foo",
		makeOption([]string{"foo", "bar", "foo", "bar"}, "foo"),
	))
	t.Run("empty description string", assertOption_NewOption(
		[]string{"foo", "bar"},
		"",
		makeOption([]string{"foo", "bar"}, ""),
	))
	t.Run("nil aliases", assertNoOptionAliasProvidedError_NewOption(
		nil,
		"foo",
		makeError("usage: option must have at least one alias"),
	))
	t.Run("no aliases", assertNoOptionAliasProvidedError_NewOption(
		make([]string, 0),
		"foo",
		makeError("usage: option must have at least one alias"),
	))
	t.Run("single empty alias string", assertEmptyOptionAliasStringError_NewOption(
		[]string{""},
		"foo",
		makeError("usage: alias string must not be empty"),
	))
	t.Run("multiple empty alias strings", assertEmptyOptionAliasStringError_NewOption(
		[]string{"foo", "", "bar", ""},
		"foo",
		makeError("usage: alias string must not be empty"),
	))
}
