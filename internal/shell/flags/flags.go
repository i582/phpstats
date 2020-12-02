package flags

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/i582/phpstats/internal/utils"
)

type Flag struct {
	Name  string
	Value string

	Help    string
	Default string

	Required  bool
	WithValue bool
}

func (f *Flag) String() string {
	var res string

	var withValueSpan string
	if f.WithValue {
		withValueSpan = " <value>"
	}

	var defaultSpan string
	if f.Default != "" {
		defaultSpan = fmt.Sprintf(" (default: %s)", f.Default)
	}

	if f.Required {
		res += fmt.Sprintf(" %s%-*s  %s%s", f.Name, 29-len(f.Name)-1, withValueSpan, f.Help, defaultSpan)
	} else {
		res += fmt.Sprintf("[%s%s]%-*s %s%s", f.Name, withValueSpan, 29-len(f.Name)-len(withValueSpan)-1, "", f.Help, defaultSpan)
	}

	return res
}

type Flags struct {
	Flags map[string]*Flag
}

func NewFlags(flags ...*Flag) *Flags {
	flagsMap := make(map[string]*Flag, len(flags))

	for _, flag := range flags {
		flagsMap[flag.Name] = flag
	}

	return &Flags{
		Flags: flagsMap,
	}
}

func (f *Flags) Contains(flagName string) bool {
	_, ok := f.Flags[flagName]
	return ok
}

func (f *Flags) Get(flagName string) (*Flag, bool) {
	flag, ok := f.Flags[flagName]
	return flag, ok
}

func flagParsed(arg string) bool {
	if strings.Contains(arg, "=") {
		return false
	}
	return true
}

func ParseFlag(flag string) []string {
	if flagParsed(flag) {
		return []string{flag}
	}

	return strings.Split(flag, "=")
}

func isFlag(arg string) bool {
	if !strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
		return false
	}

	argV := strings.TrimLeft(arg, "-")

	_, err := strconv.ParseInt(argV, 0, 64)
	if err == nil {
		return false
	}

	return true
}

func ParseFlags(args []string, allowed *Flags) (flags *Flags, argsWithoutFlags []string) {
	flags = &Flags{
		Flags: map[string]*Flag{},
	}

	parsedArgs := make([]string, 0, len(args))

	for _, arg := range args {
		if isFlag(arg) {
			parsedArgs = append(parsedArgs, ParseFlag(arg)...)
			continue
		}

		parsedArgs = append(parsedArgs, arg)
	}

	for i := 0; i < len(parsedArgs); i++ {
		if isFlag(parsedArgs[i]) {
			var val string
			name := parsedArgs[i]

			var needValue bool
			if f, ok := allowed.Get(name); ok {
				needValue = f.WithValue
			}

			if needValue {
				if i+1 < len(parsedArgs) && !isFlag(parsedArgs[i+1]) {
					val = parsedArgs[i+1]
					i++
				}
			}

			flag := &Flag{
				Name:  name,
				Value: utils.Unquote(val),
			}

			flags.Flags[name] = flag
		} else {
			argsWithoutFlags = append(argsWithoutFlags, parsedArgs[i])
		}
	}

	return flags, argsWithoutFlags
}
