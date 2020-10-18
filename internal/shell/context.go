package shell

import (
	"fmt"

	"github.com/gookit/color"

	"github.com/i582/phpstats/internal/shell/flags"
)

type Context struct {
	Args  []string
	Flags *flags.Flags

	Exec *Executor
}

func (c *Context) Error(err error) {
	color.Red.Printf("Error: %v\n", err)
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

func (c *Context) ShowHelpPage() {
	fmt.Println("Usage:")
	fmt.Println(c.Exec.HelpPage(0))
}
