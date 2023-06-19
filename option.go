package main

import "text/template"

type Option struct {
	Description string
	Tmpl        *template.Template
	aliases     []string
	args        []string
}

func (o Option) Args() []string {
	return nil
}

func (o Option) Aliases() []string {
	return nil
}

func (o *Option) AddArg(arg string) error {
	return nil
}

func (o *Option) SetAliases(aliases []string) error {
	return nil
}

func (o Option) Usage() (string, error) {
	return "", nil
}
