package usage

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Usage struct {
	Name    string
	entries map[string]Entry
	options []Option
	args    []string
}

func (u Usage) Entries() []Entry {
	output := make([]Entry, 0)
	for _, v := range u.entries {
		output = append(output, v)
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i].Name < output[j].Name
	})
	return output
}

func (u *Usage) AddEntry(e *Entry) error {
	if e == nil {
		return nilEntryProvidedErr()
	}
	if e.Name == "" {
		return emptyEntryNameStringErr()
	}
	if len(u.args) > 0 {
		return existingArgsErr()
	}
	u.entries[e.Name] = *e
	return nil
}

func (u *Usage) AddOption(o *Option) error {
	if o == nil {
		return nilOptionProvidedErr()
	}
	if len(o.Aliases) == 0 {
		return noOptionAliasProvidedErr()
	}
	for _, alias := range o.Aliases {
		if len(alias) == 0 {
			return emptyOptionAliasStringErr()
		}
	}
	u.options = append(u.options, *o)
	return nil
}

func (u *Usage) AddArg(arg string) error {
	if len(arg) == 0 {
		return emptyArgStringErr()
	}
	if len(u.entries) > 0 {
		return existingEntriesErr()
	}
	u.args = append(u.args, arg)
	return nil
}

func (u Usage) Global() string {
	var summary, summaryExt strings.Builder
	summary.WriteString(u.Name)
	summaryExt.WriteString("\n\nTo learn more about the available options" +
		" for each command, use the --help flag like so:\n\n" + u.Name)
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

		entrySummary := fmt.Sprintf("\n    %s%s", e.Name, args)
		usage.WriteString(entrySummary)
		for _, line := range chopMultipleParagraphs(e.Description, 64) {
			usage.WriteString("\n        " + line)
		}
		usage.WriteString("\n")
	}
	return usage.String()
}

func (u Usage) Lookup(entry string) string {
	return "lookup: " + entry
}

func NewUsage(name string) *Usage {
	return &Usage{
		Name:    name,
		entries: make(map[string]Entry),
		options: make([]Option, 0),
	}
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
