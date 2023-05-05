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

func newArgSlice(argString string) argSlice {
	if argString == "" {
		return make(argSlice, 0)
	}
	return argSlice(strings.Split(argString[1:len(argString)-1], "> <"))
}
