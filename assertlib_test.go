package usage

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"text/template"
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

func assertNoOptionError(t *testing.T, got, want error) {
	if got == nil {
		t.Fatal("no error returned with no option provided")
	}
	assertError(t, got, want)
}

func assertNoEntryError(t *testing.T, got, want error) {
	if got == nil {
		t.Fatal("no error returned with no entry provided")
	}
	assertError(t, got, want)
}

func assertExistingArgsError(t *testing.T, got, want error) {
	if got == nil {
		t.Fatal("no error returned when an entry is added with existing args")
	}
	assertError(t, got, want)
}

func assertExistingEntriesError(t *testing.T, got, want error) {
	if got == nil {
		t.Fatal("no error returned when an arg is added with existing entries")
	}
	assertError(t, got, want)
}

func assertUninitializedError(t *testing.T, got, want error) {
	if got == nil {
		t.Fatal("no error returned with uninitialized global usage")
	}
	assertError(t, got, want)
}

func assertUninitializedPanic(t *testing.T, want error) {
	r := recover()
	if r == nil {
		t.Fatal("no panic with uninitialized global usage")
	}
	assertUninitializedError(t, r.(error), want)
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

func assertTemplate(t *testing.T, got, want *template.Template) {
	if got != want {
		t.Errorf("template is %+v but should be %+v", got, want)
	}
}

func assertUsage(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("usage is %q but should be %q", got, want)
	}
}

func assertChildren(t *testing.T, got, want map[string]*Entry) {
	if len(got) != len(want) {
		t.Fatalf("%d children returned but wanted %d", len(got), len(want))
	}
	for gotName, gotChild := range got {
		if gotChild != want[gotName] {
			t.Errorf("child %q is %+v but should be %+v", gotName, gotChild, want[gotName])
		}
	}
}

