package usage

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func assertError(t *testing.T, got, want error) {
	if !errors.Is(got, want) {
		t.Errorf("got %q error but wanted %q", got, want)
	}
}

func assertNilError(t *testing.T, got error) {
	if got != nil {
		t.Errorf("got %q error but should be nil", got)
	}
}

func assertEmptyArgStringError(t *testing.T, got, want error) {
	if got == nil {
		t.Fatal("no error returned with an empty arg string")
	}
	assertError(t, got, want)
}

func assertNoAliasesError(t *testing.T, got, want error) {
	if got == nil {
		t.Fatal("no error returned with no provided aliases")
	}
	assertError(t, got, want)
}

func assertEmptyAliasStringError(t *testing.T, got, want error) {
	if got == nil {
		t.Fatal("no error returned with an empty alias string")
	}
	assertError(t, got, want)
}

func assertEmptyNameStringError(t *testing.T, got, want error) {
	if got == nil {
		t.Fatal("no error returned with an empty name string")
	}
	assertError(t, got, want)
}

func assertName(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("name is %q but should be %q", got, want)
	}
}

func assertDescription(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("description is %q but should be %q", got, want)
	}
}

func assertUsage(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("usage is %q but should be %q", got, want)
	}
}

func assertAliases(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d aliases returned but wanted %d", len(got), len(want))
	}
	for i, gotAlias := range got {
		if gotAlias != want[i] {
			t.Errorf("alias is %q but should be %q", gotAlias, want[i])
		}
	}
}

func assertArgs(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d args returned but wanted %d", len(got), len(want))
	}
	for i, gotArg := range got {
		if gotArg != want[i] {
			t.Errorf("arg is %q but should be %q", gotArg, want[i])
		}
	}
}

func assertOption(t *testing.T, got, want *Option) {
	assertAliases(t, got.aliases, want.aliases)
	assertDescription(t, got.Description, want.Description)
	if got.args == nil || len(got.args) != 0 {
		t.Error("args not initialized to an empty slice")
	}
	if got.Tmpl == "" {
		t.Error("template not initialized to a template string")
	}
	got.Usage()
}

func assertNilOption(t *testing.T, got *Option) {
	if got != nil {
		t.Errorf("got %+v option but should be nil", got)
	}
}

func assertErrorString(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("error is %q but should be %q", got, want)
	}
}

func assertErrorEquality(t *testing.T, got, want bool) {
	if got != want {
		t.Errorf("error equality is %t but should be %t", got, want)
	}
}

func stringToOption(str string) *Option {
	indent := "    "
	aliasesAndArgs, choppedDescription, _ := strings.Cut(str, "\n"+strings.Repeat(indent, 2))
	aliasesAndArgs = strings.TrimPrefix(aliasesAndArgs, indent)

	aliasesString, argsString, _ := strings.Cut(aliasesAndArgs, " | ")
	aliases := strings.Split(aliasesString, " ")
	args := make([]string, 0)
	if len(argsString) > 0 {
		args = strings.Split(argsString, " ")
	}

	var description strings.Builder
	for _, line := range strings.Split(choppedDescription, "\n"+strings.Repeat(indent, 2)) {
		if line == "" {
			description.WriteString("\n\n")
		} else {
			description.WriteString(" " + line)
		}
	}

	return &Option{
		Description: strings.Replace(strings.TrimSpace(description.String()), "> <", fmt.Sprintf("> %s <", strings.Repeat("a", 72)), 1),
		Tmpl: `    {{join .Aliases " "}}{{if .Args}} | {{join .Args " "}}{{end}}{{if .Description}}
        {{with chop .Description 64}}{{join . "\n        "}}{{end}}{{end}}`,
		aliases: aliases,
		args:    args,
	}
}
