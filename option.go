package usage

type option struct {
	Description string
	aliases     []string
	args        []string
}

func (o option) Args() []string {
	return o.args
}

func (o *option) AddArg(arg string) error {
	if arg == "" {
		return emptyArgStringErr()
	}
	o.args = append(o.args, arg)
	return nil
}

func (o *option) SetAliases(aliases []string) error {
	if len(aliases) == 0 {
		return noAliasProvidedErr()
	}
	for _, alias := range aliases {
		if len(alias) == 0 {
			return emptyAliasStringErr()
		}
	}
	o.aliases = aliases
	return nil
}

func NewOption(aliases []string, description string) (*option, error) {
	if len(aliases) == 0 {
		return nil, noAliasProvidedErr()
	}
	for _, alias := range aliases {
		if len(alias) == 0 {
			return nil, emptyAliasStringErr()
		}
	}

	return &option{
		aliases:     aliases,
		Description: description,
		args:        make([]string, 0),
	}, nil
}
