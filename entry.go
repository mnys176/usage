package usage

import (
	"errors"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

type Entry struct {
	Description string
	Tmpl        string
	name        string
	args        []string
	options     []Option
	children    map[string]*Entry
	parent      *Entry
}

func (e Entry) Args() []string {
	return e.args
}

func (e Entry) Options() []Option {
	return e.options
}

func (e Entry) Entries() []Entry {
	output := make([]Entry, 0)
	for _, v := range e.children {
		output = append(output, *v)
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].name < output[j].name
	})
	return output
}

func (e Entry) Name() string {
	return e.name
}

func (e Entry) Ancestry() []string {
	ancestry := []string{e.name}
	for ptr := &e; ptr.parent != nil; ptr = ptr.parent {
		ancestry = append(ancestry, ptr.parent.name)
	}
	return ancestry
}

func (e *Entry) AddArg(arg string) error {
	if len(e.children) > 0 {
		return &UsageError{errors.New("cannot add arg with child entries present")}
	}
	if arg == "" {
		return &UsageError{errors.New("arg string must not be empty")}
	}
	e.args = append(e.args, arg)
	return nil
}

func (e *Entry) AddOption(option *Option) error {
	if option == nil {
		return &UsageError{errors.New("no option provided")}
	}
	if len(option.aliases) == 0 {
		return &UsageError{errors.New("option must have at least one alias")}
	}
	for _, alias := range option.aliases {
		if len(alias) == 0 {
			return &UsageError{errors.New("alias string must not be empty")}
		}
	}
	e.options = append(e.options, *option)
	return nil
}

func (e *Entry) AddEntry(entry *Entry) error {
	if entry == nil {
		return &UsageError{errors.New("no entry provided")}
	}
	if entry.name == "" {
		return &UsageError{errors.New("name string must not be empty")}
	}
	if len(e.args) > 0 {
		return &UsageError{errors.New("cannot add child entry with args present")}
	}
	entry.parent = e
	e.children[entry.name] = entry
	return nil
}

func (e *Entry) SetName(name string) error {
	if name == "" {
		return &UsageError{errors.New("name string must not be empty")}
	}
	e.name = name
	return nil
}

func (e Entry) Usage() string {
	fn := template.FuncMap{
		"join":    strings.Join,
		"reverse": reverseAncestryChain,
		"summary": deriveSummaryString,
		"chop":    chopEssay,
	}
	t := template.Must(template.New(e.name).Funcs(fn).Parse(e.Tmpl))
	var b strings.Builder
	t.Execute(&b, e)
	return b.String()
}

func (e Entry) Lookup(lookupPath string) (string, error) {
	return "", nil
}

func NewEntry(name, desc string) (*Entry, error) {
	if name == "" {
		return nil, &UsageError{errors.New("name string must not be empty")}
	}
	tmpl := `Usage:
    {{summary .}}{{if .Entries}}

    To learn more about the available options for each command,
    use the --help flag like so:

    {{.Name}} <command> --help

Commands:{{range $command := .Entries}}
    {{$command.Name}}{{if $command.Args}} {{join $command.Args " "}}{{end}}{{if $command.Description}}
        {{with chop $command.Description 64}}{{join . "\n        "}}{{end}}{{end}}{{end}}{{end}}{{if .Options}}

Options:{{range $option := .Options}}
    {{$option.Usage}}{{end}}{{end}}`
	return &Entry{
		Description: desc,
		Tmpl:        tmpl,
		name:        name,
		args:        make([]string, 0),
		options:     make([]Option, 0),
		children:    make(map[string]*Entry),
	}, nil
}

func chopParagraph(paragraph string, length int) []string {
	paragraph = strings.TrimSpace(paragraph)
	splitter := regexp.MustCompile(`\s+`)
	words := splitter.Split(paragraph, -1)
	lines := make([]string, 0)

	var b strings.Builder
	for _, w := range words {
		if len(w) > length {
			continue
		}
		if b.Len()+len(w) > length {
			lines = append(lines, strings.TrimSpace(b.String()))
			b.Reset()
		}
		b.WriteString(w + " ")
	}
	lines = append(lines, strings.TrimSpace(b.String()))
	return lines
}

func chopEssay(essay string, length int) []string {
	lines := make([]string, 0)
	splitter := regexp.MustCompile("\n+")
	for _, p := range splitter.Split(essay, -1) {
		if len(p) > 0 {
			pLines := chopParagraph(p, length)
			pLines = append(pLines, "")
			lines = append(lines, pLines...)
		}
	}
	if len(lines) == 0 {
		return lines
	}
	return lines[:len(lines)-1]
}

func deriveSummaryString(entry Entry) string {
	var b strings.Builder
	b.WriteString(strings.Join(reverseAncestryChain(entry.Ancestry()), " "))
	if len(entry.children) > 0 {
		b.WriteString(" <command>")
	}
	if len(entry.options) > 0 {
		b.WriteString(" [options]")
	}
	if len(entry.children) > 0 {
		foundArgs := false
		visit(&entry, func(e *Entry) {
			foundArgs = len(e.args) > 0
		})
		if foundArgs {
			b.WriteString(" <args>")
		}
	} else if len(entry.args) > 0 {
		b.WriteString(" " + strings.Join(entry.args, " "))
	}
	return b.String()
}

func visit(entry *Entry, fn func(e *Entry)) {
	fn(entry)
	for _, c := range entry.children {
		visit(c, fn)
	}
}

func reverseAncestryChain(ancestry []string) []string {
	if len(ancestry) == 0 {
		return ancestry
	}
	reversed := make([]string, len(ancestry))
	for i := 0; i <= len(ancestry)/2; i++ {
		reversed[i] = ancestry[len(ancestry)-i-1]
		reversed[len(ancestry)-i-1] = ancestry[i]
	}
	return reversed
}
