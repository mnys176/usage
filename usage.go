package usage

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const Indent string = "    "

type usage struct {
	name    string
	entries map[string]entry
	options []option
	args    argSlice
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

func (u usage) Usage() string {
	hasEntries, hasOptions, hasArgs := len(u.entries) > 0, len(u.options) > 0, len(u.args) > 0

	var usage strings.Builder
	usage.WriteString("Usage:\n" + Indent)

	var summary strings.Builder
	summary.WriteString(u.name)
	if hasEntries {
		summary.WriteString(" <command>")
	}
	if hasOptions {
		summary.WriteString(" [options]")
	}
	if hasEntries {
		for _, e := range u.entries {
			if len(e.args) > 0 {
				summary.WriteString(" <args>")
				break
			}
		}
	} else if hasArgs {
		summary.WriteString(" " + u.args.String())
	}

	if hasEntries {
		summary.WriteRune('\n')
		extension := "To learn more about the available options" +
			" for each command, use the --help flag like so:"
		for _, line := range chopSingleParagraph(extension, 68) {
			summary.WriteString("\n" + Indent + line)
		}
		summary.WriteString(fmt.Sprintf("\n\n%s%s <command> --help", Indent, u.name))
	}
	usage.WriteString(summary.String() + "\n")

	if hasOptions {
		usage.WriteString("\nOptions:")
		for _, o := range u.options {
			usage.WriteString(fmt.Sprintf("\n%s\n", o.String()))
		}
	}

	if hasEntries {
		usage.WriteString("\nCommands:")
		for _, e := range u.Entries() {
			usage.WriteString(fmt.Sprintf("\n%s\n", e.String()))
		}
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
	if len(lines) == 0 {
		return lines
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
// 		b.WriteString("Usa%sfreeformgen " + e.Nam Indent,e() + " ")
// 		b.WriteString(args.String())
// 	}
// 	return ""
// }
