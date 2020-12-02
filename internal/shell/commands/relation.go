package commands

import (
	"fmt"
	"strings"

	"github.com/i582/phpstats/internal/relations"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
)

func Relation() *shell.Executor {
	relationAllExecutor := &shell.Executor{
		Name: "all",
		Help: "shows the relationship between specific classes and functions",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "--related",
				Help:      "related class",
				WithValue: true,
			},
			&flags.Flag{
				Name:      "--classes",
				Help:      "target class",
				WithValue: true,
			},
			&flags.Flag{
				Name:      "--funcs",
				Help:      "target class",
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

	relationExecutor.AddExecutor(relationAllExecutor)

	return relationExecutor
}
