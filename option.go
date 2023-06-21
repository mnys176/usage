package main

import (
	"errors"
	"regexp"
	"strings"
	"text/template"
)

type Option struct {
	Description string
	Tmpl        string
	aliases     []string
	args        []string
}

func (o Option) Args() []string {
	return o.args
}

func (o Option) Aliases() []string {
	return o.aliases
}

func (o *Option) AddArg(arg string) error {
	if arg == "" {
		return &UsageError{errors.New("arg string must not be empty")}
	}
	o.args = append(o.args, arg)
	return nil
}

func (o *Option) SetAliases(aliases []string) error {
	if len(aliases) == 0 {
		return &UsageError{errors.New("option must have at least one alias")}
	}
	for _, alias := range aliases {
		if len(alias) == 0 {
			return &UsageError{errors.New("alias string must not be empty")}
		}
	}
	o.aliases = aliases
	return nil
}

func (o Option) Usage() (string, error) {
	fn := template.FuncMap{
		"join": strings.Join,
		"chop": chop,
	}
	t := template.Must(template.New(strings.Join(o.aliases, "/")).Funcs(fn).Parse(o.Tmpl))
	var b strings.Builder
	err := t.Execute(&b, o)
	return b.String(), err
}

func NewOption(aliases []string, desc string) (*Option, error) {
	if len(aliases) == 0 {
		return nil, &UsageError{errors.New("option must have at least one alias")}
	}
	for _, alias := range aliases {
		if len(alias) == 0 {
			return nil, &UsageError{errors.New("alias string must not be empty")}
		}
	}

	tmpl := `    {{join .Aliases ", "}}{{if .Args}} {{join .Args " "}}{{end}}
        {{with chop .Description 84}}{{join . "        "}}{{end}}
`

	return &Option{
		Description: desc,
		Tmpl:        tmpl,
		aliases:     aliases,
		args:        make([]string, 0),
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

func chop(str string, length int) []string {
	return chopEssay(str, length)
}
