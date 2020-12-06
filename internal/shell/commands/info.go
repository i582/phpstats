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
		Help:      "shows info about a specific class",
		WithValue: true,
		Flags:     flags.NewFlags(),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			fmt.Printf("Show information about %s class\n\n", c.Args[0])

			class, err := walkers.GlobalCtx.Classes.GetClassByPartOfName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			data := representator.GetStringClassRepr(class)
			fmt.Println(data)
		},
	}

	ifaceInfoExecutor := &shell.Executor{
		Name:      "interface",
		Help:      "shows info about a specific interface",
		WithValue: true,
		Aliases:   []string{"iface"},
		Flags:     flags.NewFlags(),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			fmt.Printf("Show information about %s interface\n\n", c.Args[0])

			class, err := walkers.GlobalCtx.Classes.GetInterfaceByPartOfName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			data := representator.GetStringClassRepr(class)
			fmt.Println(data)
		},
	}

	traitInfoExecutor := &shell.Executor{
		Name:      "trait",
		Help:      "shows info about a specific trait",
		WithValue: true,
		Flags:     flags.NewFlags(),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			fmt.Printf("Show information about %s trait\n\n", c.Args[0])

			trait, err := walkers.GlobalCtx.Classes.GetTraitByPartOfName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			data := representator.GetStringClassRepr(trait)
			fmt.Println(data)
		},
	}

	funcInfoExecutor := &shell.Executor{
		Name:      "func",
		Help:      "shows info about a specific function or method",
		WithValue: true,
		Aliases:   []string{"method"},
		Flags:     flags.NewFlags(),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			fmt.Printf("Show information about %s function\n\n", c.Args[0])

			fn, err := walkers.GlobalCtx.Functions.GetFunctionByPartOfName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			data := representator.GetStringFunctionRepr(fn)
			fmt.Println(data)
		},
	}

	fileInfoExecutor := &shell.Executor{
		Name:      "file",
		Help:      "shows info about a specific file",
		WithValue: true,
		Flags:     flags.NewFlags(),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			fmt.Printf("Show information about %s file\n\n", c.Args[0])

			file, err := walkers.GlobalCtx.Files.GetFileByPartOfName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			data := representator.GetStringFileRepr(file)
			fmt.Println(data)
		},
	}

	namespaceInfoExecutor := &shell.Executor{
		Name:      "namespace",
		Help:      "shows info about a specific namespace",
		WithValue: true,
		Flags:     flags.NewFlags(),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			fmt.Printf("Show information about %s namespace\n\n", c.Args[0])

			ns, ok := walkers.GlobalCtx.Namespaces.GetNamespace(c.Args[0])
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
		Help: "shows info",
		Func: func(c *shell.Context) {
			c.ShowHelpPage()
		},
	}

	infoExecutor.AddExecutor(classInfoExecutor)
	infoExecutor.AddExecutor(ifaceInfoExecutor)
	infoExecutor.AddExecutor(traitInfoExecutor)
	infoExecutor.AddExecutor(funcInfoExecutor)
	infoExecutor.AddExecutor(fileInfoExecutor)
	infoExecutor.AddExecutor(namespaceInfoExecutor)

	return infoExecutor
}