func assertParent(t *testing.T, got, want *Entry) {
	if got != want {
		t.Errorf("parent is %+v but should be %+v", got, want)
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

func assertAncestry(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("%d ancestors returned but wanted %d", len(got), len(want))
	}
	for i, gotAncestor := range got {
		if gotAncestor != want[i] {
			t.Errorf("ancestor is %q but should be %q", gotAncestor, want[i])
		}
	}
}

func assertOption(t *testing.T, got, want *Option) {
	assertDescription(t, got.Description, want.Description)
	assertTemplate(t, got.tmpl, want.tmpl)
	assertAliases(t, got.aliases, want.aliases)
	assertArgs(t, got.args, want.args)
}

func assertDefaultOption(t *testing.T, got, want *Option) {
	assertAliases(t, got.aliases, want.aliases)
	assertDescription(t, got.Description, want.Description)
	if got.args == nil || len(got.args) != 0 {
		t.Error("args not initialized to an empty slice")
	}
	if got.tmpl == nil {
		t.Error("template not initialized to a template object")
	}
}

func assertNilOption(t *testing.T, got *Option) {
	if got != nil {
		t.Errorf("got %+v option but should be nil", got)
	}
}

func assertOptions(t *testing.T, got, want []Option) {
	if len(got) != len(want) {
		t.Fatalf("%d options returned but wanted %d", len(got), len(want))
	}
	for i, gotOption := range got {
		assertOption(t, &gotOption, &want[i])
	}
}

func assertEntry(t *testing.T, got, want *Entry) {
	assertName(t, got.name, want.name)
	assertDescription(t, got.Description, want.Description)
	assertTemplate(t, got.tmpl, want.tmpl)
	assertArgs(t, got.args, want.args)
	assertOptions(t, got.options, want.options)
	assertChildren(t, got.children, want.children)
	assertParent(t, got.parent, want.parent)
}

func assertDefaultEntry(t *testing.T, got, want *Entry) {
	assertName(t, got.name, want.name)
	assertDescription(t, got.Description, want.Description)
	if got.args == nil || len(got.args) != 0 {
		t.Error("args not initialized to an empty slice")
	}
	if got.tmpl == nil {
		t.Error("template not initialized to a template object")
	}
	if got.options == nil || len(got.options) != 0 {
		t.Error("options not initialized to an empty slice")
	}
	if got.children == nil || len(got.children) != 0 {
		t.Error("children not initialized to an empty map")
	}
	if got.parent != nil {
		t.Error("parent not initialized to nil")
	}
}

func assertNilEntry(t *testing.T, got *Entry) {
	if got != nil {
		t.Errorf("got %+v entry but should be nil", got)
	}
}

func assertEntries(t *testing.T, got, want []Entry) {
	if len(got) != len(want) {
		t.Fatalf("%d entries returned but wanted %d", len(got), len(want))
	}
	for i, gotEntry := range got {
		assertEntry(t, &gotEntry, &want[i])
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
	const indent = "    "
	aliasesAndArgsString, choppedDescription, _ := strings.Cut(str, "\n"+indent)
	aliasesString, argsString, _ := strings.Cut(aliasesAndArgsString, " ")

	aliases := strings.Split(aliasesString, ",")
	var descriptionBuilder strings.Builder
	if choppedDescription != "" {
		for _, line := range strings.Split(choppedDescription, "\n"+indent) {
			if line == "" {
				descriptionBuilder.WriteString("\n\n")
			} else {
				descriptionBuilder.WriteString(" " + line)
			}
		}
	}
	description := strings.ReplaceAll(descriptionBuilder.String(), "\n\n ", "\n\n")

	rawTmpl := fmt.Sprintf(`{{join .Aliases ","}}{{if .Args}} <args>{{end}}{{if .Description}}
%s{{with chop .Description 64}}{{join . "\n%s"}}{{end}}{{end}}`, indent, indent)
	fn := template.FuncMap{"join": strings.Join, "chop": chopEssay}
	tmpl := template.Must(template.New("").Funcs(fn).Parse(rawTmpl))

	output := &Option{
		Description: strings.TrimPrefix(description, " "),
		tmpl:        tmpl,
		aliases:     aliases,
	}
	if strings.Contains(argsString, "<args>") {
		output.args = make([]string, 1)
	}
	return output
}

func stringToEntry(str string) *Entry {
	const indent = "    "
	ancestryAndTraitString, choppedDescription, _ := strings.Cut(str, "\n"+indent)
	ancestryString, traitString, _ := strings.Cut(ancestryAndTraitString, " ")

	ancestry := strings.Split(ancestryString, ":")
	entries := make([]Entry, len(ancestry))
	for i := len(ancestry) - 1; i >= 0; i-- {
		entries[i] = Entry{name: ancestry[i]}
	}
	for i := len(entries) - 1; i >= 1; i-- {
		entries[i].parent = &entries[i-1]
	}
	var descriptionBuilder strings.Builder
	if choppedDescription != "" {
		for _, line := range strings.Split(choppedDescription, "\n"+indent) {
			if line == "" {
				descriptionBuilder.WriteString("\n\n")
			} else {
				descriptionBuilder.WriteString(" " + line)
			}
		}
	}
	description := strings.ReplaceAll(descriptionBuilder.String(), "\n\n ", "\n\n")

	rawTmpl := fmt.Sprintf(`{{join (reverse .Ancestry) ":"}}{{if .Options}} [options]{{end}}{{if .Entries}} <command>{{end}}{{if .Args}} <args>{{end}}{{if .Description}}
%s{{with chop .Description 64}}{{join . "\n%s"}}{{end}}{{end}}`, indent, indent)
	fn := template.FuncMap{
		"join":    strings.Join,
		"reverse": reverseAncestryChain,
		"summary": deriveSummaryString,
		"chop":    chopEssay,
	}
	tmpl := template.Must(template.New("").Funcs(fn).Parse(rawTmpl))

	output := &entries[len(entries)-1]
	output.Description = strings.TrimPrefix(description, " ")
	output.tmpl = tmpl
	if strings.Contains(traitString, "<command>") {
		output.children = map[string]*Entry{"foo": {}}
	}
	if strings.Contains(traitString, "[options]") {
		output.options = make([]Option, 1)
	}
	if strings.Contains(traitString, "<args>") {
		output.args = make([]string, 1)
	}
	return output
}
