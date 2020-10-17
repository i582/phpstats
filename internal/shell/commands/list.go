package commands

import (
	"fmt"
	"strconv"

	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats"
)

func List() *shell.Executor {
	listFuncExecutor := &shell.Executor{
		Name: "funcs",
		Help: "show list of functions",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "offset in list",
				Default:   "0",
			},
			&flags.Flag{
				Name: "-e",
				Help: "show embedded functions",
			},
		),
		Func: func(c *shell.Context) {
			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			withEmbeddedFuncs := c.Flags.Contains("-e")

			funcs := stats.GlobalCtx.Funcs.GetAll(false, true, false, count, offset, true, withEmbeddedFuncs)

			for _, fn := range funcs {
				fmt.Println(fn.FullString())
			}
		},
	}

	listMethodExecutor := &shell.Executor{
		Name: "methods",
		Help: "show list of methods",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "offset in list",
				Default:   "0",
			},
		),
		Func: func(c *shell.Context) {
			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			funcs := stats.GlobalCtx.Funcs.GetAll(true, false, false, count, offset, true, false)

			for _, fn := range funcs {
				fmt.Println(fn.FullString())
			}
		},
	}

	listFilesExecutor := &shell.Executor{
		Name: "files",
		Help: "show list of files",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "offset in list",
				Default:   "0",
			},
			&flags.Flag{
				Name: "-f",
				Help: "show full information",
			},
		),
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")

			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			files := stats.GlobalCtx.Files.GetAll(count, offset, true)

			for _, file := range files {
				if full {
					fmt.Print(file.FullString(0))
				} else {
					fmt.Print(file.ExtraShortString(0))
				}
			}
		},
	}

	listClassesExecutor := &shell.Executor{
		Name: "classes",
		Help: "show list of classes",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "offset in list",
				Default:   "0",
			},
			&flags.Flag{
				Name: "-f",
				Help: "show full information",
			},
		),
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")

			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			classes := stats.GlobalCtx.Classes.GetAllClasses(count, offset, true)

			for _, class := range classes {
				if full {
					fmt.Println(class.ExtraFullString(0))
				} else {
					fmt.Println(class.FullString(0, true))
				}
			}
		},
	}

	listInterfaceExecutor := &shell.Executor{
		Name: "interfaces",
		Help: "show list of interfaces",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "offset in list",
				Default:   "0",
			},
			&flags.Flag{
				Name: "-f",
				Help: "show full information",
			},
		),
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")

			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			classes := stats.GlobalCtx.Classes.GetAllInterfaces(count, offset, true)

			for _, class := range classes {
				if full {
					fmt.Println(class.ExtraFullString(0))
				} else {
					fmt.Println(class.FullString(0, true))
				}
			}
		},
	}

	listExecutor := &shell.Executor{
		Name: "list",
		Help: "list of",
		Func: func(c *shell.Context) {
			c.ShowHelpPage()
		},
	}

	listExecutor.AddExecutor(listFuncExecutor)
	listExecutor.AddExecutor(listMethodExecutor)
	listExecutor.AddExecutor(listFilesExecutor)
	listExecutor.AddExecutor(listClassesExecutor)
	listExecutor.AddExecutor(listInterfaceExecutor)

	return listExecutor
}
