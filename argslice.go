package usage

import "strings"

type argSlice []string

func (args argSlice) String() string {
	if len(args) == 0 {
		return ""
	}
	var argsBuilder strings.Builder
	for _, arg := range args {
		argsBuilder.WriteString(" <" + arg + ">")
	}
	return argsBuilder.String()[1:]
}
