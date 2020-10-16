package main

import (
	"fmt"
	"strconv"

	"phpstats/internal/cli"
	"phpstats/internal/shell"
	"phpstats/internal/stats"
)

func main() {
	cli.RunPhplinterTool(&cli.PhplinterTool{
		Name:    "stats",
		Collect: stats.CollectMain,
		Process: nil,
	})

	s := shell.NewShell()

	classInfoExecutor := &shell.Executor{
		Name:      "class",
		Help:      "info about some class",
		WithValue: true,
		Flags: shell.NewFlags(
			&shell.Flag{
				Name: "-f",
				Help: "output full information",
			},
		),
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")

			if len(c.Args) != 1 {
				c.Error(fmt.Errorf("команда принимает ровно один аргумент"))
				return
			}

			classNames, err := stats.GlobalCtx.Classes.GetFullClassName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			var className string

			if len(classNames) > 1 {
				// choice := c.MultiChoice(classNames, "Какой класс вы имели ввиду?")
				// className = classNames[choice]
				className = classNames[0]
			} else {
				className = classNames[0]
			}

			class, _ := stats.GlobalCtx.Classes.Get(className)

			if full {
				fmt.Println(class.FullString(0))
			} else {
				fmt.Println(class.ShortString(0))
			}
		},
	}

	funcInfoExecutor := &shell.Executor{
		Name:      "func",
		Help:      "info about some func",
		WithValue: true,
		Flags: shell.NewFlags(
			&shell.Flag{
				Name: "-f",
				Help: "output full information",
			},
		),
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")

			if len(c.Args) != 1 {
				c.Error(fmt.Errorf("команда принимает ровно один аргумент"))
				return
			}

			funcNameKeys, err := stats.GlobalCtx.Funcs.GetFullFuncName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			var funcKeyIndex int

			if len(funcNameKeys) > 1 {
				funcManes := make([]string, 0, len(funcNameKeys))
				for _, key := range funcNameKeys {
					funcManes = append(funcManes, key.String())
				}

				// funcKeyIndex = c.MultiChoice(funcManes, "Какую функцию вы имели ввиду?")
				funcKeyIndex = 0
			} else {
				funcKeyIndex = 0
			}

			fn, _ := stats.GlobalCtx.Funcs.Get(funcNameKeys[funcKeyIndex])

			if full {
				fmt.Println(fn.FullString())
			} else {
				fmt.Println(fn.ShortString())
			}
		},
	}

	fileInfoExecutor := &shell.Executor{
		Name:      "file",
		Help:      "info about some file",
		WithValue: true,
		Flags: shell.NewFlags(
			&shell.Flag{
				Name: "-f",
				Help: "output full information",
			},
			&shell.Flag{
				Name:      "-r",
				Help:      "output recursive",
				Default:   "5",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")
			recursiveFlag, recursive := c.Flags.Get("-r")

			if len(c.Args) != 1 {
				c.Error(fmt.Errorf("команда принимает ровно один аргумент"))
				return
			}

			patches, err := stats.GlobalCtx.Files.GetFullFileName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			var patch string

			if len(patches) > 1 {
				// choice := c.MultiChoice(patches, "Какой файл вы имели ввиду?")
				// patch = patches[choice]
				patch = patches[0]
			} else {
				patch = patches[0]
			}

			file, _ := stats.GlobalCtx.Files.Get(patch)

			if recursive {
				count, err := strconv.ParseInt(recursiveFlag.Value, 0, 64)
				if err != nil {
					c.Error(fmt.Errorf("значение флага должно быть числом"))
				}

				fmt.Println(file.FullStringRecursive(int(count)))
			}

			if full {
				fmt.Println(file.FullString(0))
			} else {
				fmt.Println(file.ShortString(0))
			}
		},
	}

	infoExecutor := &shell.Executor{
		Name: "info",
		Help: "info about something",
		Func: func(c *shell.Context) {

		},
	}

	infoExecutor.AddExecutor(classInfoExecutor)
	infoExecutor.AddExecutor(funcInfoExecutor)
	infoExecutor.AddExecutor(fileInfoExecutor)

	s.AddExecutor(infoExecutor)

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
					fmt.Println(class.FullString(0))
				}
			}
		},
	}

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
					fmt.Println(class.FullString(0))
				}
			}
		},
	}

	listExecutor := &shell.Executor{
		Name: "list",
		Help: "list of something",
		Func: func(c *shell.Context) {

		},
	}

	listExecutor.AddExecutor(listFuncExecutor)
	listExecutor.AddExecutor(listMethodExecutor)
	listExecutor.AddExecutor(listFilesExecutor)
	listExecutor.AddExecutor(listClassesExecutor)
	listExecutor.AddExecutor(listInterfaceExecutor)

	s.AddExecutor(listExecutor)

	graphFuncExecutor := &shell.Executor{
		Name:      "func",
		Help:      "graph some func",
		WithValue: true,
		Flags: shell.NewFlags(
			&shell.Flag{
				Name:      "-o",
				WithValue: true,
				Required:  true,
				Help:      "output file",
			},
			&shell.Flag{
				Name:      "-r",
				WithValue: true,
				Help:      "recursive level",
				Default:   "5",
			},
		),
		Func: func(c *shell.Context) {
			recursiveFlag := c.GetFlagValue("-r")

			fmt.Println(recursiveFlag)
		},
	}

	graphExecutor := &shell.Executor{
		Name: "graph",
		Help: "graph view",
		Func: func(c *shell.Context) {

		},
	}

	graphExecutor.AddExecutor(graphFuncExecutor)

	s.AddExecutor(graphExecutor)

	s.Run()
}
