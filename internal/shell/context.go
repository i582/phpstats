package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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

func (c *Context) GetIntFlagValue(flag string) int64 {
	val := c.GetFlagValue(flag)
	intVal, _ := strconv.ParseInt(val, 0, 64)
	return intVal
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

func (c *Context) ValidateFilePath(path string) (*os.File, error) {
	if path == "" {
		return nil, fmt.Errorf("empty filepath")
	}

	err := os.MkdirAll(filepath.Dir(path), 0677)
	if err != nil {
		return nil, fmt.Errorf("dirs not created %v", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0677)
	if err != nil {
		return nil, fmt.Errorf("file not open %v", err)
	}

	return file, nil
}

func (c *Context) ValidateFile(flag string) (*os.File, error) {
	path := c.GetFlagValue(flag)
	return c.ValidateFilePath(path)
}
