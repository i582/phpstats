package shell

import (
	"fmt"
	"sort"
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
		execFlags := make([]*flags.Flag, 0, len(e.SubExecs))
		for _, f := range e.Flags.Flags {
			execFlags = append(execFlags, f)
		}
		sort.Slice(execFlags, func(i, j int) bool {
			if execFlags[i].WithValue && !execFlags[j].WithValue {
				return true
			}
			if !execFlags[i].WithValue && execFlags[j].WithValue {
				return false
			}

			if strings.Contains(execFlags[i].Name, "--") && !strings.Contains(execFlags[j].Name, "--") {
				return false
			}
			if !strings.Contains(execFlags[i].Name, "--") && strings.Contains(execFlags[j].Name, "--") {
				return true
			}
			if strings.Contains(execFlags[i].Name, "--") && strings.Contains(execFlags[j].Name, "--") {
				return execFlags[i].Name < execFlags[j].Name
			}

			return execFlags[i].Name < execFlags[j].Name
		})

		for _, flag := range execFlags {
			res += fmt.Sprintf("%s    %s\n", utils.GenIndent(level), flag)
		}
	}

	execs := make([]*Executor, 0, len(e.SubExecs))
	for _, e := range e.SubExecs {
		execs = append(execs, e)
	}
	sort.Slice(execs, func(i, j int) bool {
		return execs[i].Name < execs[j].Name
	})

	for _, e := range execs {
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

	e.Func(ctx)
}

func (e *Executor) AddExecutor(exec *Executor) {
	if e.SubExecs == nil {
		e.SubExecs = Executors{}
	}

	e.SubExecs[exec.Name] = exec
}
