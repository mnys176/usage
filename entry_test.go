package usage

import (
	"testing"
)

func TestDefaultEntryAddArg(t *testing.T) {
	t.Run("baseline", assertEntryArgs_DefaultEntryAddArg("foo"))
	t.Run("repeated arg strings", assertRepeatedEntryArgs_DefaultEntryAddArg("foo"))
	t.Run("empty arg string", assertEmptyArgStringError_DefaultEntryAddArg("", makeError("usage: arg string must not be empty")))
}

func TestDefaultEntryAddOption(t *testing.T) {
	t.Run("baseline", assertEntryOptions_DefaultEntryAddOption(makeDefaultOption([]string{"foo", "bar"}, "foo")))
	t.Run("repeated options", assertRepeatedEntryOptions_DefaultEntryAddOption(makeDefaultOption([]string{"foo", "bar"}, "foo")))
	t.Run("nil option aliases", assertNoOptionAliasProvidedError_DefaultEntryAddOption(
		makeDefaultOption(nil, "foo"),
		makeError("usage: option must have at least one alias"),
	))
	t.Run("no option aliases", assertNoOptionAliasProvidedError_DefaultEntryAddOption(
		makeDefaultOption(make([]string, 0), "foo"),
		makeError("usage: option must have at least one alias"),
	))
	t.Run("single empty option alias string", assertEmptyOptionAliasStringError_DefaultEntryAddOption(
		makeDefaultOption([]string{""}, "foo"),
		makeError("usage: alias string must not be empty"),
	))
	t.Run("multiple empty option alias strings", assertEmptyOptionAliasStringError_DefaultEntryAddOption(
		makeDefaultOption([]string{"foo", "", "bar", ""}, "foo"),
		makeError("usage: alias string must not be empty"),
	))
}

func TestNewEntry(t *testing.T) {
	t.Run("baseline", assertEntry_NewEntry(
		"foo",
		"foo",
		makeDefaultEntry("foo", "foo"),
	))
	t.Run("empty description string", assertEntry_NewEntry(
		"foo",
		"",
		makeDefaultEntry("foo", ""),
	))
	t.Run("empty name string", assertEmptyNameStringError_NewEntry(
		"",
		"foo",
		makeError("usage: name string must not be empty"),
	))
}
