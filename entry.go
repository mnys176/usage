package usage

import (
	"errors"
	"sort"
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

func (e *Entry) AddArg(arg string) error {
	if arg == "" {
		return &UsageError{errors.New("arg string must not be empty")}
	}
	e.args = append(e.args, arg)
	return nil
}

func (e *Entry) AddOption(o *Option) error {
	if o == nil {
		return &UsageError{errors.New("no option provided")}
	}
	e.options = append(e.options, *o)
	return nil
}

func (e *Entry) AddEntry(entry *Entry) error {
	return nil
}

func (e *Entry) SetName(name string) error {
	if name == "" {
		return &UsageError{errors.New("name string must not be empty")}
	}
	e.name = name
	return nil
}

func (e Entry) Usage() (string, error) {
	return "", nil
}

func (e Entry) Lookup(lookupPath string) (string, error) {
	return "", nil
}

func NewEntry(name, desc string) (*Entry, error) {
	if name == "" {
		return nil, &UsageError{errors.New("name string must not be empty")}
	}

	tmpl := `foo`

	return &Entry{
		Description: desc,
		Tmpl:        tmpl,
		name:        name,
		args:        make([]string, 0),
		options:     make([]Option, 0),
		children:    make(map[string]*Entry),
	}, nil
}
