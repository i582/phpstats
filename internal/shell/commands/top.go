package commands

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats"
)

func Top() *shell.Executor {
	topFuncsExecutor := &shell.Executor{
		Name: "funcs",
		Help: "show top of functions",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-by-deps",
				WithValue: true,
				Help:      "top functions by dependencies",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-by-as-dep",
				WithValue: true,
				Help:      "top functions by as dependency",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-by-uses",
				WithValue: true,
				Help:      "top functions by uses count",
				Default:   "10",
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

			allFuncs := stats.GlobalCtx.Funcs.GetAll(true, true, true, -1, 0, false, true)

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

			for _, fn := range allFuncs {
				fmt.Println(fn.FullString())
			}
		},
	}

	topClassesExecutor := &shell.Executor{
		Name: "classes",
		Help: "show top of classes",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "-by-aff",
				WithValue: true,
				Help:      "top classes by afferent coupling",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-by-eff",
				WithValue: true,
				Help:      "top classes by efferent coupling",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-by-stab",
				WithValue: true,
				Help:      "top classes by stability",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-by-lcom",
				WithValue: true,
				Help:      "top classes by Lack of cohesion in methods",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-by-lcom4",
				WithValue: true,
				Help:      "top classes by Lack of cohesion in methods 4",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-by-deps",
				WithValue: true,
				Help:      "top classes by dependencies",
				Default:   "10",
			},
			&flags.Flag{
				Name:      "-by-as-dep",
				WithValue: true,
				Help:      "top classes by as dependency",
				Default:   "10",
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
		),
		Func: func(c *shell.Context) {
			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			reverse := c.Flags.Contains("-r")

			byAff := c.Flags.Contains("-by-aff")
			byEff := c.Flags.Contains("-by-eff")
			byStab := c.Flags.Contains("-by-stab")
			byLcom := c.Flags.Contains("-by-lcom")
			byLcom4 := c.Flags.Contains("-by-lcom4")
			byDeps := c.Flags.Contains("-by-deps")
			byAsDeps := c.Flags.Contains("-by-as-deps")

			allClasses := stats.GlobalCtx.Classes.GetAll(false, -1, 0, false)

			sort.Slice(allClasses, func(i, j int) bool {
				affI, effI, stabI := stats.AfferentEfferentStabilityOfClass(allClasses[i])
				affJ, effJ, stabJ := stats.AfferentEfferentStabilityOfClass(allClasses[j])

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
				case byStab:
					if reverse {
						stabI, stabJ = stabJ, stabI
					}
					return stabI > stabJ
				case byLcom:
					lcomI, _ := stats.LackOfCohesionInMethodsOfCLass(allClasses[i])
					lcomJ, _ := stats.LackOfCohesionInMethodsOfCLass(allClasses[j])
					if reverse {
						lcomI, lcomJ = lcomJ, lcomI
					}
					return lcomI > lcomJ
				case byLcom4:
					lcom4I := stats.Lcom4(allClasses[i])
					lcom4J := stats.Lcom4(allClasses[j])
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

			for _, fn := range allClasses {
				fmt.Println(fn.FullString(0, true))
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
