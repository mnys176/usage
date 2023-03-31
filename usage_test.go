package main

import (
	"errors"
	"regexp"
	"sort"
	"testing"
)

func TestUsageEntries(t *testing.T) {
	usage := &Usage{entries: make(map[string]Entry)}
	entry1, entry2 := &defaultEntry{name: "foo"}, &defaultEntry{name: "bar"}
	usage.entries[entry1.name] = entry1
	usage.entries[entry2.name] = entry2

	got := usage.Entries()
	if len(got) != len(usage.entries) {
		t.Errorf("expected %d entries but got %d", len(usage.entries), len(got))
		return
	}
	sorted := sort.SliceIsSorted(got, func(i, j int) bool {
		return got[i].Name() < got[j].Name()
	})
	if !sorted {
		t.Error("entries are not sorted alphabetically")
	}
}

func TestUsageAddEntry(t *testing.T) {
	tests := []struct {
		name              string
		e                 Entry
		missingEntryErr   error
		emptyEntryNameErr error
		existingArgsErr   error
	}{
		{
			name: "normal",
			e:    &defaultEntry{name: "foo"},
		},
		{
			name:            "missing entry",
			missingEntryErr: errors.New("no entry provided"),
		},
		{
			name:              "empty entry name",
			e:                 &defaultEntry{},
			emptyEntryNameErr: errors.New("entry name must not be empty"),
		},
		{
			name:            "existing global args",
			e:               &defaultEntry{name: "foo"},
			existingArgsErr: errors.New("cannot use subcommands with global args"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			usage := &Usage{entries: make(map[string]Entry)}

			if tc.name == "missing entry" {
				err := usage.AddEntry(tc.e)
				if err == nil {
					t.Error("no error returned with a missing entry")
				} else if err.Error() != tc.missingEntryErr.Error() {
					t.Errorf("got %q error but wanted %q", err, tc.missingEntryErr)
				}
				return
			}

			if tc.name == "empty entry name" {
				err := usage.AddEntry(tc.e)
				if err == nil {
					t.Error("no error returned with an empty entry name string")
				} else if err.Error() != tc.emptyEntryNameErr.Error() {
					t.Errorf("got %q error but wanted %q", err, tc.emptyEntryNameErr)
				}
				return
			}

			if tc.name == "existing global args" {
				usage.args = append(usage.args, "foo")
				err := usage.AddEntry(tc.e)
				if err == nil {
					t.Error("no error returned with a existing global args")
				} else if err.Error() != tc.existingArgsErr.Error() {
					t.Errorf("got %q error but wanted %q", err, tc.existingArgsErr)
				}
				return
			}

			usage.AddEntry(tc.e)
			if l := len(usage.entries); l != 1 {
				t.Errorf("usage entries should have 1 entry but has %d", l)
			}
		})
	}
}

func TestUsageAddOption(t *testing.T) {
	tests := []struct {
		name                string
		o                   Option
		missingOptionErr    error
		missingAliasErr     error
		emptyOptionAliasErr error
	}{
		{
			name: "normal",
			o:    &defaultEntryOption{aliases: []string{"foo"}},
		},
		{
			name:             "missing option",
			missingOptionErr: errors.New("no option provided"),
		},
		{
			name:            "missing alias",
			o:               &defaultEntryOption{aliases: make([]string, 0)},
			missingAliasErr: errors.New("option must have at least one alias"),
		},
		{
			name:                "empty alias",
			o:                   &defaultEntryOption{aliases: []string{""}},
			emptyOptionAliasErr: errors.New("alias string must not be empty"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			usage := &Usage{options: make([]Option, 0)}

			if tc.name == "missing option" {
				err := usage.AddOption(tc.o)
				if err == nil {
					t.Error("no error returned with a missing option")
				} else if err.Error() != tc.missingOptionErr.Error() {
					t.Errorf("got %q error but wanted %q", err, tc.missingOptionErr)
				}
				return
			}

			if tc.name == "missing alias" {
				err := usage.AddOption(tc.o)
				if err == nil {
					t.Error("no error returned with zero aliases")
				} else if err.Error() != tc.missingAliasErr.Error() {
					t.Errorf("got %q error but wanted %q", err, tc.missingAliasErr)
				}
				return
			}

			if tc.name == "empty alias" {
				err := usage.AddOption(tc.o)
				if err == nil {
					t.Error("no error returned with an empty alias string")
				} else if err.Error() != tc.emptyOptionAliasErr.Error() {
					t.Errorf("got %q error but wanted %q", err, tc.emptyOptionAliasErr)
				}
				return
			}

			usage.AddOption(tc.o)
			if l := len(usage.options); l != 1 {
				t.Errorf("usage options should have 1 option but has %d", l)
			}
		})
	}
}

