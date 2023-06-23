package usage

func Init(name string) error {
	return nil
}

func Args() []string {
	return nil
}

func Options() []Option {
	return nil
}

func Entries() []Entry {
	return nil
}

func AddArg(arg string) error {
	return nil
}

func AddOption(o *Option) error {
	return nil
}

func AddEntry(e *Entry) error {
	return nil
}

func SetName(name string) error {
	return nil
}

func Usage() (string, error) {
	return "", nil
}

func Lookup(lookupPath string) (string, error) {
	return "", nil
}
