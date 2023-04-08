package usage

type Entry interface {
	argCollector
	optionCollector
	Name() string
	Description() string
}

type defaultEntry struct {
	name        string
	description string
	args        []string
	options     []Option
}

func (e defaultEntry) Name() string {
	return e.name
}

func (e defaultEntry) Description() string {
	return e.description
}

func (e defaultEntry) Args() []string {
	return e.args
}

func (e *defaultEntry) AddArg(arg string) error {
	if len(arg) == 0 {
		return emptyArgStringErr()
	}
	e.args = append(e.args, arg)
	return nil
}

func (e defaultEntry) Options() []Option {
	return e.options
}

func (e *defaultEntry) AddOption(o Option) error {
	if len(o.Aliases()) == 0 {
		return noOptionAliasProvidedErr()
	}
	for _, alias := range o.Aliases() {
		if len(alias) == 0 {
			return emptyOptionAliasStringErr()
		}
	}

	e.options = append(e.options, o)
	return nil
}

func NewEntry(name, description string) (Entry, error) {
	if name == "" {
		return nil, emptyEntryNameStringErr()
	}
	return &defaultEntry{
		name:        name,
		description: description,
		args:        make([]string, 0),
		options:     make([]Option, 0),
	}, nil
}
