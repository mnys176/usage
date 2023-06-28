package usage

import (
	"errors"
	"text/template"
)

var global *Entry

func Init(name string) error {
	glob, err := NewEntry(name, "")
	if err != nil {
		return err
	}
	global = glob
	return nil
}

func Args() []string {
	checkInit()
	return global.Args()
}

func Options() []Option {
	checkInit()
	return global.Options()
}

func Entries() []Entry {
	checkInit()
	return global.Entries()
}

func AddArg(arg string) error {
	checkInit()
	return global.AddArg(arg)
}

func AddOption(option *Option) error {
	checkInit()
	return global.AddOption(option)
}

func AddEntry(entry *Entry) error {
	checkInit()
	return global.AddEntry(entry)
}

func SetName(name string) error {
	checkInit()
	return global.SetName(name)
}

func Usage() string {
	checkInit()
	return global.Usage()
}

func Lookup(lookup string) string {
	checkInit()
	return global.Lookup(lookup)
}

func SetEntryTemplate(tmpl *template.Template) {
	checkInit()
	visit(global, func(e *Entry) {
		e.setTemplate(tmpl)
	})
}

func checkInit() {
	if global == nil {
		panic(&UsageError{errors.New("global usage not initialized")})
	}
}
