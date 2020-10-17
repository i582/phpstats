package commands

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"phpstats/internal/shell"
	"phpstats/internal/shell/flags"
	"phpstats/internal/stats"
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

			outputPath := c.GetFlagValue("-o")
			if outputPath == "" {
				c.Error(fmt.Errorf("invalid filepath\n"))
				return
			}

			paths, err := stats.GlobalCtx.Files.GetFullFileName(c.Args[0])
			if err != nil {
				fmt.Printf("Файл %s не найден!\n", c.Args[0])
				return
			}

			var res string

			file, _ := stats.GlobalCtx.Files.Get(paths[0])

			outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				log.Fatalf("file not open %v", err)
			}

			res += file.GraphvizRecursive(recursiveLevel, root, block)

			fmt.Fprint(outputFile, res)
			outputFile.Close()

			if show {
				fmt.Println(res)
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

			outputPath := c.GetFlagValue("-o")
			if outputPath == "" {
				c.Error(fmt.Errorf("invalid filepath\n"))
				return
			}

			classes, err := stats.GlobalCtx.Classes.GetFullClassName(c.Args[0])
			if err != nil {
				fmt.Printf("Класс %s не найден!\n", c.Args[0])
				return
			}

			var res string

			class, _ := stats.GlobalCtx.Classes.Get(classes[0])

			outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				log.Fatalf("file not open %v", err)
			}

			res += class.GraphvizRecursive(0, recursiveLevel, map[string]struct{}{})

			fmt.Fprint(outputFile, res)
			outputFile.Close()

			if show {
				fmt.Println(res)
			}
		},
	}

	graphFuncExecutor := &shell.Executor{
		Name:      "func",
		Help:      "dependency graph for function",
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
		Func: func(c *shell.Context) {
			recursiveLevelValue := c.GetFlagValue("-r")
			recursiveLevel, _ := strconv.ParseInt(recursiveLevelValue, 0, 64)

			show := c.Flags.Contains("-show")

			outputPath := c.GetFlagValue("-o")
			if outputPath == "" {
				c.Error(fmt.Errorf("invalid filepath\n"))
				return
			}

			funcs, err := stats.GlobalCtx.Funcs.GetFullFuncName(c.Args[0])
			if err != nil {
				fmt.Printf("Функция %s не найден!\n", c.Args[0])
				return
			}

			var res string

			fun, _ := stats.GlobalCtx.Funcs.Get(funcs[0])

			outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				c.Error(fmt.Errorf("file not open %v", err))
			}

			res += fun.GraphvizRecursive(0, recursiveLevel, map[string]struct{}{})

			fmt.Fprint(outputFile, res)
			outputFile.Close()

			if show {
				fmt.Println(res)
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
