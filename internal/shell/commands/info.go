package commands

import (
	"fmt"

	"github.com/i582/phpstats/internal/representator"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
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
			classNames, err := walkers.GlobalCtx.Classes.GetFullClassName(c.Args[0])
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

			class, _ := walkers.GlobalCtx.Classes.Get(className)
			data := representator.GetStringClassRepr(class)

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

			funcNameKeys, err := walkers.GlobalCtx.Functions.GetFullFuncName(c.Args[0])
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

			fn, _ := walkers.GlobalCtx.Functions.Get(funcNameKeys[funcKeyIndex])

			data := representator.GetStringFunctionRepr(fn)
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
		),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")

			patches, err := walkers.GlobalCtx.Files.GetFullFileName(c.Args[0])
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

			file, _ := walkers.GlobalCtx.Files.Get(patch)

			var data string
			if full {
				data = representator.GetStringFileRepr(file)
			} else {
				data = representator.GetShortStringFileRepr(file)
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

			ns, ok := walkers.GlobalCtx.Namespaces.GetNamespace(namespace)
			if !ok {
				c.Error(fmt.Errorf("namespace %s not found", c.Args[0]))
				return
			}

			data := representator.GetStringNamespaceRepr(ns)
			fmt.Println(data)
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
