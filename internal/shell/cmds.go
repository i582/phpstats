package shell

//
// import (
// 	"fmt"
// 	"strconv"
//
// 	"github.com/abiosoft/ishell"
//
// 	"phpstats/internal/stats"
// )
//
// func Run() {
// 	shell := ishell.New()
// 	shell.Println("Sample Interactive Shell")
//
// 	classInfoCmd := &ishell.Cmd{
// 		Name: "class",
// 		Help: "info about some class",
// 		Func: func(c *ishell.Context) {
// 			flags, args := getFlags(c.Args)
// 			full := flags.Contains("-f")
//
// 			if len(args) != 1 {
// 				c.Err(fmt.Errorf("команда принимает ровно один аргумент"))
// 				return
// 			}
//
// 			classNames, err := stats.GlobalCtx.Classes.GetFullClassName(args[0])
// 			if err != nil {
// 				c.Err(err)
// 				return
// 			}
//
// 			var className string
//
// 			if len(classNames) > 1 {
// 				choice := c.MultiChoice(classNames, "Какой класс вы имели ввиду?")
// 				className = classNames[choice]
//
// 				c.Println()
// 			} else {
// 				className = classNames[0]
// 			}
//
// 			class, _ := stats.GlobalCtx.Classes.Get(className)
//
// 			if full {
// 				fmt.Println(class.FullString(0))
// 			} else {
// 				fmt.Println(class.ShortString(0))
// 			}
// 		},
// 	}
//
// 	funcInfoCmd := &ishell.Cmd{
// 		Name: "func",
// 		Help: "info about some func",
// 		Func: func(c *ishell.Context) {
// 			flags, args := getFlags(c.Args)
// 			full := flags.Contains("-f")
//
// 			if len(args) != 1 {
// 				c.Err(fmt.Errorf("команда принимает ровно один аргумент"))
// 				return
// 			}
//
// 			funcNameKeys, err := stats.GlobalCtx.Funcs.GetFullFuncName(args[0])
// 			if err != nil {
// 				c.Err(err)
// 				return
// 			}
//
// 			var funcKeyIndex int
//
// 			if len(funcNameKeys) > 1 {
// 				funcManes := make([]string, 0, len(funcNameKeys))
// 				for _, key := range funcNameKeys {
// 					funcManes = append(funcManes, key.String())
// 				}
//
// 				funcKeyIndex = c.MultiChoice(funcManes, "Какую функцию вы имели ввиду?")
// 				c.Println()
// 			} else {
// 				funcKeyIndex = 0
// 			}
//
// 			fn, _ := stats.GlobalCtx.Funcs.Get(funcNameKeys[funcKeyIndex])
//
// 			if full {
// 				fmt.Println(fn.FullString())
// 			} else {
// 				fmt.Println(fn.ShortString())
// 			}
// 		},
// 	}
//
// 	fileInfoCmd := &ishell.Cmd{
// 		Name: "file",
// 		Help: "info about some file",
// 		Func: func(c *ishell.Context) {
// 			flags, args := getFlags(c.Args)
// 			full := flags.Contains("-f")
// 			recursiveFlag, recursive := flags.Get("-r")
//
// 			if len(args) != 1 {
// 				c.Err(fmt.Errorf("команда принимает ровно один аргумент"))
// 				return
// 			}
//
// 			patches, err := stats.GlobalCtx.Files.GetFullFileName(args[0])
// 			if err != nil {
// 				c.Err(err)
// 				return
// 			}
//
// 			var patch string
//
// 			if len(patches) > 1 {
// 				choice := c.MultiChoice(patches, "Какой файл вы имели ввиду?")
// 				patch = patches[choice]
//
// 				c.Println()
// 			} else {
// 				patch = patches[0]
// 			}
//
// 			file, _ := stats.GlobalCtx.Files.Get(patch)
//
// 			if recursive {
// 				count, err := strconv.ParseInt(recursiveFlag.Value, 0, 64)
// 				if err != nil {
// 					c.Err(fmt.Errorf("значение флага должно быть числом"))
// 				}
//
// 				fmt.Println(file.FullStringRecursive(int(count)))
// 			}
//
// 			if full {
// 				fmt.Println(file.FullString(0))
// 			} else {
// 				fmt.Println(file.ShortString(0))
// 			}
// 		},
// 	}
//
// 	infoCmd := &ishell.Cmd{
// 		Name: "info",
// 		Help: "info about something",
// 		Func: func(c *ishell.Context) {
//
// 		},
// 	}
//
// 	infoCmd.AddCmd(classInfoCmd)
// 	infoCmd.AddCmd(funcInfoCmd)
// 	infoCmd.AddCmd(fileInfoCmd)
//
// 	shell.AddCmd(infoCmd)
//
// 	shell.Run()
// }
