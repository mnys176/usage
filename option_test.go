package usage

import (
	"errors"
	"testing"
)

func TestDefaultOptionAliases(t *testing.T) {
	aliases := []string{"foo", "bar"}
	option := &defaultOption{aliases: aliases}
	results := option.Aliases()
	if len(results) != len(aliases) {
		t.Fatalf("%d aliases returned but wanted %d", len(results), len(aliases))
	}
	for i := range aliases {
		if results[i] != aliases[i] {
			t.Errorf("alias %q should be %q", results[i], aliases[i])
		}
	}
}

func TestDefaultOptionDescription(t *testing.T) {
	description := "foo"
	option := &defaultOption{description: description}
	if got := option.Description(); got != description {
		t.Errorf("desctiption should be %q but got %q", description, got)
	}
}

func TestDefaultOptionArgs(t *testing.T) {
	args := []string{"foo", "bar"}
	option := &defaultOption{args: args}
	results := option.Args()
	if len(results) != len(args) {
		t.Fatalf("%d aliases returned but wanted %d", len(results), len(args))
	}
	for i := range args {
		if results[i] != args[i] {
			t.Errorf("arg %q should be %q", results[i], args[i])
		}
	}
}

func TestDefaultOptionAddArg(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		err  error
	}{
		{
			name: "default",
			arg:  "foo",
		},
		{
			name: "empty arg string",
			err:  emptyArgStringErr(),
		},
		{
			name: "arg already exists",
			arg:  "foo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			option := &defaultOption{aliases: []string{"foo"}}
			err := option.AddArg(tc.arg)

			if tc.name == "empty arg string" {
				if err == nil {
					t.Fatal("no error returned with an empty arg provided")
				}
				if !errors.Is(err, tc.err) {
					t.Errorf("got %q error but wanted %q", err, tc.err)
				}
				return
			}

			if tc.name == "arg already exists" {
				// Should not make any changes.
				option.AddArg(tc.arg)
			}

			if len(option.args) != 1 {
				t.Errorf("%d args returned but wanted 1", len(option.args))
			}
		})
	}
}

func TestNewOption(t *testing.T) {
	tests := []struct {
		name        string
		aliases     []string
		description string
		err         error
	}{
		{
			name:    "default",
			aliases: []string{"foo"},
		},
		{
			name:    "no aliases provided",
			aliases: make([]string, 0),
			err:     noOptionAliasProvidedErr(),
		},
		{
			name:    "empty alias string",
			aliases: []string{""},
			err:     emptyOptionAliasStringErr(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results, err := NewOption(tc.aliases, "")

			if tc.name == "no aliases provided" {
				if err == nil {
					t.Fatal("no error returned with no aliases provided")
				}
				if !errors.Is(err, tc.err) {
					t.Errorf("got %q error but wanted %q", err, tc.err)
				}
				return
			}

			if tc.name == "empty alias string" {
				if err == nil {
					t.Fatal("no error returned with an empty alias provided")
				}
				if !errors.Is(err, tc.err) {
					t.Errorf("got %q error but wanted %q", err, tc.err)
				}
				return
			}

			if len(results.Aliases()) != len(tc.aliases) {
				t.Fatalf("option aliases should have length of %d but has %d", len(tc.aliases), len(results.Aliases()))
			}
			for i, alias := range results.Aliases() {
				if alias != tc.aliases[i] {
					t.Errorf("option alias should be %q but is %q", tc.aliases[i], alias)
				}
			}

			if results.Description() != tc.description {
				t.Errorf("option desctiption should be %q but is %q", tc.description, results.Description())
			}

			if results.Args() == nil || len(results.Args()) != 0 {
				t.Error("option arg slice not initialized to empty")
			}
		})
	}
}
