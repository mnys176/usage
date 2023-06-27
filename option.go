package usage

import (
	"errors"
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
		"chop": chopEssay,
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
	tmpl := `{{join .Aliases ", "}}{{if .Args}} {{join .Args " "}}{{end}}{{if .Description}}
        {{with chop .Description 64}}{{join . "\n        "}}{{end}}{{end}}`
	return &Option{
		Description: desc,
		Tmpl:        tmpl,
		aliases:     aliases,
		args:        make([]string, 0),
	}, nil
}
