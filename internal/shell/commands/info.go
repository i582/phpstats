package commands

import (
	"fmt"
	"strconv"

	"phpstats/internal/shell"
	"phpstats/internal/stats"
)

func Info() *shell.Executor {
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

	return infoExecutor
}
