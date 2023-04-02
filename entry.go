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

func (e *defaultEntry) AddArg(arg string) (err error) {
	if len(arg) == 0 {
		err = emptyArgStringErr()
		return
	}
	e.args = append(e.args, arg)
	return
}

func (e defaultEntry) Options() []Option {
	return e.options
}

func (e *defaultEntry) AddOption(o Option) (err error) {
	if len(o.Aliases()) == 0 {
		err = noOptionAliasProvidedErr()
		return
	}
	for _, alias := range o.Aliases() {
		if len(alias) == 0 {
			err = emptyOptionAliasStringErr()
			return
		}
	}
	e.options = append(e.options, o)
	return
}

func NewEntry(name, description string) (e Entry, err error) {
	if name == "" {
		err = emptyEntryNameStringErr()
		return
	}
	e = &defaultEntry{
		name:        name,
		description: description,
		args:        make([]string, 0),
		options:     make([]Option, 0),
	}
	return
}
