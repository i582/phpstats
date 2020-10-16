package shell

import (
	"strings"
)

type Flag struct {
	Name  string
	Value string
}

type Flags struct {
	Flags map[string]Flag
}

func NewFlags() *Flags {
	return &Flags{
		Flags: map[string]Flag{},
	}
}

func (f *Flags) Contains(flagName string) bool {
	_, ok := f.Flags[flagName]
	return ok
}

func (f *Flags) Get(flagName string) (Flag, bool) {
	flag, ok := f.Flags[flagName]
	return flag, ok
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--")
}

func getFlags(args []string) (flags *Flags, argsWithoutFlags []string) {
	flags = NewFlags()

	for i := 0; i < len(args); i++ {
		if isFlag(args[i]) {
			var val string
			name := args[i]

			if i+1 < len(args) && !isFlag(args[i+1]) {
				val = args[i+1]
				i++
			}

			flag := Flag{
				Name:  name,
				Value: val,
			}

			flags.Flags[name] = flag
		} else {
			argsWithoutFlags = append(argsWithoutFlags, args[i])
		}
	}

	return flags, argsWithoutFlags
}
