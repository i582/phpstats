package commands

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/i582/phpstats/internal/representator"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
)

func List() *shell.Executor {
	listFuncExecutor := &shell.Executor{
		Name: "funcs",
		Help: "show list of functions",
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
				Name: "--output",
				Help: "output json file",
			},
		),
		Func: func(c *shell.Context) {
			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			withEmbeddedFuncs := c.Flags.Contains("-e")

			var f *os.File
			output := c.GetFlagValue("--output")
			if output != "" {
				var err error
				f, err = c.ValidateFile("--output")
				if err != nil {
					log.Fatal(err)
				}
			}

			funcs := walkers.GlobalCtx.Functions.GetAll(false, true, false, count, offset, true, withEmbeddedFuncs)

			if f != nil {
				data, err := representator.GetPrettifyJsonFunctionsRepr(funcs)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(f, data)
				f.Close()
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
		Help: "show list of methods",
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

			var f *os.File
			output := c.GetFlagValue("--output")
			if output != "" {
				var err error
				f, err = c.ValidateFile("--output")
				if err != nil {
					log.Fatal(err)
				}
			}

			funcs := walkers.GlobalCtx.Functions.GetAll(true, false, false, count, offset, true, false)

			if f != nil {
				data, err := representator.GetPrettifyJsonFunctionsRepr(funcs)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(f, data)
				f.Close()
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
		Help: "show list of files",
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
				Name:      "--output",
				Help:      "output json file",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			full := c.Flags.Contains("-f")

			countValue := c.GetFlagValue("-c")
			count, _ := strconv.ParseInt(countValue, 0, 64)

			offsetValue := c.GetFlagValue("-o")
			offset, _ := strconv.ParseInt(offsetValue, 0, 64)

			var f *os.File
			output := c.GetFlagValue("--output")
			if output != "" {
				var err error
				f, err = c.ValidateFile("--output")
				if err != nil {
					log.Fatal(err)
				}
			}

			files := walkers.GlobalCtx.Files.GetAll(count, offset, true)

			if f != nil {
				data, err := representator.GetPrettifyJsonFilesRepr(files)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(f, data)
				f.Close()
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
		Help: "show list of classes",
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

			var f *os.File
			output := c.GetFlagValue("--output")
			if output != "" {
				var err error
				f, err = c.ValidateFile("--output")
				if err != nil {
					log.Fatal(err)
				}
			}

			classes := walkers.GlobalCtx.Classes.GetAllClasses(count, offset, true)

			if f != nil {
				data, err := representator.GetPrettifyJsonClassesRepr(classes)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(f, data)
				f.Close()
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
		Help: "show list of interfaces",
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

			var f *os.File
			output := c.GetFlagValue("--output")
			if output != "" {
				var err error
				f, err = c.ValidateFile("--output")
				if err != nil {
					log.Fatal(err)
				}
			}

			ifaces := walkers.GlobalCtx.Classes.GetAllInterfaces(count, offset, true)

			if f != nil {
				data, err := representator.GetPrettifyJsonClassesRepr(ifaces)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(f, data)
				f.Close()
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
		Help: "show list of namespaces on specific level",
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

			levelValue := c.GetFlagValue("-l")
			level, _ := strconv.ParseInt(levelValue, 0, 64)

			var f *os.File
			output := c.GetFlagValue("--output")
			if output != "" {
				var err error
				f, err = c.ValidateFile("--output")
				if err != nil {
					log.Fatal(err)
				}
			}

			nss := walkers.GlobalCtx.Namespaces.GetNamespacesWithSpecificLevel(level)

			if count+offset < int64(len(nss)) {
				nss = nss[:count+offset]
			}

			if offset < int64(len(nss)) {
				nss = nss[offset:]
			}

			if f != nil {
				data, err := representator.GetPrettifyJsonNamespacesRepr(nss)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(f, data)
				f.Close()
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
		Help: "list of",
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
