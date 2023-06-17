package usage

import "strings"

type ArgSlice []string

func (args ArgSlice) String() string {
	if len(args) == 0 {
		return ""
	}
	var argsBuilder strings.Builder
	for _, arg := range args {
		argsBuilder.WriteString(" <" + arg + ">")
	}
	return argsBuilder.String()[1:]
}

func NewArgSlice(argString string) ArgSlice {
	if argString == "" {
		return make(ArgSlice, 0)
	}
	return ArgSlice(strings.Split(argString[1:len(argString)-1], "> <"))
}
