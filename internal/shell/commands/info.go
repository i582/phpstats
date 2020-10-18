package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats"
)

func Info() *shell.Executor {
	classInfoExecutor := &shell.Executor{
		Name:      "class",
		Help:      "info about class or interface",
		WithValue: true,
		Aliases:   []string{"interface"},
		Flags: flags.NewFlags(
			&flags.Flag{
				Name: "-f",
				Help: "output full information",
			},
			&flags.Flag{
				Name: "-metrics",
				Help: "output only metrics",
			},
		),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")
			onlyMetrics := c.Flags.Contains("-metrics")

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

			if onlyMetrics {
				fmt.Println(class.OnlyMetricsString())
				return
			}

			if full {
				fmt.Println(class.ExtraFullString(0))
			} else {
				fmt.Println(class.FullString(0, true))
			}
		},
	}

	funcInfoExecutor := &shell.Executor{
		Name:      "func",
		Help:      "info about function or method",
		WithValue: true,
		Aliases:   []string{"method"},
		Flags: flags.NewFlags(
			&flags.Flag{
				Name: "-f",
				Help: "output full information",
			},
		),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")

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
		Help:      "info about file",
		WithValue: true,
		Flags: flags.NewFlags(
			&flags.Flag{
				Name: "-f",
				Help: "output full information",
			},
			&flags.Flag{
				Name:      "-r",
				Help:      "output recursive",
				Default:   "5",
				WithValue: true,
			},
		),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")
			recursiveFlag, recursive := c.Flags.Get("-r")

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
					c.Error(fmt.Errorf("flag value must be a number"))
				}

				fmt.Println(file.FullStringRecursive(int(count)))
				return
			}

			if full {
				fmt.Println(file.FullString(0))
			} else {
				fmt.Println(file.ShortString(0))
			}
		},
	}

	namespaceInfoExecutor := &shell.Executor{
		Name:      "namespace",
		Help:      "info about namespace",
		WithValue: true,
		Flags:     flags.NewFlags(),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			namespace := c.Args[0]

			classes := stats.NewClasses()
			for _, class := range stats.GlobalCtx.Classes.Classes {
				if strings.Contains(class.Name, namespace) {
					classes.Add(class)
				}
			}

			var aff float64
			var eff float64

			for _, class := range classes.Classes {
				for _, dep := range class.Deps.Classes {
					// если зависимость вне пространства имен
					if !strings.Contains(dep.Name, namespace) {
						aff++
					}
				}

				for _, depBy := range class.DepsBy.Classes {
					// если зависимость вне пространства имен
					if !strings.Contains(depBy.Name, namespace) {
						eff++
					}
				}
			}

			var stability float64
			if eff+aff == 0 {
				stability = 0
			} else {
				stability = eff / (eff + aff)
			}

			var res string

			res += fmt.Sprintf("Пространство имен %s:\n", namespace)

			res += fmt.Sprintf(" Афферентность: %.2f\n", aff)
			res += fmt.Sprintf(" Эфферентность: %.2f\n", eff)
			res += fmt.Sprintf(" Стабильность:  %.2f\n", stability)

			fmt.Println(res)
		},
	}

	infoExecutor := &shell.Executor{
		Name: "info",
		Help: "info about",
		Func: func(c *shell.Context) {
			c.ShowHelpPage()
		},
	}

	infoExecutor.AddExecutor(classInfoExecutor)
	infoExecutor.AddExecutor(funcInfoExecutor)
	infoExecutor.AddExecutor(fileInfoExecutor)
	infoExecutor.AddExecutor(namespaceInfoExecutor)

	return infoExecutor
}
