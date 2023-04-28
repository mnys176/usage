package usage

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type usage struct {
	name    string
	entries map[string]entry
	options []option
	args    []string
}

func (u usage) Args() []string {
	return u.args
}

func (u usage) Options() []option {
	return u.options
}

func (u usage) Entries() []entry {
	output := make([]entry, 0)
	for _, v := range u.entries {
		output = append(output, v)
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].name < output[j].name
	})
	return output
}

func (u *usage) AddArg(arg string) error {
	if arg == "" {
		return emptyArgStringErr()
	}
	if len(u.entries) > 0 {
		return existingEntriesErr()
	}
	u.args = append(u.args, arg)
	return nil
}

func (u *usage) AddOption(o *option) error {
	if o == nil {
		return nilOptionProvidedErr()
	}
	u.options = append(u.options, *o)
	return nil
}

func (u *usage) AddEntry(e *entry) error {
	if e == nil {
		return nilEntryProvidedErr()
	}
	if len(u.args) > 0 {
		return existingArgsErr()
	}
	u.entries[e.name] = *e
	return nil
}

func (u usage) Global() string {
	var summary, summaryExt strings.Builder
	summary.WriteString(u.name)
	summaryExt.WriteString("\n\nTo learn more about the available options" +
		" for each command, use the --help flag like so:\n\n" + u.name)
	if len(u.entries) > 0 {
		summary.WriteString(" <command>")
		summaryExt.WriteString(" <command>")
	}
	if len(u.options) > 0 {
		summary.WriteString(" [options]")
	}
	if len(u.args) > 0 {
		summary.WriteString(" <args>")
	}
	summaryExt.WriteString(" --help")

	// If there are no subcommands, this section does not need to
	// be appended.
	if len(u.entries) > 0 {
		summary.WriteString(summaryExt.String())
	}

	var usage strings.Builder
	usage.WriteString("Usage:")
	for _, line := range chopSingleParagraph(summary.String(), 68) {
		usage.WriteString("\n    " + line)
	}
	usage.WriteString("\n")

	// No commands, no command section.
	if len(u.entries) == 0 {
		return usage.String()
	}

	usage.WriteString("\nCommands:")

	for _, e := range u.Entries() {
		var b strings.Builder
		for _, arg := range e.Args() {
			b.WriteString(" <" + arg + ">")
		}
		args := b.String()

		entrySummary := fmt.Sprintf("\n    %s%s", e.name, args)
		usage.WriteString(entrySummary)
		for _, line := range chopMultipleParagraphs(e.Description, 64) {
			usage.WriteString("\n        " + line)
		}
		usage.WriteString("\n")
	}
	return usage.String()
}

func (u usage) Lookup(entry string) string {
	return "lookup: " + entry
}

func (u *usage) SetName(name string) error {
	if name == "" {
		return emptyNameStringErr()
	}
	u.name = name
	return nil
}

func NewUsage(name string) (*usage, error) {
	if name == "" {
		return nil, emptyNameStringErr()
	}
	return &usage{
		name:    name,
		entries: make(map[string]entry),
		options: make([]option, 0),
		args:    make([]string, 0),
	}, nil
}

func chopSingleParagraph(p string, length int) []string {
	if length < 0 {
		panic("length cannot be negative")
	}
	p = strings.TrimSpace(p)
	splitter := regexp.MustCompile(`\s+`)
	words := splitter.Split(p, -1)
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

func chopMultipleParagraphs(ps string, length int) []string {
	lines := make([]string, 0)
	splitter := regexp.MustCompile("\n+")
	for _, p := range splitter.Split(ps, -1) {
		if len(p) > 0 {
			pLines := chopSingleParagraph(p, length)
			pLines = append(pLines, "")
			lines = append(lines, pLines...)
		}
	}
	return lines[:len(lines)-1]
}

// func Lookup(name string) string {
// 	if e, ok := inventory[name]; ok {
// 		var b strings.Builder
// 		for _, name := range e.Names() {
// 			b.WriteString(name + ", ")
// 		}
// 		names := strings.TrimSuffix(b.String(), ", ")

// 		b.Reset()
// 		for _, arg := range e.Args() {
// 			b.WriteString("<" + arg + "> ")
// 		}
// 		args := strings.TrimSuffix(b.String(), " ")

// 		var options string
// 		if len(e.Options()) > 0 {
// 			options = "[options] "
// 		}

// 		summary := fmt.Sprintf("freeformgen %s %s %s", names, options, args)

// 		b.Reset()
// 		b.WriteString("Usage:\n    freeformgen " + e.Name() + " ")
// 		b.WriteString(args.String())
// 	}
// 	return ""
// }
