package commands

import (
	"fmt"
	"log"

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
				Name: "-e",
				Help: "show embedded functions",
			},
			&flags.Flag{
				Name:      "--json",
				Help:      "output to json file",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			count := c.GetIntFlagValue("-c")
			offset := c.GetIntFlagValue("-o")

			withEmbeddedFuncs := c.Flags.Contains("-e")
			toJson, jsonFile := handleOutputInJson(c)

			funcs := walkers.GlobalCtx.Functions.GetAll(false, true, false, count, offset, true, withEmbeddedFuncs)

			if toJson {
				data, err := representator.GetPrettifyJsonFunctionsRepr(funcs)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
			} else {
				for _, fn := range funcs {
					data := representator.GetStringFunctionRepr(fn)
					fmt.Println(data)
				}
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
				Name:      "--json",
				Help:      "output to json file",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			count := c.GetIntFlagValue("-c")
			offset := c.GetIntFlagValue("-o")

			toJson, jsonFile := handleOutputInJson(c)

			funcs := walkers.GlobalCtx.Functions.GetAll(true, false, false, count, offset, true, false)

			if toJson {
				data, err := representator.GetPrettifyJsonFunctionsRepr(funcs)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
			} else {
				for _, fn := range funcs {
					data := representator.GetStringFunctionRepr(fn)
					fmt.Println(data)
				}
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
				Name: "-f",
				Help: "show full information",
			},
			&flags.Flag{
				Name:      "--json",
				Help:      "output to json file",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")
			count := c.GetIntFlagValue("-c")
			offset := c.GetIntFlagValue("-o")

			toJson, jsonFile := handleOutputInJson(c)

			files := walkers.GlobalCtx.Files.GetAll(count, offset, true)

			if toJson {
				data, err := representator.GetPrettifyJsonFilesRepr(files)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
			} else {
				for _, file := range files {
					var data string
					if full {
						data = representator.GetStringFileRepr(file)
					} else {
						data = representator.GetShortStringFileRepr(file)
					}
					fmt.Println(data)
				}
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
				Name:      "--json",
				Help:      "output to json file",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			count := c.GetIntFlagValue("-c")
			offset := c.GetIntFlagValue("-o")

			toJson, jsonFile := handleOutputInJson(c)

			classes := walkers.GlobalCtx.Classes.GetAllClasses(count, offset, true)

			if toJson {
				data, err := representator.GetPrettifyJsonClassesRepr(classes)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
			} else {
				for _, class := range classes {
					data := representator.GetStringClassRepr(class)
					fmt.Println(data)
				}
			}
		},
	}

	listInterfaceExecutor := &shell.Executor{
		Name: "interfaces",
		Help: "shows list of interfaces",
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
				Name: "-f",
				Help: "show full information",
			},
			&flags.Flag{
				Name:      "--json",
				Help:      "output to json file",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			count := c.GetIntFlagValue("-c")
			offset := c.GetIntFlagValue("-o")

			toJson, jsonFile := handleOutputInJson(c)

			ifaces := walkers.GlobalCtx.Classes.GetAllInterfaces(count, offset, true)

			if toJson {
				data, err := representator.GetPrettifyJsonClassesRepr(ifaces)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
			} else {
				for _, iface := range ifaces {
					data := representator.GetStringClassRepr(iface)
					fmt.Println(data)
				}
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
				Name: "-f",
				Help: "show full information",
			},
			&flags.Flag{
				Name:      "--json",
				Help:      "output to json file",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			count := c.GetIntFlagValue("-c")
			offset := c.GetIntFlagValue("-o")
			level := c.GetIntFlagValue("-l")

			toJson, jsonFile := handleOutputInJson(c)

			nss := walkers.GlobalCtx.Namespaces.GetNamespacesWithSpecificLevel(level, count, offset)

			if toJson {
				data, err := representator.GetPrettifyJsonNamespacesRepr(nss)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(jsonFile, data)
				jsonFile.Close()
			} else {
				for _, ns := range nss {
					data := representator.GetStringNamespaceRepr(ns)
					fmt.Println(data)
				}
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
	listExecutor.AddExecutor(listNamespacesByLevelExecutor)

	return listExecutor
}
