package commands

import (
	"fmt"
	"strings"

	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/relations"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
)

func Relation() *shell.Executor {
	relationFuncReachabilityExecutor := &shell.Executor{
		Name: "func-reachability",
		Help: "shows the reachability between specific functions",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-c",
				WithValue: true,
				Help:      "count in list",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-o",
				WithValue: true,
				Help:      "offset in list",
				Default:   "0",
			},
			&flags.Flag{
				Name: "--show",
				Help: "show paths in console",
			},
			&flags.Flag{
				Name:      "--parent",
				Help:      "the name of the function from which the reachability will be checked",
				WithValue: true,
			},
			&flags.Flag{
				Name:      "--child",
				Help:      "name of the function for which reachability will be checked",
				WithValue: true,
			},
			&flags.Flag{
				Name:      "--exclude",
				Help:      "comma-separated list of functions to be excluded from found paths",
				WithValue: true,
			},
			&flags.Flag{
				Name:      "--depth",
				WithValue: true,
				Help:      "max search depth",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "--json",
				Help:      "path to the file where the data will be saved in json format",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			count := c.GetIntFlagValue("-c")
			offset := c.GetIntFlagValue("-o")
			showPaths := c.ContainsFlag("--show")
			maxDepth := c.GetIntFlagValue("--depth")
			parentFunctionName := c.GetFlagValue("--parent")
			childFunctionName := c.GetFlagValue("--child")
			excludedFunctionsSlice := strings.Split(c.GetFlagValue("--exclude"), ",")
			excludedFunctions := relations.ReachabilityExcludedMap{}

			if c.GetFlagValue("--exclude") != "" {
				for _, excludedFunctionName := range excludedFunctionsSlice {
					excludedFunction, err := walkers.GlobalCtx.Functions.GetFunctionByPartOfName(excludedFunctionName)
					if err != nil {
						c.Error(err)
						return
					}
					excludedFunctions[excludedFunction] = excludedFunction
				}
			}

			parentFunction, err := walkers.GlobalCtx.Functions.GetFunctionByPartOfName(parentFunctionName)
			if err != nil {
				c.Error(err)
				return
			}

			childFunction, err := walkers.GlobalCtx.Functions.GetFunctionByPartOfName(childFunctionName)
			if err != nil {
				c.Error(err)
				return
			}

			toJson, jsonFile := handleOutputInJson(c)

			rel := relations.GetReachabilityFunction(parentFunction, childFunction, excludedFunctions, maxDepth)

			rel.PrintPaths = showPaths
			rel.PrintCount = count
			rel.PrintOffset = offset

			if toJson {
				data, err := rel.Json()
				if err != nil {
					c.Error(fmt.Errorf("writing paths to file: %v", err))
				}
				fmt.Fprintln(jsonFile, string(data))
				jsonFile.Close()
				cfmt.Printf("The paths was {{successfully}}::green saved to file {{'%s'}}::blue\n", jsonFile.Name())
			} else {
				fmt.Println(rel)
			}
		},
	}

	relationAllExecutor := &shell.Executor{
		Name: "all",
		Help: "shows the relationship between specific classes and functions",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "--classes",
				Help:      "comma-separated list of classes without spaces for which you want to find a relationship with other classes or functions",
				WithValue: true,
			},
			&flags.Flag{
				Name:      "--funcs",
				Help:      "comma-separated list of functions without spaces for which you want to find a relationship with other classes or functions",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			classes := strings.Split(c.GetFlagValue("--classes"), ",")
			if c.GetFlagValue("--classes") == "" {
				classes = nil
			}
			for i := range classes {
				classes[i] = strings.TrimSpace(classes[i])
			}

			funcs := strings.Split(c.GetFlagValue("--funcs"), ",")
			if c.GetFlagValue("--funcs") == "" {
				funcs = nil
			}
			for i := range funcs {
				funcs[i] = strings.TrimSpace(funcs[i])
			}

			for i := 0; i < len(classes); i++ {
				for j := i + 1; j < len(classes); j++ {
					targetClass, err := walkers.GlobalCtx.Classes.GetClassByPartOfName(classes[i])
					if err != nil {
						c.Error(err)
						return
					}

					relatedClass, err := walkers.GlobalCtx.Classes.GetClassByPartOfName(classes[j])
					if err != nil {
						c.Error(err)
						return
					}

					rel := relations.GetClass2ClassRelation(targetClass, relatedClass)
					fmt.Println(rel)
				}
			}

			for i := 0; i < len(funcs); i++ {
				for j := i + 1; j < len(funcs); j++ {
					targetFunction, err := walkers.GlobalCtx.Functions.GetFunctionByPartOfName(funcs[i])
					if err != nil {
						c.Error(err)
						return
					}

					relatedFunction, err := walkers.GlobalCtx.Functions.GetFunctionByPartOfName(funcs[j])
					if err != nil {
						c.Error(err)
						return
					}

					rel := relations.GetFunc2FuncRelation(targetFunction, relatedFunction)
					fmt.Println(rel)
				}
			}

			for i := 0; i < len(classes); i++ {
				for j := 0; j < len(funcs); j++ {
					targetClass, err := walkers.GlobalCtx.Classes.GetClassByPartOfName(classes[i])
					if err != nil {
						c.Error(err)
						return
					}

					relatedFunction, err := walkers.GlobalCtx.Functions.GetFunctionByPartOfName(funcs[j])
					if err != nil {
						c.Error(err)
						return
					}

					rel := relations.GetClass2FuncRelation(targetClass, relatedFunction)
					fmt.Println(rel)
				}
			}
		},
	}

	relationExecutor := &shell.Executor{
		Name:  "relation",
		Help:  "shows relation",
		Flags: flags.NewFlags(),
		Func: func(c *shell.Context) {
			c.ShowHelpPage()
		},
	}

	relationExecutor.AddExecutor(relationFuncReachabilityExecutor)
	relationExecutor.AddExecutor(relationAllExecutor)

	return relationExecutor
}
