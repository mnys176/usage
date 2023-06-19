package main

import "text/template"

type Entry struct {
	Description string
	Tmpl        *template.Template
	name        string
	args        []string
	options     []Option
	children    map[string]*Entry
	parent      *Entry
}

func (e Entry) Args() []string {
	return nil
}

func (e Entry) Options() []Option {
	return nil
}

func (e Entry) Entries() []Entry {
	return nil
}

func (e Entry) Name() string {
	return ""
}

func (e *Entry) AddArg(arg string) error {
	return nil
}

func (e *Entry) AddOption(option *Option) error {
	return nil
}

func (e *Entry) AddEntry(entry *Entry) error {
	return nil
}

func (e *Entry) SetName(name string) error {
	return nil
}

func (e Entry) Usage() (string, error) {
	return "", nil
}

func (e Entry) Lookup(lookupPath string) (string, error) {
	return "", nil
}
