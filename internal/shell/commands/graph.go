package commands

import (
	"fmt"
	"strconv"

	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats"
)

func Graph() *shell.Executor {
	graphFileExecutor := &shell.Executor{
		Name:      "file",
		Help:      "dependency graph for file",
		WithValue: true,
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Required:  true,
				Help:      "output file",
			},
			&flags.Flag{
				Name:      "-r",
				WithValue: true,
				Help:      "recursive level",
				Default:   "5",
			},
			&flags.Flag{
				Name: "-root",
				Help: "only root require",
			},
			&flags.Flag{
				Name: "-block",
				Help: "only block require",
			},
			&flags.Flag{
				Name: "-show",
				Help: "show graph file in console",
			},
		),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			recursiveLevelValue := c.GetFlagValue("-r")
			recursiveLevel, _ := strconv.ParseInt(recursiveLevelValue, 0, 64)

			root := c.Flags.Contains("-root")
			block := c.Flags.Contains("-block")
			show := c.Flags.Contains("-show")

			output, err := c.ValidateFile("-o")
			if err != nil {
				c.Error(err)
			}
			defer output.Close()

			paths, err := stats.GlobalCtx.Files.GetFullFileName(c.Args[0])
			if err != nil {
				fmt.Printf("Файл %s не найден!\n", c.Args[0])
				return
			}

			file, _ := stats.GlobalCtx.Files.Get(paths[0])
			graph := file.GraphvizRecursive(recursiveLevel, root, block)

			fmt.Fprint(output, graph)

			if show {
				fmt.Println(graph)
			}
		},
	}

	graphClassExecutor := &shell.Executor{
		Name:      "class",
		Help:      "dependency graph for class",
		WithValue: true,
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Required:  true,
				Help:      "output file",
			},
			&flags.Flag{
				Name:      "-r",
				WithValue: true,
				Help:      "recursive level",
				Default:   "5",
			},
			&flags.Flag{
				Name: "-show",
				Help: "show graph file in console",
			},
		),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			recursiveLevelValue := c.GetFlagValue("-r")
			recursiveLevel, _ := strconv.ParseInt(recursiveLevelValue, 0, 64)

			show := c.Flags.Contains("-show")

			output, err := c.ValidateFile("-o")
			if err != nil {
				c.Error(err)
			}
			defer output.Close()

			classes, err := stats.GlobalCtx.Classes.GetFullClassName(c.Args[0])
			if err != nil {
				fmt.Printf("Класс %s не найден!\n", c.Args[0])
				return
			}

			class, _ := stats.GlobalCtx.Classes.Get(classes[0])
			graph := class.GraphvizRecursive(0, recursiveLevel, map[string]struct{}{})

			fmt.Fprint(output, graph)

			if show {
				fmt.Println(graph)
			}
		},
	}

	graphFuncExecutor := &shell.Executor{
		Name:      "func",
		Help:      "dependency graph for function",
		WithValue: true,
		CountArgs: 1,
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Required:  true,
				Help:      "output file",
			},
			&flags.Flag{
				Name:      "-r",
				WithValue: true,
				Help:      "recursive level",
				Default:   "5",
			},
			&flags.Flag{
				Name: "-show",
				Help: "show graph file in console",
			},
		),
		Func: func(c *shell.Context) {
			recursiveLevelValue := c.GetFlagValue("-r")
			recursiveLevel, _ := strconv.ParseInt(recursiveLevelValue, 0, 64)

			show := c.Flags.Contains("-show")

			outputPath := c.GetFlagValue("-o")
			if outputPath == "" {
				c.Error(fmt.Errorf("invalid filepath"))
				return
			}

			output, err := c.ValidateFile("-o")
			if err != nil {
				c.Error(err)
			}
			defer output.Close()

			funcs, err := stats.GlobalCtx.Funcs.GetFullFuncName(c.Args[0])
			if err != nil {
				fmt.Printf("Функция %s не найден!\n", c.Args[0])
				return
			}

			fun, _ := stats.GlobalCtx.Funcs.Get(funcs[0])
			graph := fun.GraphvizRecursive(0, recursiveLevel, map[string]struct{}{})

			fmt.Fprint(output, graph)

			if show {
				fmt.Println(graph)
			}
		},
	}

	graphExecutor := &shell.Executor{
		Name: "graph",
		Help: "dependencies graph view",
		Func: func(c *shell.Context) {
			c.ShowHelpPage()
		},
	}

	graphExecutor.AddExecutor(graphFileExecutor)
	graphExecutor.AddExecutor(graphClassExecutor)
	graphExecutor.AddExecutor(graphFuncExecutor)

	return graphExecutor
}
