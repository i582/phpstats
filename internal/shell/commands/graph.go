package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/grapher"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
	"github.com/i582/phpstats/internal/utils"
)

// main grapher
var g = grapher.NewGrapher()

func Graph() *shell.Executor {
	graphFileExecutor := &shell.Executor{
		Name:      "file",
		Help:      "output dependency graph for file",
		WithValue: true,
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "output file",
			},
			&flags.Flag{
				Name:      "-r",
				WithValue: true,
				Help:      "recursive level",
				Default:   "5",
			},
			&flags.Flag{
				Name: "--root",
				Help: "only root require",
			},
			&flags.Flag{
				Name: "--block",
				Help: "only block require",
			},
			&flags.Flag{
				Name: "--web",
				Help: "show graph in browser",
			},
		),
		CountArgs: 1,
		Func: func(c *shell.Context) {
			recursiveLevel := c.GetIntFlagValue("-r")

			root := c.Flags.Contains("--root")
			block := c.Flags.Contains("--block")
			inBrowser := c.Flags.Contains("--web")

			if !validateOutputPath(c, inBrowser) {
				return
			}

			file, err := walkers.GlobalCtx.Files.GetFileByPartOfName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			graphData := g.FileDeps(file, recursiveLevel, root, block)
			handleGraphOutputWithWeb(c, inBrowser, graphData)
		},
	}

	graphClassExecutor := &shell.Executor{
		Name:      "class",
		Help:      "output dependency graph for class or interface",
		WithValue: true,
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "output file",
			},
			&flags.Flag{
				Name:      "-r",
				WithValue: true,
				Help:      "recursive level",
				Default:   "5",
			},
			&flags.Flag{
				Name: "--inheritance",
				Help: "show the inheritance graph",
			},
			// TODO:
			// &flags.Flag{
			// 	Name: "--superglobals",
			// 	Help: "show the superglobals deps",
			// },
			// &flags.Flag{
			// 	Name: "-g",
			// 	Help: "show classes in groups",
			// },
			&flags.Flag{
				Name: "--web",
				Help: "show graph in browser",
			},
		),
		CountArgs: 1,
		Aliases:   []string{"interface"},
		Func: func(c *shell.Context) {
			recursiveLevel := c.GetIntFlagValue("-r")
			onlyInheritance := c.Flags.Contains("--inheritance")
			onlySuperGlobals := c.Flags.Contains("--superglobals")
			withGroups := c.Flags.Contains("-g")
			inBrowser := c.Flags.Contains("--web")

			if !validateOutputPath(c, inBrowser) {
				return
			}

			class, err := walkers.GlobalCtx.Classes.GetClassByPartOfName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			var graphData string

			switch {
			case onlyInheritance:
				graphData = g.ClassImplementsExtendsDeps(class, recursiveLevel)
			case onlySuperGlobals:
				graphData = g.ClassSuperGlobalsDeps(class)
			default:
				graphData = g.ClassDeps(class, recursiveLevel, withGroups)
			}

			handleGraphOutputWithWeb(c, inBrowser, graphData)
		},
	}

	graphFuncExecutor := &shell.Executor{
		Name:      "func",
		Help:      "output dependency graph for function or method",
		WithValue: true,
		CountArgs: 1,
		Aliases:   []string{"method"},
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "output file",
			},
			&flags.Flag{
				Name:      "-r",
				WithValue: true,
				Help:      "recursive level",
				Default:   "5",
			},
			&flags.Flag{
				Name: "--web",
				Help: "show graph in browser",
			},
		),
		Func: func(c *shell.Context) {
			recursiveLevel := c.GetIntFlagValue("-r")
			inBrowser := c.Flags.Contains("--web")

			if !validateOutputPath(c, inBrowser) {
				return
			}

			fun, err := walkers.GlobalCtx.Functions.GetClassByPartOfName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			graphData := g.NewFuncDeps(fun, recursiveLevel)
			handleGraphOutputWithWeb(c, inBrowser, graphData)
		},
	}

	graphLcom4Executor := &shell.Executor{
		Name:      "lcom4",
		Help:      "output lcom4 connected class components",
		WithValue: true,
		CountArgs: 1,
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "output file",
			},
			&flags.Flag{
				Name: "--web",
				Help: "show graph in browser",
			},
		),
		Func: func(c *shell.Context) {
			inBrowser := c.Flags.Contains("--web")

			if !validateOutputPath(c, inBrowser) {
				return
			}

			class, err := walkers.GlobalCtx.Classes.GetClassByPartOfName(c.Args[0])
			if err != nil {
				c.Error(err)
				return
			}

			graphData := g.Lcom4(class)
			handleGraphOutputWithWeb(c, inBrowser, graphData)
		},
	}

	graphExecutor := &shell.Executor{
		Name: "graph",
		Help: "dependencies graph view in svg",
		Func: func(c *shell.Context) {
			c.ShowHelpPage()
		},
	}

	graphExecutor.AddExecutor(graphFileExecutor)
	graphExecutor.AddExecutor(graphClassExecutor)
	graphExecutor.AddExecutor(graphFuncExecutor)
	graphExecutor.AddExecutor(graphLcom4Executor)

	return graphExecutor
}

