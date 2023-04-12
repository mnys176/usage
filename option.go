package usage

type Option struct {
	Aliases     []string
	Description string
	args        []string
}

func (o Option) Args() []string {
	return o.args
}

func (o *Option) AddArg(arg string) error {
	if len(arg) == 0 {
		return emptyArgStringErr()
	}
	o.args = append(o.args, arg)
	return nil
}

func NewOption(aliases []string, description string) (*Option, error) {
	if len(aliases) == 0 {
		return nil, noOptionAliasProvidedErr()
	}
	for _, alias := range aliases {
		if len(alias) == 0 {
			return nil, emptyOptionAliasStringErr()
		}
	}

	return &Option{
		Aliases:     aliases,
		Description: description,
		args:        make([]string, 0),
	}, nil
}