func TestUsageAddArg(t *testing.T) {
	tests := []struct {
		name               string
		arg                string
		emptyArgErr        error
		existingEntriesErr error
	}{
		{
			name: "normal",
			arg:  "foo",
		},
		{
			name:        "empty arg",
			emptyArgErr: errors.New("arg string must not be empty"),
		},
		{
			name:               "existing entries",
			arg:                "foo",
			existingEntriesErr: errors.New("cannot use global args with subcommands"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			usage := &Usage{args: make([]string, 0)}

			if tc.name == "empty arg" {
				err := usage.AddArg(tc.arg)
				if err == nil {
					t.Error("no error returned with an empty arg")
				} else if err.Error() != tc.emptyArgErr.Error() {
					t.Errorf("got %q error but wanted %q", err, tc.emptyArgErr)
				}
				return
			}

			if tc.name == "existing entries" {
				usage.entries = make(map[string]Entry)
				usage.entries[""] = &defaultEntry{}
				err := usage.AddArg(tc.arg)
				if err == nil {
					t.Error("no error returned with existing subcommands")
				} else if err.Error() != tc.existingEntriesErr.Error() {
					t.Errorf("got %q error but wanted %q", err, tc.existingEntriesErr)
				}
				return
			}

			usage.AddArg(tc.arg)
			if l := len(usage.args); l != 1 {
				t.Errorf("usage args should have 1 arg but has %d", l)
			}
		})
	}
}

func TestNewUsage(t *testing.T) {
	name := "foo"
	got := NewUsage(name)
	if got.Name != name {
		t.Errorf("created usage name is %q but should be %q", got.Name, name)
	}
}

func TestChopSingleParagraph(t *testing.T) {
	tests := []struct {
		name   string
		p      string
		length int
	}{
		{
			name: "normal",
			p: "This is just a couple sentences for" +
				" the `chopSingleParagraph` test. Note" +
				" the use of the two backticks (`) and" +
				" even the parantheses to show how unique" +
				" characters provided some useful edge" +
				" cases. Hopefully this will be into lines" +
				" no longer than 32 characters.",
			length: 32,
		},
		{
			name:   "short sentence",
			p:      "This is a short sentence.",
			length: 32,
		},
		{
			name: "long word",
			p: "There should be nothing in between these" +
				" arrows -> reallyreallyreallyreallylongword <-" +
				" because the word is too long",
			length: 30,
		},
		{
			name:   "empty paragraph",
			length: 32,
		},
		{
			name:   "lots of whitespace",
			p:      "             	\t\tThis \n\n\t is a short sentence.\t\t       ",
			length: 32,
		},
		{
			name: "no viable words",
			p: "The length is set to 0, so there should be" +
				" a single empty string returned in the slice.",
		},
		{
			name:   "negative length",
			length: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "negative length" {
				defer func() {
					if r := recover(); r == nil {
						t.Error("no panic with negative length")
					}
				}()
			}

			pattern := regexp.MustCompile(`^(((\S+ )+)?\S+)?$`)
			got := chopSingleParagraph(tc.p, tc.length)
			for _, line := range got {
				if !pattern.MatchString(line) {
					t.Errorf("line %q has excess whitespace", line)
					return
				}
				if len(line) > tc.length {
					t.Errorf(
						"line %q has a length of %d which is greater than %d",
						line,
						len(line),
						tc.length,
					)
					return
				}
			}
		})
	}
}

func TestChopMultipleParagraph(t *testing.T) {
	tests := []struct {
		name   string
		ps     string
		length int
	}{
		{
			name: "normal",
			ps: "This is just a couple sentences for" +
				" the `chopMultipleParagraph` test. Note" +
				" the use of the two backticks (`) and" +
				" even the parantheses to show how unique" +
				" characters provided some useful edge" +
				" cases. Hopefully this will be into lines" +
				" no longer than 32 characters.\n\nThis" +
				" is just another paragraph that serves to" +
				" test the \"multiple\" part of the new" +
				" function because otherwise everything is" +
				" the same.",
			length: 32,
		},
		{
			name: "lots of whitespace",
			ps: "\n\n\t\t   This is    the first   paragraph.\n" +
				" \nThis is the second paragraph.  \t\t\n\n\n\n" +
				" This is the third and final paragraph.\n\t\n\t",
			length: 32,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pattern := regexp.MustCompile(`^((((((\S+ )+)?\S+)\n\n)+)?(((\S+ )+)?\S+))?$`)
			got := chopMultipleParagraphs(tc.ps, tc.length)
			for _, line := range got {
				if !pattern.MatchString(line) {
					t.Errorf("line %q has excess whitespace", line)
					return
				}
				if len(line) > tc.length {
					t.Errorf(
						"line %q has a length of %d which is greater than %d",
						line,
						len(line),
						tc.length,
					)
					return
				}
			}
		})
	}
}
