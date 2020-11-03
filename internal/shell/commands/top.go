package commands

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/i582/phpstats/internal/representator"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/metrics"
	"github.com/i582/phpstats/internal/stats/walkers"
)

func Top() *shell.Executor {
	topFuncsExecutor := &shell.Executor{
		Name:    "funcs",
		Aliases: []string{"methods"},
		Help:    "show top of functions or methods",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name: "-by-deps",
				Help: "top functions by dependencies",
			},
			&flags.Flag{
				Name: "-by-as-dep",
				Help: "top functions by as dependency",
			},
			&flags.Flag{
				Name: "-by-uses",
				Help: "top functions by uses count",
			},
			&flags.Flag{
				Name: "-by-cc",
				Help: "top functions by cyclomatic complexity",
			},
			&flags.Flag{
				Name: "-r",
				Help: "sort reverse",
			},
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
				Name:      "--output",
				Help:      "output json file",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			reverse := c.Flags.Contains("-r")

			byDeps := c.Flags.Contains("-by-deps")
			byAsDeps := c.Flags.Contains("-by-as-deps")
			byUses := c.Flags.Contains("-by-uses")
			byCC := c.Flags.Contains("-by-cc")

			allFuncs := walkers.GlobalCtx.Functions.GetAll(true, true, true, -1, 0, false, true)

			sort.Slice(allFuncs, func(i, j int) bool {
				switch {
				case byDeps:
					depsI := allFuncs[i].CountDeps()
					depsJ := allFuncs[j].CountDeps()
					if reverse {
						depsI, depsJ = depsJ, depsI
					}
					return depsI > depsJ
				case byAsDeps:
					depsI := allFuncs[i].CountDeps()
					depsJ := allFuncs[j].CountDeps()
					if reverse {
						depsI, depsJ = depsJ, depsI
					}
					return depsI > depsJ
				case byUses:
					usesI := allFuncs[i].CountDeps()
					usesJ := allFuncs[j].CountDeps()
					if reverse {
						usesI, usesJ = usesJ, usesI
					}
					return usesI > usesJ
				case byCC:
					ccI := allFuncs[i].CyclomaticComplexity
					ccJ := allFuncs[j].CyclomaticComplexity
					if reverse {
						ccI, ccJ = ccJ, ccI
					}
					return ccI > ccJ
				}

				nameI := allFuncs[i].Name.Name
				nameJ := allFuncs[j].Name.Name
				if reverse {
					nameI, nameJ = nameJ, nameI
				}
				return nameI < nameJ
			})

			if offset < 0 {
				offset = 0
			}

			if offset > int64(len(allFuncs))-1 {
				offset = int64(len(allFuncs) - 1)
			}

			allFuncs = allFuncs[offset:]

			if count < 0 {
				count = 0
			}

			if count == -1 {
				count = int64(len(allFuncs) - 1)
			}

			if count > int64(len(allFuncs))-1 {
				count = int64(len(allFuncs) - 1)
			}

			allFuncs = allFuncs[:count]

			var f *os.File
			output := c.GetFlagValue("--output")
			if output != "" {
				var err error
				f, err = c.ValidateFile("--output")
				if err != nil {
					log.Fatal(err)
				}
			}

			if f != nil {
				data, err := representator.GetPrettifyJsonFunctionsRepr(allFuncs)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(f, data)
				f.Close()
			} else {
				for _, fn := range allFuncs {
					data := representator.GetStringFunctionRepr(fn)
					fmt.Println(data)
				}
			}
		},
	}

	topClassesExecutor := &shell.Executor{
		Name:    "classes",
		Aliases: []string{"interfaces"},
		Help:    "show top of classes or interfaces",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name: "-by-aff",
				Help: "top classes by afferent coupling",
			},
			&flags.Flag{
				Name: "-by-eff",
				Help: "top classes by efferent coupling",
			},
			&flags.Flag{
				Name: "-by-instab",
				Help: "top classes by instability",
			},
			&flags.Flag{
				Name: "-by-lcom",
				Help: "top classes by Lack of cohesion in methods",
			},
			&flags.Flag{
				Name: "-by-lcom4",
				Help: "top classes by Lack of cohesion in methods 4",
			},
			&flags.Flag{
				Name: "-by-deps",
				Help: "top classes by dependencies",
			},
			&flags.Flag{
				Name: "-by-as-dep",
				Help: "top classes by as dependency",
			},
			&flags.Flag{
				Name: "-r",
				Help: "sort reverse",
			},
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
				Name:      "--output",
				Help:      "output json file",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			reverse := c.Flags.Contains("-r")

			byAff := c.Flags.Contains("-by-aff")
			byEff := c.Flags.Contains("-by-eff")
			byInstab := c.Flags.Contains("-by-instab")
			byLcom := c.Flags.Contains("-by-lcom")
			byLcom4 := c.Flags.Contains("-by-lcom4")
			byDeps := c.Flags.Contains("-by-deps")
			byAsDeps := c.Flags.Contains("-by-as-deps")

			allClasses := walkers.GlobalCtx.Classes.GetAll(false, -1, 0, false)

			sort.Slice(allClasses, func(i, j int) bool {
				affI, effI, instabI := metrics.AfferentEfferentStabilityOfClass(allClasses[i])
				affJ, effJ, instabJ := metrics.AfferentEfferentStabilityOfClass(allClasses[j])

				switch {
				case byAff:
					if reverse {
						affI, affJ = affJ, affI
					}
					return affI > affJ
				case byEff:
					if reverse {
						effI, effJ = effJ, effI
					}
					return effI > effJ
				case byInstab:
					if reverse {
						instabI, instabJ = instabJ, instabI
					}
					return instabI > instabJ
				case byLcom:
					lcomI, _ := metrics.LackOfCohesionInMethodsOfCLass(allClasses[i])
					lcomJ, _ := metrics.LackOfCohesionInMethodsOfCLass(allClasses[j])
					if reverse {
						lcomI, lcomJ = lcomJ, lcomI
					}
					return lcomI > lcomJ
				case byLcom4:
					lcom4I := metrics.Lcom4(allClasses[i])
					lcom4J := metrics.Lcom4(allClasses[j])
					if reverse {
						lcom4I, lcom4J = lcom4J, lcom4I
					}
					return lcom4I > lcom4J
				case byDeps:
					depsI := allClasses[i].Deps.Len()
					depsJ := allClasses[j].Deps.Len()
					if reverse {
						depsI, depsJ = depsJ, depsI
					}
					return depsI > depsJ
				case byAsDeps:
					depsI := allClasses[i].DepsBy.Len()
					depsJ := allClasses[j].DepsBy.Len()
					if reverse {
						depsI, depsJ = depsJ, depsI
					}
					return depsI > depsJ
				}

				nameI := allClasses[i].Name
				nameJ := allClasses[j].Name
				if reverse {
					nameI, nameJ = nameJ, nameI
				}
				return nameI < nameJ
			})

			if offset < 0 {
				offset = 0
			}

			if offset > int64(len(allClasses))-1 {
				offset = int64(len(allClasses) - 1)
			}

			allClasses = allClasses[offset:]

			if count < 0 {
				count = 0
			}

			if count == -1 {
				count = int64(len(allClasses) - 1)
			}

			if count > int64(len(allClasses))-1 {
				count = int64(len(allClasses) - 1)
			}

			allClasses = allClasses[:count]

			var f *os.File
			output := c.GetFlagValue("--output")
			if output != "" {
				var err error
				f, err = c.ValidateFile("--output")
				if err != nil {
					log.Fatal(err)
				}
			}

			if f != nil {
				data, err := representator.GetPrettifyJsonClassesRepr(allClasses)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(f, data)
				f.Close()
			} else {
				for _, class := range allClasses {
					data := representator.GetStringClassRepr(class)
					fmt.Println(data)
				}
			}
		},
	}

	topExecutor := &shell.Executor{
		Name:  "top",
		Help:  "shows top of",
		Flags: flags.NewFlags(),
		Func: func(c *shell.Context) {
			c.ShowHelpPage()
		},
	}

	topExecutor.AddExecutor(topFuncsExecutor)
	topExecutor.AddExecutor(topClassesExecutor)

	return topExecutor
}
