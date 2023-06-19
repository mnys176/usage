package main

import (
	"errors"
	"text/template"
)

type Option struct {
	Description string
	Tmpl        *template.Template
	aliases     []string
	args        []string
}

func (o Option) Args() []string {
	return o.args
}

func (o Option) Aliases() []string {
	return nil
}

func (o *Option) AddArg(arg string) error {
	if arg == "" {
		return &UsageError{
			Context: "usage",
			Err:     errors.New("arg string must not be empty"),
		}
	}
	o.args = append(o.args, arg)
	return nil
}

func (o *Option) SetAliases(aliases []string) error {
	return nil
}

func (o Option) Usage() (string, error) {
	return "", nil
}

func NewOption(aliases []string, desc string) (*Option, error) {
	if len(aliases) == 0 {
		return nil, &UsageError{
			Context: "usage",
			Err:     errors.New("option must have at least one alias"),
		}
	}
	for _, alias := range aliases {
		if len(alias) == 0 {
			return nil, &UsageError{
				Context: "usage",
				Err:     errors.New("alias string must not be empty"),
			}
		}
	}

	return &Option{
		aliases:     aliases,
		Description: desc,
		args:        make([]string, 0),
	}, nil
}
