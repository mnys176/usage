package usage

import "errors"

func makeError(s string) error {
	return errors.New(s)
}

func makeOption(aliases []string, description string) *Option {
	return &Option{
		Aliases:     aliases,
		Description: description,
		args:        make([]string, 0),
	}
}
func makeEntry(name string, description string) *Entry {
	return &Entry{
		Name:        name,
		Description: description,
		options:     make([]Option, 0),
		args:        make([]string, 0),
	}
}
