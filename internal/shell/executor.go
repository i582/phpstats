package shell

import (
	"fmt"

	"phpstats/internal/stats"
)

type Context struct {
	Args  []string
	Flags *Flags

	Exec *Executor
}

func (c *Context) Error(err error) {
	fmt.Println("Error:", err)
}

func (c *Context) ContainsFlag(flag string) bool {
	ok := c.Flags.Contains(flag)
	if ok {
		return true
	}

	f, ok := c.Exec.Flags.Get(flag)
	if ok {
		if f.Default == "" {
			return false
		}

		return true
	}

	return false
}

func (c *Context) GetFlagValue(flag string) string {
	f, ok := c.Flags.Get(flag)
	if ok {
		return f.Value
	}

	f, ok = c.Exec.Flags.Get(flag)
	if ok {
		if f.Default == "" {
			return ""
		}

		return f.Default
	}

	return ""
}

type Executors map[string]*Executor

type Executor struct {
	Name string
	Help string

	WithValue bool

	Flags *Flags

	SubExecs Executors

	Func func(*Context)
}

func (e *Executor) HelpPage(level int) string {
	var res string

	var withValueSpan string
	if e.WithValue {
		withValueSpan = "<value>"
	}

	res += fmt.Sprintf("%s  %s %-*s%s\n", stats.GenIndent(level), e.Name, 20-len(stats.GenIndent(level))-len(e.Name), withValueSpan, e.Help)

	if e.Flags != nil {
		for _, flag := range e.Flags.Flags {
			res += fmt.Sprintf("%s    %s\n", stats.GenIndent(level), flag)
		}
	}

	for _, e := range e.SubExecs {
		res += fmt.Sprintln(e.HelpPage(level + 1))
	}

	return res
}

func (e *Executor) Execute(ctx *Context) {
	if e.Flags == nil {
		e.Flags = NewFlags()
	}

	if len(ctx.Args) > 0 {
		arg := ctx.Args[0]

		if exec, ok := e.SubExecs[arg]; ok {
			exec.Execute(&Context{
				Args:  ctx.Args[1:],
				Flags: NewFlags(),
			})
			return
		}

		ctx.Flags, ctx.Args = getFlags(ctx.Args, e.Flags)
	}

	ctx.Exec = e

	e.Func(ctx)
}

func (e *Executor) AddExecutor(exec *Executor) {
	if e.SubExecs == nil {
		e.SubExecs = Executors{}
	}

	e.SubExecs[exec.Name] = exec
}
