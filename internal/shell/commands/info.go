package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/i582/phpstats/internal/representator"
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
		Flags:     flags.NewFlags(),
		CountArgs: 1,
		Func: func(c *shell.Context) {
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
			data := representator.GetClassRepr(class)

			fmt.Println(data)
		},
	}

	funcInfoExecutor := &shell.Executor{
		Name:      "func",
		Help:      "info about function or method",
		WithValue: true,
		Aliases:   []string{"method"},
		Flags:     flags.NewFlags(),
		CountArgs: 1,
		Func: func(c *shell.Context) {

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

			dataJson, _ := representator.GetJsonFunctionReprWithFlag(fn)
			fmt.Println(dataJson)
			data := representator.GetFunctionRepr(fn)
			fmt.Println(data)
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

			var data string
			if full {
				data = representator.GetFileRepr(file)
			} else {
				data = representator.GetShortFileRepr(file)
			}
			fmt.Println(data)
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
					if !strings.Contains(dep.Name, namespace) {
						aff++
					}
				}

				for _, depBy := range class.DepsBy.Classes {
					if !strings.Contains(depBy.Name, namespace) {
						eff++
					}
				}
			}

			var instability float64
			if eff+aff == 0 {
				instability = 0
			} else {
				instability = eff / (eff + aff)
			}

			var res string

			res += fmt.Sprintf("Namespace %s:\n", namespace)

			res += fmt.Sprintf(" Afferent:  %.2f\n", aff)
			res += fmt.Sprintf(" Efferent:  %.2f\n", eff)
			res += fmt.Sprintf(" Instability: %.2f\n", instability)

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
