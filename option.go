package usage

type option struct {
	Aliases     []string
	Description string
	args        []string
}

func (o option) Args() []string {
	return o.args
}

func (o *option) AddArg(arg string) error {
	if len(arg) == 0 {
		return emptyArgStringErr()
	}
	o.args = append(o.args, arg)
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
		Aliases:     aliases,
		Description: description,
		args:        make([]string, 0),
	}, nil
}
