package usage

import "strings"

type option struct {
	Description string
	aliases     []string
	args        argSlice
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

func (o option) String() string {
	var optionBuilder, aliasBuilder strings.Builder
	for _, alias := range o.aliases {
		if len(alias) == 1 {
			aliasBuilder.WriteString("-" + alias)
		} else {
			aliasBuilder.WriteString("--" + alias)
		}
		aliasBuilder.WriteString(", ")
	}
	optionBuilder.WriteString(Indent + aliasBuilder.String()[:len(aliasBuilder.String())-2])

	if len(o.args) > 0 {
		optionBuilder.WriteString(" " + o.args.String())
	}
	for _, line := range chopMultipleParagraphs(o.Description, 64) {
		optionBuilder.WriteString("\n" + strings.Repeat(Indent, 2) + line)
	}
	return optionBuilder.String()
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
