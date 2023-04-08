package usage

import "testing"

func TestEntryAddArg(t *testing.T) {
	t.Run("baseline", assertEntryArgs_EntryAddArg("foo"))
	t.Run("repeated arg strings", assertRepeatedEntryArgs_EntryAddArg("foo"))
	t.Run("empty arg string", assertEmptyArgStringError_EntryAddArg("", makeError("usage: arg string must not be empty")))
}

func TestEntryAddOption(t *testing.T) {
	t.Run("baseline", assertEntryOptions_EntryAddOption(makeOption([]string{"foo", "bar"}, "foo")))
	t.Run("repeated options", assertRepeatedEntryOptions_EntryAddOption(makeOption([]string{"foo", "bar"}, "foo")))
	t.Run("nil option aliases", assertNoOptionAliasProvidedError_EntryAddOption(
		makeOption(nil, "foo"),
		makeError("usage: option must have at least one alias"),
	))
	t.Run("no option aliases", assertNoOptionAliasProvidedError_EntryAddOption(
		makeOption(make([]string, 0), "foo"),
		makeError("usage: option must have at least one alias"),
	))
	t.Run("single empty option alias string", assertEmptyOptionAliasStringError_EntryAddOption(
		makeOption([]string{""}, "foo"),
		makeError("usage: alias string must not be empty"),
	))
	t.Run("multiple empty option alias strings", assertEmptyOptionAliasStringError_EntryAddOption(
		makeOption([]string{"foo", "", "bar", ""}, "foo"),
		makeError("usage: alias string must not be empty"),
	))
}

func TestNewEntry(t *testing.T) {
	t.Run("baseline", assertEntry_NewEntry(
		"foo",
		"foo",
		makeEntry("foo", "foo"),
	))
	t.Run("empty description string", assertEntry_NewEntry(
		"foo",
		"",
		makeEntry("foo", ""),
	))
	t.Run("empty name string", assertEmptyNameStringError_NewEntry(
		"",
		"foo",
		makeError("usage: name string must not be empty"),
	))
}
