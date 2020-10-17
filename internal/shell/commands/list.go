package commands

import (
	"fmt"
	"strconv"

	"phpstats/internal/shell"
	"phpstats/internal/stats"
)

func List() *shell.Executor {
	listFuncExecutor := &shell.Executor{
		Name: "funcs",
		Help: "show list funcs",
		Flags: shell.NewFlags(
			&shell.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&shell.Flag{
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

			funcs := stats.GlobalCtx.Funcs.GetAll(false, true, false, count, offset, true)

			for _, fn := range funcs {
				fmt.Print(fn.FullString())
			}
		},
	}

	listMethodExecutor := &shell.Executor{
		Name: "methods",
		Help: "show list methods",
		Flags: shell.NewFlags(
			&shell.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&shell.Flag{
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

			funcs := stats.GlobalCtx.Funcs.GetAll(true, false, false, count, offset, true)

			for _, fn := range funcs {
				fmt.Print(fn.FullString())
			}
		},
	}

	listFilesExecutor := &shell.Executor{
		Name: "files",
		Help: "show list of files",
		Flags: shell.NewFlags(
			&shell.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&shell.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "offset in list",
				Default:   "0",
			},
			&shell.Flag{
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
		Flags: shell.NewFlags(
			&shell.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&shell.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "offset in list",
				Default:   "0",
			},
			&shell.Flag{
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
	} // \VK\API\Library\DeprecatedWrappers

	listInterfaceExecutor := &shell.Executor{
		Name: "ifaces",
		Help: "show list of interfaces",
		Flags: shell.NewFlags(
			&shell.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&shell.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "offset in list",
				Default:   "0",
			},
			&shell.Flag{
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
		Help: "list of something",
		Func: func(c *shell.Context) {
			fmt.Println("Usage:")
			fmt.Println(c.Exec.HelpPage(0))
		},
	}

	listExecutor.AddExecutor(listFuncExecutor)
	listExecutor.AddExecutor(listMethodExecutor)
	listExecutor.AddExecutor(listFilesExecutor)
	listExecutor.AddExecutor(listClassesExecutor)
	listExecutor.AddExecutor(listInterfaceExecutor)

	return listExecutor
}
