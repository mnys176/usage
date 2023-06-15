package usage

import (
	"fmt"
	"strings"
)

type Entry struct {
	Description string
	name        string
	args        argSlice
	options     []Option
}

func (e Entry) Args() []string {
	return e.args
}

func (e Entry) Options() []Option {
	return e.options
}

func (e *Entry) AddArg(arg string) error {
	if arg == "" {
		return emptyArgStringErr()
	}
	e.args = append(e.args, arg)
	return nil
}

func (e *Entry) AddOption(o *Option) error {
	if o == nil {
		return nilOptionProvidedErr()
	}
	e.options = append(e.options, *o)
	return nil
}

func (e *Entry) SetName(name string) error {
	if name == "" {
		return emptyNameStringErr()
	}
	e.name = name
	return nil
}

func (e Entry) Usage() string {
	hasOptions, hasArgs := len(e.options) > 0, len(e.args) > 0

	var usage strings.Builder
	usage.WriteString("Usage:\n" + Indent + "%s ")

	var summary strings.Builder
	summary.WriteString(e.name)
	if hasOptions {
		summary.WriteString(" [options]")
	}
	if hasArgs {
		summary.WriteString(" " + e.args.String())
	}
	usage.WriteString(summary.String() + "\n")

	if hasOptions {
		usage.WriteString("\nOptions:")
		for _, o := range e.options {
			usage.WriteString(fmt.Sprintf("\n%s\n", o.String()))
		}
	}
	return usage.String()
}

func (e Entry) String() string {
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

func NewEntry(name, description string) (*Entry, error) {
	if name == "" {
		return nil, emptyNameStringErr()
	}
	return &Entry{
		name:        name,
		Description: description,
		args:        make([]string, 0),
		options:     make([]Option, 0),
	}, nil
}
