package usage

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type EntryOption interface {
	Aliases() []string
	Description() string
	Args() []string
}

type Entry interface {
	Name() string
	Description() string
	Args() []string
	Options() []EntryOption
	AddOption(EntryOption)
}

type Usage struct {
	Name    string
	entries map[string]Entry
	options []EntryOption
}

func (u *Usage) AddEntry(e Entry) error {
	if e == nil {
		return errors.New("no entry provided")
	}
	if e.Name() == "" {
		return errors.New("entry must have a name")
	}
	u.entries[e.Name()] = e
	return nil
}

func (u *Usage) AddOption(o EntryOption) error {
	if o == nil {
		return errors.New("no option provided")
	}
	u.options = append(u.options, o)
	return nil
}

func (u Usage) Global() string {
	var usage strings.Builder
	usage.WriteString(`Usage:
    freeformgen <command> [options] <args>

    To learn more about the available options for each command, use the
    --help flag like so:

    freeformgen <command> --help

Commands:`)

	for _, e := range u.entries {
		var b strings.Builder
		for _, arg := range e.Args() {
			b.WriteString("<" + arg + "> ")
		}
		args := strings.TrimSuffix(b.String(), " ")

		summary := fmt.Sprintf("\n    %s %s", e.Name(), args)
		usage.WriteString(summary)
		for _, line := range blockText(e.Description(), blockTextSize) {
			usage.WriteString("\n        " + line)
		}
		usage.WriteString("\n")
	}
	return usage.String()
}

func (u Usage) Lookup(entry string) string {
	return "lookup: " + entry
}

func New(name string) *Usage {
	return &Usage{
		Name:    name,
		entries: make(map[string]Entry),
		options: make([]EntryOption, 0),
	}
}

const blockTextSize int = 64

func blockText(str string, length int) []string {
	chopLines := func(s string) []string {
		re := regexp.MustCompile(" +")
		words := re.Split(s, -1)
		lines := make([]string, 0)
		var lineBuilder strings.Builder
		for _, w := range words {
			if len(w) > length {
				continue
			}
			if lineBuilder.Len()+len(w) > length {
				lines = append(lines, strings.TrimSpace(lineBuilder.String()))
				lineBuilder.Reset()
			}
			lineBuilder.WriteString(w + " ")
		}
		lines = append(lines, strings.TrimSpace(lineBuilder.String()))
		return lines
	}

	allLines := make([]string, 0)
	re := regexp.MustCompile("\n+")
	for _, paragraph := range re.Split(str, -1) {
		if len(paragraph) > 0 {
			lines := chopLines(paragraph)
			lines = append(lines, "")
			allLines = append(allLines, lines...)
		}
	}
	return allLines[:len(allLines)-1]
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
