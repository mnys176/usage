package usage

import "errors"

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
	return nil
}

func Usage() string {
	return ""
}

func Lookup(lookup string) string {
	return ""
}

func checkInit() {
	if global == nil {
		panic(&UsageError{errors.New("global usage not initialized")})
	}
}
