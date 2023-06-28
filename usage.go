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
	return nil
}

func AddOption(optopm *Option) error {
	return nil
}

func AddEntry(e *Entry) error {
	return nil
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
