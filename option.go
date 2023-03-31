package usage

type Option interface {
	argCollector
	Aliases() []string
	Description() string
}

type defaultOption struct {
	aliases     []string
	description string
	args        []string
}

func (o defaultOption) Aliases() []string {
	return o.aliases
}

func (o defaultOption) Description() string {
	return o.description
}

func (o defaultOption) Args() []string {
	return o.args
}

func (o *defaultOption) AddArg(arg string) error {
	if len(arg) == 0 {
		return emptyArgStringErr()
	}

	exists := false
	for i := range o.args {
		exists = o.args[i] == arg
	}
	if !exists {
		o.args = append(o.args, arg)
	}
	return nil
}

func NewOption(aliases []string, description string) (Option, error) {
	if len(aliases) == 0 {
		return nil, noOptionAliasProvidedErr()
	}
	for _, alias := range aliases {
		if len(alias) == 0 {
			return nil, emptyOptionAliasStringErr()
		}
	}
	return &defaultOption{
		aliases:     aliases,
		description: description,
		args:        make([]string, 0),
	}, nil
}
