package shell

import (
	"fmt"
	"strings"

	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/utils"
)

type Executors map[string]*Executor

type Executor struct {
	Name string
	Help string

	Aliases []string

	WithValue bool

	Flags     *flags.Flags
	CountArgs int

	SubExecs Executors

	Func func(*Context)
}

func (e *Executor) HelpPage(level int) string {
	var res string

	var withValueSpan string
	if e.WithValue {
		withValueSpan = "<value>"
	}

	aliases := strings.Join(e.Aliases, ",")
	if aliases != "" {
		aliases = "(or " + aliases + ")"
	}

	res += fmt.Sprintf("%s  %s %s %-*s%s\n", utils.GenIndent(level), e.Name, aliases, 35-len(utils.GenIndent(level))-len(e.Name)-len(aliases)-1, withValueSpan, e.Help)

	if e.Flags != nil {
		for _, flag := range e.Flags.Flags {
			res += fmt.Sprintf("%s    %s\n", utils.GenIndent(level), flag)
		}
	}

	for _, e := range e.SubExecs {
		res += fmt.Sprintln(e.HelpPage(level + 1))
	}

	return res
}

func (e *Executor) findSubExec(name string) (*Executor, bool) {
	for _, exec := range e.SubExecs {
		if exec.Name == name {
			return exec, true
		} else {
			for _, alias := range exec.Aliases {
				if alias == name {
					return exec, true
				}
			}
		}
	}

	return nil, false
}

func (e *Executor) Execute(ctx *Context) {
	if e.Flags == nil {
		e.Flags = flags.NewFlags()
	}

	if len(ctx.Args) > 0 {
		command := ctx.Args[0]

		if exec, ok := e.findSubExec(command); ok {
			exec.Execute(&Context{
				Args:  ctx.Args[1:],
				Flags: flags.NewFlags(),
			})
			return
		}

		ctx.Flags, ctx.Args = flags.ParseFlags(ctx.Args, e.Flags)
	}

	ctx.Exec = e

	if e.CountArgs == 0 && len(e.SubExecs) != 0 {
		ctx.ShowHelpPage()
		return
	}

	if e.CountArgs != -1 && len(ctx.Args) != e.CountArgs {
		ctx.Error(fmt.Errorf("command %s takes exactly %d argument\n", e.Name, e.CountArgs))
		return
	}

	fmt.Println()
	e.Func(ctx)
}

func (e *Executor) AddExecutor(exec *Executor) {
	if e.SubExecs == nil {
		e.SubExecs = Executors{}
	}

	e.SubExecs[exec.Name] = exec
}