func validateOutputPath(c *shell.Context, inBrowser bool) bool {
	if !inBrowser {
		output, err := c.ValidateFile("-o")
		if err != nil {
			c.Error(err)
			return false
		}
		output.Close()
	}
	return true
}

func handleGraphOutputWithWeb(c *shell.Context, inBrowser bool, graphData string) {
	if inBrowser {
		output, ok := c.Flags.Get("-o")
		if !ok {
			output = &flags.Flag{
				Name: "-o",
			}
			c.Flags.Flags["-o"] = output
		}
		output.Value = filepath.Join(utils.DefaultGraphsDir(), "graph.svg")
	}

	handleGraphOutput(c, graphData)
	transformSvgGraph(c)

	if inBrowser {
		err := utils.OpenFile("file:///" + c.GetFlagValue("-o"))
		if err != nil {
			log.Print("error open graph file:", err)
		}
	}
}

func transformSvgGraph(c *shell.Context) {
	name := c.GetFlagValue("-o")
	data, err := ioutil.ReadFile(name)
	if err == nil {
		needStr := `xmlns:xlink="http://www.w3.org/1999/xlink">`
		startGraphData := bytes.Index(data, []byte(needStr))
		startGraphData += len(needStr)
		startSvg := bytes.Index(data, []byte("<svg ")) + 5
		startViewBox := bytes.Index(data, []byte(" viewBox"))
		startEndSvg := bytes.Index(data, []byte("</svg>"))
		var newData []byte
		newData = append(newData, data[0:startSvg]...)
		newData = append(newData, []byte("width=\"100%\" height=\"100%\"")...)
		newData = append(newData, data[startViewBox:startGraphData]...)
		newData = append(newData, []byte(graph.WebAdditionHeader)...)
		newData = append(newData, data[startGraphData:startEndSvg]...)
		newData = append(newData, []byte(graph.WebAdditionFooter)...)
		newData = append(newData, data[startEndSvg:]...)
		err := ioutil.WriteFile(name, newData, 0677)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func handleGraphOutput(c *shell.Context, graph string) {
	name := c.GetFlagValue("-o")
	graphFileName := name + ".gv"
	graphFile, err := c.ValidateFilePath(graphFileName)
	if err != nil {
		c.Error(err)
		return
	}

	fmt.Fprint(graphFile, graph)
	graphFile.Close()

	dot := &grapher.Dot{
		Format:     grapher.Svg,
		InputFile:  graphFileName,
		OutputName: name,
	}
	err = dot.Execute()
	if err != nil {
		c.Error(err)
		return
	}
	return
}
