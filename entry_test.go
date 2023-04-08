package usage

import (
	"testing"
)

func TestDefaultEntryAddArg(t *testing.T) {
	t.Run("baseline", assertDefaultEntryAddArgEntryArgs("foo"))
	t.Run("repeated arg strings", assertDefaultEntryAddArgRepeatedEntryArgs("foo"))
	t.Run("empty arg string", assertDefaultEntryAddArgEmptyArgStringError("", makeError("usage: arg string must not be empty")))
}

func TestDefaultEntryAddOption(t *testing.T) {
	t.Run("baseline", assertDefaultEntryAddOptionEntryOptions(makeDefaultOption([]string{"foo", "bar"}, "foo")))
	t.Run("repeated options", assertDefaultEntryAddOptionRepeatedEntryOptions(makeDefaultOption([]string{"foo", "bar"}, "foo")))
	t.Run("nil option aliases", assertDefaultEntryAddOptionNoOptionAliasProvidedError(
		makeDefaultOption(nil, "foo"),
		makeError("usage: option must have at least one alias"),
	))
	t.Run("no option aliases", assertDefaultEntryAddOptionNoOptionAliasProvidedError(
		makeDefaultOption(make([]string, 0), "foo"),
		makeError("usage: option must have at least one alias"),
	))
	t.Run("single empty option alias string", assertDefaultEntryAddOptionEmptyOptionAliasStringError(
		makeDefaultOption([]string{""}, "foo"),
		makeError("usage: alias string must not be empty"),
	))
	t.Run("multiple empty option alias strings", assertDefaultEntryAddOptionEmptyOptionAliasStringError(
		makeDefaultOption([]string{"foo", "", "bar", ""}, "foo"),
		makeError("usage: alias string must not be empty"),
	))
}

func TestNewEntry(t *testing.T) {
	t.Run("baseline", assertNewEntryEntry(
		"foo",
		"foo",
		makeDefaultEntry("foo", "foo"),
	))
	t.Run("empty description string", assertNewEntryEntry(
		"foo",
		"",
		makeDefaultEntry("foo", ""),
	))
	t.Run("empty name string", assertNewEntryEmptyNameStringError(
		"",
		"foo",
		makeError("usage: name string must not be empty"),
	))
}
