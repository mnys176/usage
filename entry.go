package usage

import (
	"strings"
)

type entry struct {
	Description string
	name        string
	args        argSlice
	options     []option
}

func (e entry) Args() []string {
	return e.args
}

func (e entry) Options() []option {
	return e.options
}

func (e *entry) AddArg(arg string) error {
	if arg == "" {
		return emptyArgStringErr()
	}
	e.args = append(e.args, arg)
	return nil
}

func (e *entry) AddOption(o *option) error {
	if o == nil {
		return nilOptionProvidedErr()
	}
	e.options = append(e.options, *o)
	return nil
}

func (e *entry) SetName(name string) error {
	if name == "" {
		return emptyNameStringErr()
	}
	e.name = name
	return nil
}

func (e entry) String() string {
	var entryBuilder strings.Builder
	entryBuilder.WriteString(Indent + e.name)
	if len(e.args) > 0 {
		entryBuilder.WriteString(" " + e.args.String())
	}
	for _, line := range chopMultipleParagraphs(e.Description, 64) {
		entryBuilder.WriteString("\n" + strings.Repeat(Indent, 2) + line)
	}
	return entryBuilder.String()
}

func NewEntry(name, description string) (*entry, error) {
	if name == "" {
		return nil, emptyNameStringErr()
	}
	return &entry{
		name:        name,
		Description: description,
		args:        make([]string, 0),
		options:     make([]option, 0),
	}, nil
}
