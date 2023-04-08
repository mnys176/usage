package usage

import "errors"

func makeError(s string) error {
	return errors.New(s)
}

func makeOption(aliases []string, description string) Option {
	return &defaultOption{
		aliases:     aliases,
		description: description,
		args:        make([]string, 0),
	}
}
func makeEntry(name string, description string) Entry {
	return &defaultEntry{
		name:        name,
		description: description,
		options:     make([]Option, 0),
		args:        make([]string, 0),
	}
}
