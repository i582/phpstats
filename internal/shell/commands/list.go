package commands

import (
	"fmt"

	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/getter"
	"github.com/i582/phpstats/internal/representator"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
)

func List() *shell.Executor {
	listFuncExecutor := &shell.Executor{
		Name: "funcs",
		Help: "shows list of functions",
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
				Name:      "--sort",
				WithValue: true,
				Help:      "column number by which sorting will be performed",
				Default:   "2",
			},
			&flags.Flag{
				Name: "-r",
				Help: "reverse sort",
			},
			&flags.Flag{
				Name: "-e",
				Help: "show embedded functions",
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
			sortColumn := c.GetIntFlagValue("--sort")
			reverseSort := c.Flags.Contains("-r")

			withEmbeddedFuncs := c.Flags.Contains("-e")
			toJson, jsonFile := handleOutputInJson(c)

			funcs := getter.GetFunctionsByOptions(walkers.GlobalCtx.Functions, getter.FunctionsGetOptions{
				OnlyFuncs:         true,
				Count:             count,
				Offset:            offset,
				WithEmbeddedFuncs: withEmbeddedFuncs,
				SortColumn:        sortColumn,
				ReverseSort:       reverseSort,
			})

			if toJson {
				data, err := representator.GetPrettifyJsonFunctionsRepr(funcs)
				if err != nil {
					c.Error(fmt.Errorf("writing list to file: %v", err))
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
				cfmt.Printf("The functions list was {{successfully}}::green saved to file {{'%s'}}::blue\n", jsonFile.Name())
			} else {
				fmt.Printf("Showing %d functions out of %d starting from %d\n\n", len(funcs), walkers.GlobalCtx.Functions.CountFunctions(withEmbeddedFuncs), offset+1)
				data := representator.GetTableFunctionsRepr(funcs, offset)
				fmt.Println(data)
			}
		},
	}

	listMethodExecutor := &shell.Executor{
		Name: "methods",
		Help: "shows list of methods",
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
				Name:      "--sort",
				WithValue: true,
				Help:      "column number by which sorting will be performed",
				Default:   "2",
			},
			&flags.Flag{
				Name: "-r",
				Help: "reverse sort",
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
			sortColumn := c.GetIntFlagValue("--sort")
			reverseSort := c.Flags.Contains("-r")

			toJson, jsonFile := handleOutputInJson(c)

			methods := getter.GetFunctionsByOptions(walkers.GlobalCtx.Functions, getter.FunctionsGetOptions{
				OnlyMethods: true,
				Count:       count,
				Offset:      offset,
				SortColumn:  sortColumn,
				ReverseSort: reverseSort,
			})

			if toJson {
				data, err := representator.GetPrettifyJsonFunctionsRepr(methods)
				if err != nil {
					c.Error(fmt.Errorf("writing list to file: %v", err))
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
				cfmt.Printf("The methods list was {{successfully}}::green saved to file {{'%s'}}::blue\n", jsonFile.Name())
			} else {
				fmt.Printf("Showing %d methods out of %d starting from %d\n\n", len(methods), walkers.GlobalCtx.Functions.CountMethods(), offset+1)
				data := representator.GetTableFunctionsRepr(methods, offset)
				fmt.Println(data)
			}
		},
	}

	listFilesExecutor := &shell.Executor{
		Name: "files",
		Help: "shows list of files",
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
				Name:      "--sort",
				WithValue: true,
				Help:      "column number by which sorting will be performed",
				Default:   "2",
			},
			&flags.Flag{
				Name: "-r",
				Help: "reverse sort",
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
			sortColumn := c.GetIntFlagValue("--sort")
			reverseSort := c.Flags.Contains("-r")

			toJson, jsonFile := handleOutputInJson(c)

			files := getter.GetFilesByOptions(walkers.GlobalCtx.Files, getter.FilesGetOptions{
				Count:       count,
				Offset:      offset,
				SortColumn:  sortColumn,
				ReverseSort: reverseSort,
			})

			if toJson {
				data, err := representator.GetPrettifyJsonFilesRepr(files)
				if err != nil {
					c.Error(fmt.Errorf("writing list to file: %v", err))
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
				cfmt.Printf("The files list was {{successfully}}::green saved to file {{'%s'}}::blue\n", jsonFile.Name())
			} else {
				fmt.Printf("Showing %d files out of %d starting from %d\n\n", len(files), walkers.GlobalCtx.Files.Len(), offset+1)
				data := representator.GetTableFilesRepr(files, offset)
				fmt.Println(data)
			}
		},
	}

	listClassesExecutor := &shell.Executor{
		Name: "classes",
		Help: "shows list of classes",
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
				Name:      "--sort",
				WithValue: true,
				Help:      "column number by which sorting will be performed",
				Default:   "2",
			},
			&flags.Flag{
				Name: "-r",
				Help: "reverse sort",
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
			sortColumn := c.GetIntFlagValue("--sort")
			reverseSort := c.Flags.Contains("-r")

			toJson, jsonFile := handleOutputInJson(c)

			classes := getter.GetClassesByOption(walkers.GlobalCtx.Classes, getter.ClassesGetOptions{
				OnlyClasses: true,
				Count:       count,
				Offset:      offset,
				SortColumn:  sortColumn,
				ReverseSort: reverseSort,
			})

			if toJson {
				data, err := representator.GetPrettifyJsonClassesRepr(classes)
				if err != nil {
					c.Error(fmt.Errorf("writing list to file: %v", err))
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
				cfmt.Printf("The classes list was {{successfully}}::green saved to file {{'%s'}}::blue\n", jsonFile.Name())
			} else {
				fmt.Printf("Showing %d classes out of %d starting from %d\n\n", len(classes), walkers.GlobalCtx.Classes.CountClasses(), offset+1)
				data := representator.GetTableClassesRepr(classes, offset)
				fmt.Println(data)
			}
		},
	}

	listInterfaceExecutor := &shell.Executor{
		Name:    "interfaces",
		Aliases: []string{"ifaces"},
		Help:    "shows list of interfaces",
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
				Name:      "--sort",
				WithValue: true,
				Help:      "column number by which sorting will be performed",
				Default:   "2",
			},
			&flags.Flag{
				Name: "-r",
				Help: "reverse sort",
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
			sortColumn := c.GetIntFlagValue("--sort")
			reverseSort := c.Flags.Contains("-r")

			toJson, jsonFile := handleOutputInJson(c)

			ifaces := getter.GetClassesByOption(walkers.GlobalCtx.Classes, getter.ClassesGetOptions{
				OnlyInterfaces: true,
				Count:          count,
				Offset:         offset,
				SortColumn:     sortColumn,
				ReverseSort:    reverseSort,
			})

			if toJson {
				data, err := representator.GetPrettifyJsonClassesRepr(ifaces)
				if err != nil {
					c.Error(fmt.Errorf("writing list to file: %v", err))
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
				cfmt.Printf("The interfaces list was {{successfully}}::green saved to file {{'%s'}}::blue\n", jsonFile.Name())
			} else {
				fmt.Printf("Showing %d interfaces out of %d starting from %d\n\n", len(ifaces), walkers.GlobalCtx.Classes.CountIfaces(), offset+1)
				data := representator.GetTableClassesRepr(ifaces, offset)
				fmt.Println(data)
			}
		},
	}

	listTraitsExecutor := &shell.Executor{
		Name: "traits",
		Help: "shows list of traits",
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
				Name:      "--sort",
				WithValue: true,
				Help:      "column number by which sorting will be performed",
				Default:   "2",
			},
			&flags.Flag{
				Name: "-r",
				Help: "reverse sort",
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
			sortColumn := c.GetIntFlagValue("--sort")
			reverseSort := c.Flags.Contains("-r")

			toJson, jsonFile := handleOutputInJson(c)

			traits := getter.GetClassesByOption(walkers.GlobalCtx.Classes, getter.ClassesGetOptions{
				OnlyTraits:  true,
				Count:       count,
				Offset:      offset,
				SortColumn:  sortColumn,
				ReverseSort: reverseSort,
			})

			if toJson {
				data, err := representator.GetPrettifyJsonClassesRepr(traits)
				if err != nil {
					c.Error(fmt.Errorf("writing list to file: %v", err))
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
				cfmt.Printf("The traits list was {{successfully}}::green saved to file {{'%s'}}::blue\n", jsonFile.Name())
			} else {
				fmt.Printf("Showing %d traits out of %d starting from %d\n\n", len(traits), walkers.GlobalCtx.Classes.CountTraits(), offset+1)
				data := representator.GetTableClassesRepr(traits, offset)
				fmt.Println(data)
			}
		},
	}

	listNamespacesByLevelExecutor := &shell.Executor{
		Name: "namespaces",
		Help: "shows list of namespaces on specific level",
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
				Name:      "-l",
				WithValue: true,
				Help:      "level of namespaces",
				Default:   "0",
			},
			&flags.Flag{
				Name:      "--sort",
				WithValue: true,
				Help:      "column number by which sorting will be performed",
				Default:   "2",
			},
			&flags.Flag{
				Name: "-r",
				Help: "reverse sort",
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
			level := c.GetIntFlagValue("-l")
			sortColumn := c.GetIntFlagValue("--sort")
			reverseSort := c.Flags.Contains("-r")

			toJson, jsonFile := handleOutputInJson(c)

			nss := getter.GetNamespacesByOptions(walkers.GlobalCtx.Namespaces, getter.NamespacesGetOptions{
				Level:       level,
				Count:       count,
				Offset:      offset,
				SortColumn:  sortColumn,
				ReverseSort: reverseSort,
			})

			if toJson {
				data, err := representator.GetPrettifyJsonNamespacesRepr(nss)
				if err != nil {
					c.Error(fmt.Errorf("writing list to file: %v", err))
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
				cfmt.Printf("The namespaces list was {{successfully}}::green saved to file {{'%s'}}::blue\n", jsonFile.Name())
			} else {
				fmt.Printf("Showing %d namespaces (level %d) out of %d starting from %d\n\n", len(nss), level, len(walkers.GlobalCtx.Namespaces.GetNamespacesWithSpecificLevel(level, 100000, 0)), offset+1)
				data := representator.GetTableNamespacesRepr(nss, offset)
				fmt.Println(data)
			}
		},
	}

	listExecutor := &shell.Executor{
		Name: "list",
		Help: "shows list",
		Func: func(c *shell.Context) {
			c.ShowHelpPage()
		},
	}

	listExecutor.AddExecutor(listFuncExecutor)
	listExecutor.AddExecutor(listMethodExecutor)
	listExecutor.AddExecutor(listFilesExecutor)
	listExecutor.AddExecutor(listClassesExecutor)
	listExecutor.AddExecutor(listInterfaceExecutor)
	listExecutor.AddExecutor(listTraitsExecutor)
	listExecutor.AddExecutor(listNamespacesByLevelExecutor)

	return listExecutor
}
