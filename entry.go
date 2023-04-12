package usage

type Entry struct {
	Name        string
	Description string
	args        []string
	options     []Option
}

func (e Entry) Args() []string {
	return e.args
}

func (e *Entry) AddArg(arg string) error {
	if len(arg) == 0 {
		return emptyArgStringErr()
	}
	e.args = append(e.args, arg)
	return nil
}

func (e Entry) Options() []Option {
	return e.options
}

func (e *Entry) AddOption(o *Option) error {
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

	e.options = append(e.options, *o)
	return nil
}

func NewEntry(name, description string) (*Entry, error) {
	if name == "" {
		return nil, emptyEntryNameStringErr()
	}
	return &Entry{
		Name:        name,
		Description: description,
		args:        make([]string, 0),
		options:     make([]Option, 0),
	}, nil
}
