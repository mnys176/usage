package main

import "errors"

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
	return nil
}

func (e Entry) Entries() []Entry {
	return nil
}

func (e Entry) Name() string {
	return e.name
}

func (e *Entry) AddArg(arg string) error {
	if arg == "" {
		return &UsageError{errors.New("arg string must not be empty")}
	}
	e.args = append(e.args, arg)
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
