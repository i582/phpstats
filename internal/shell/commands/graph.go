package commands

import (
	"fmt"
	"strconv"

	"github.com/i582/phpstats/internal/grapher"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
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
				return
			}
			defer output.Close()

			paths, err := walkers.GlobalCtx.Files.GetFullFileName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			file, _ := walkers.GlobalCtx.Files.Get(paths[0])

			g := grapher.NewGrapher()
			graph := g.FileDeps(file, recursiveLevel, root, block)

			fmt.Fprint(output, graph)

			if show {
				fmt.Println(graph)
			}
		},
	}

	graphClassExecutor := &shell.Executor{
		Name:      "class",
		Help:      "dependency graph for class or interface",
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
		Aliases:   []string{"interface"},
		Func: func(c *shell.Context) {
			recursiveLevelValue := c.GetFlagValue("-r")
			recursiveLevel, _ := strconv.ParseInt(recursiveLevelValue, 0, 64)

			show := c.Flags.Contains("-show")

			output, err := c.ValidateFile("-o")
			if err != nil {
				c.Error(err)
				return
			}
			defer output.Close()

			classes, err := walkers.GlobalCtx.Classes.GetFullClassName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			class, _ := walkers.GlobalCtx.Classes.Get(classes[0])

			g := grapher.NewGrapher()
			graph := g.ClassDeps(class, recursiveLevel)

			fmt.Fprint(output, graph)

			if show {
				fmt.Println(graph)
			}
		},
	}

	graphFuncExecutor := &shell.Executor{
		Name:      "func",
		Help:      "dependency graph for function or method",
		WithValue: true,
		CountArgs: 1,
		Aliases:   []string{"method"},
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Required:  true,
				Help:      "output file",
			},
			&flags.Flag{
				Name: "-show",
				Help: "show graph file in console",
			},
		),
		Func: func(c *shell.Context) {
			show := c.Flags.Contains("-show")

			output, err := c.ValidateFile("-o")
			if err != nil {
				c.Error(err)
				return
			}
			defer output.Close()

			funcs, err := walkers.GlobalCtx.Functions.GetFullFuncName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			fun, _ := walkers.GlobalCtx.Functions.Get(funcs[0])

			g := grapher.NewGrapher()
			graph := g.FuncDeps(fun)

			fmt.Fprint(output, graph)

			if show {
				fmt.Println(graph)
			}
		},
	}

	graphLcom4Executor := &shell.Executor{
		Name:      "lcom4",
		Help:      "show lcom4 connected class components",
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
				Name: "-show",
				Help: "show graph file in console",
			},
		),
		Func: func(c *shell.Context) {
			show := c.Flags.Contains("-show")

			output, err := c.ValidateFile("-o")
			if err != nil {
				c.Error(err)
				return
			}
			defer output.Close()

			classes, err := walkers.GlobalCtx.Classes.GetFullClassName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			class, _ := walkers.GlobalCtx.Classes.Get(classes[0])
			graph := class.Lcom4Graph()

			fmt.Fprint(output, graph)

			if show {
				fmt.Println(graph)
			}
		},
	}

	graphNamespacesExecutor := &shell.Executor{
		Name: "namespaces",
		Help: "show graph with all namespaces",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Required:  true,
				Help:      "output file",
			},
			&flags.Flag{
				Name: "-show",
				Help: "show graph file in console",
			},
		),
		Func: func(c *shell.Context) {
			show := c.Flags.Contains("-show")

			output, err := c.ValidateFile("-o")
			if err != nil {
				c.Error(err)
				return
			}
			defer output.Close()

			g := grapher.NewGrapher()
			graph := g.Namespaces(walkers.GlobalCtx.Namespaces)

			fmt.Fprint(output, graph)

			if show {
				fmt.Println(graph)
			}
		},
	}

	graphNamespaceExecutor := &shell.Executor{
		Name:      "namespace",
		Help:      "show graph with namespace",
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
				Name: "-show",
				Help: "show graph file in console",
			},
		),
		Func: func(c *shell.Context) {
			show := c.Flags.Contains("-show")

			output, err := c.ValidateFile("-o")
			if err != nil {
				c.Error(err)
				return
			}
			defer output.Close()

			ns, ok := walkers.GlobalCtx.Namespaces.GetNamespace(c.Args[0])
			if !ok {
				c.Error(fmt.Errorf("namespace %s not found", c.Args[0]))
				return
			}

			g := grapher.NewGrapher()
			graph := g.Namespace(ns)

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
	graphExecutor.AddExecutor(graphLcom4Executor)
	graphExecutor.AddExecutor(graphNamespacesExecutor)
	graphExecutor.AddExecutor(graphNamespaceExecutor)

	return graphExecutor
}
