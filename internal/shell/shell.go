package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Shell struct {
	Execs Executors
}

func (s *Shell) AddExecutor(exec *Executor) {
	if s.Execs == nil {
		s.Execs = Executors{}
	}

	s.Execs[exec.Name] = exec
}

func NewShell() *Shell {
	shell := &Shell{}

	shell.AddExecutor(&Executor{
		Name: "help",
		Help: "help page",
		Func: func(c *Context) {
			fmt.Println("Commands:")
			for _, e := range shell.Execs {
				fmt.Print(e.HelpPage(0))
			}
		},
	})

	shell.AddExecutor(&Executor{
		Name: "clear",
		Help: "clear screen",
		Func: func(c *Context) {
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/c", "cls")
			} else {
				cmd = exec.Command("clear")
			}
			cmd.Stdout = os.Stdout
			_ = cmd.Run()
		},
	})

	// return shell

	// execs := map[string]*Executor{}
	// execs["info"] = &Executor{Func: func(args []string) {
	// 	if len(args) < 2 {
	// 		return
	// 	}
	//
	// 	switch args[0] {
	// 	case "class":
	// 		full := args[1] == "-f"
	//
	// 		if full {
	// 			args = args[1:]
	// 		}
	//
	// 		class, ok := stats.GlobalCtx.Classes.Get(args[1])
	// 		if !ok {
	// 			fmt.Printf("Класс %s не найден!\n", args[1])
	// 			return
	// 		}
	//
	// 		if full {
	// 			fmt.Println(class.FullString(0))
	// 		} else {
	// 			fmt.Println(class.ShortString(0))
	// 		}
	//
	// 	case "func":
	// 		names, err := stats.GlobalCtx.Funcs.GetFullFuncName(args[1])
	// 		if err != nil {
	// 			fmt.Printf("Функция %s не найдена!\n", args[1])
	// 			return
	// 		}
	//
	// 		if len(names) != 1 {
	// 			fmt.Printf("Найдено несколько функций с похожим именем\n")
	// 		}
	//
	// 		if len(names) == 1 {
	// 			fn, _ := stats.GlobalCtx.Funcs.Get(names[0])
	// 			fmt.Print(fn)
	// 			return
	// 		}
	//
	// 		for _, name := range names {
	// 			fn, _ := stats.GlobalCtx.Funcs.Get(name)
	// 			fmt.Print(fn)
	// 		}
	//
	// 	case "file":
	// 		recursive := args[1] == "-r"
	// 		full := args[1] == "-f"
	//
	// 		if recursive || full {
	// 			args = args[1:]
	// 		}
	//
	// 		path, err := stats.GlobalCtx.Files.GetFullFileName(args[1])
	// 		if err != nil {
	// 			fmt.Printf("Файл %s не найден!\n", args[1])
	// 			return
	// 		}
	// 		file, _ := stats.GlobalCtx.Files.Get(path[0])
	//
	// 		if recursive {
	// 			fmt.Println(file.FullStringRecursive(5))
	// 		}
	//
	// 		if full {
	// 			fmt.Println(file.FullString(0))
	// 		} else {
	// 			fmt.Println(file.ShortString(0))
	// 		}
	//
	// 	default:
	// 		fmt.Printf("Нераспознанная команда: %s\n", args[0])
	// 	}
	// }}
	//
	// execs["list"] = &Executor{Func: func(args []string) {
	//
	// 	if len(args) == 0 {
	// 		return
	// 	}
	//
	// 	switch args[0] {
	// 	case "classes":
	// 		count, ok := tryEatCount(args)
	// 		if !ok {
	// 			count = 10
	// 		}
	//
	// 		var index int64
	// 		classes := stats.GlobalCtx.Classes.Classes
	//
	// 		for _, class := range classes {
	// 			fmt.Println(class.FullString(0))
	//
	// 			index++
	//
	// 			if index >= count {
	// 				break
	// 			}
	// 		}
	//
	// 	case "interfaces":
	// 		count, ok := tryEatCount(args)
	// 		if !ok {
	// 			count = 10
	// 		}
	//
	// 		var index int64
	// 		classes := stats.GlobalCtx.Classes.Classes
	//
	// 		for _, class := range classes {
	// 			if !class.IsInterface {
	// 				continue
	// 			}
	//
	// 			fmt.Println(class.FullString(0))
	//
	// 			index++
	//
	// 			if index >= count {
	// 				break
	// 			}
	// 		}
	//
	// 	case "funcs":
	// 		count, ok := tryEatCount(args)
	// 		if ok {
	// 			funcs := stats.GlobalCtx.Funcs.GetAll(false, false, true, int(count), true)
	//
	// 			for _, fn := range funcs {
	// 				fmt.Print(fn)
	// 			}
	// 		} else {
	// 			funcs := stats.GlobalCtx.Funcs.GetAll(false, false, true, 10, true)
	//
	// 			for _, fn := range funcs {
	// 				fmt.Print(fn)
	// 			}
	// 		}
	//
	// 	case "methods":
	// 		count, ok := tryEatCount(args)
	// 		if ok {
	// 			funcs := stats.GlobalCtx.Funcs.GetAll(false, false, true, int(count), true)
	//
	// 			for _, fn := range funcs {
	// 				if !fn.Name.IsMethod() {
	// 					continue
	// 				}
	//
	// 				fmt.Print(fn)
	// 			}
	// 		} else {
	// 			funcs := stats.GlobalCtx.Funcs.GetAll(false, false, true, 10, true)
	//
	// 			for _, fn := range funcs {
	// 				if !fn.Name.IsMethod() {
	// 					continue
	// 				}
	//
	// 				fmt.Print(fn)
	// 			}
	// 		}
	//
	// 	case "files":
	// 		full := args[1] == "-f"
	// 		if full {
	// 			args = args[1:]
	// 		}
	//
	// 		count, ok := tryEatCount(args)
	// 		if ok {
	// 			files := stats.GlobalCtx.Files.GetAll(int(count), true)
	//
	// 			for _, file := range files {
	// 				if full {
	// 					fmt.Print(file.FullString(0))
	// 				} else {
	// 					fmt.Print(file.ExtraShortString(0))
	// 				}
	// 			}
	// 		} else {
	// 			files := stats.GlobalCtx.Files.GetAll(10, true)
	//
	// 			for _, file := range files {
	// 				if full {
	// 					fmt.Print(file.FullString(0))
	// 				} else {
	// 					fmt.Print(file.ExtraShortString(0))
	// 				}
	// 			}
	// 		}
	// 	default:
	// 		fmt.Printf("Нераспознанная команда: %s", args[0])
	// 	}
	//
	// }}
	//
	// execs["used-in"] = &Executor{Func: func(args []string) {
	//
	// 	if len(args) == 0 {
	// 		return
	// 	}
	//
	// 	switch args[0] {
	// 	case "class":
	// 	case "func":
	// 		names, err := stats.GlobalCtx.Funcs.GetFullFuncName(args[1])
	// 		if err != nil {
	// 			fmt.Printf("Функция %s не найдена!\n", args[1])
	// 			return
	// 		}
	//
	// 		if len(names) != 1 {
	// 			fmt.Printf("Найдено несколько функций с похожим именем\n")
	// 		}
	//
	// 		if len(names) == 1 {
	// 			fn, _ := stats.GlobalCtx.Funcs.Get(names[0])
	//
	// 			classes := stats.NewClasses()
	// 			for _, called := range fn.Called.Funcs {
	// 				if called.Class != nil && called.IsMethod() {
	// 					classes.Add(called.Class)
	// 				}
	// 			}
	//
	// 			if classes.Len() != 0 {
	// 				fmt.Println("Функция вызывает методы следующих классов")
	// 			}
	// 			for _, class := range classes.Classes {
	// 				fmt.Printf("%s\n", class.ShortString(1))
	// 			}
	//
	// 			return
	// 		}
	//
	// 	case "file":
	// 		path, err := stats.GlobalCtx.Files.GetFullFileName(args[2])
	// 		if err != nil {
	// 			fmt.Printf("Файл %s не найден!\n", args[1])
	// 			return
	// 		}
	//
	// 		var res string
	//
	// 		file, _ := stats.GlobalCtx.Files.Get(path[0])
	//
	// 		count, ok := tryEatCount(args)
	// 		if ok {
	// 			args = args[1:]
	// 			res += file.FullStringRecursive(int(count))
	// 		} else {
	// 			res += file.FullStringRecursive(5)
	// 		}
	//
	// 		fmt.Println(res)
	//
	// 	default:
	// 		fmt.Printf("Нераспознанная команда: %s", args[0])
	// 	}
	//
	// }}
	//
	// execs["graph"] = &Executor{Func: func(args []string) {
	//
	// 	if len(args) == 0 {
	// 		return
	// 	}
	//
	// 	switch args[0] {
	// 	case "class":
	// 	case "func":
	// 	case "file":
	// 		path, err := stats.GlobalCtx.Files.GetFullFileName(args[3])
	// 		if err != nil {
	// 			fmt.Printf("Файл %s не найден!\n", args[3])
	// 			return
	// 		}
	//
	// 		var res string
	//
	// 		file, _ := stats.GlobalCtx.Files.Get(path[0])
	//
	// 		output := args[2]
	// 		outputFile, err := os.OpenFile(output, os.O_CREATE|os.O_RDWR, os.ModePerm)
	// 		if err != nil {
	// 			log.Fatalf("file not open %v", err)
	// 		}
	//
	// 		count, ok := tryEatCount(args)
	// 		if ok {
	// 			args = args[1:]
	// 			res += file.GraphvizRecursive(int(count))
	// 		} else {
	// 			res += file.GraphvizRecursive(5)
	// 		}
	//
	// 		fmt.Fprint(outputFile, res)
	// 		fmt.Println(res)
	// 		outputFile.Close()
	//
	// 	default:
	// 		fmt.Printf("Нераспознанная команда: %s", args[0])
	// 	}
	//
	// }}
	//
	// execs["stat"] = &Executor{Func: func(args []string) {
	//
	// 	if len(args) == 0 {
	// 		return
	// 	}
	//
	// 	switch args[0] {
	// 	case "namespace":
	// 		namespace := args[1]
	// 		classes := stats.NewClasses()
	// 		for _, class := range stats.GlobalCtx.Classes.Classes {
	// 			if strings.Contains(class.Name, namespace) {
	// 				classes.Add(class)
	// 			}
	// 		}
	//
	// 		var aff float64
	// 		var eff float64
	//
	// 		for _, class := range classes.Classes {
	// 			for _, dep := range class.Deps.Classes {
	// 				// если зависимость вне пространства имен
	// 				if !strings.Contains(dep.Name, namespace) {
	// 					aff++
	// 				}
	// 			}
	//
	// 			for _, depBy := range class.DepsBy.Classes {
	// 				// если зависимость вне пространства имен
	// 				if !strings.Contains(depBy.Name, namespace) {
	// 					eff++
	// 				}
	// 			}
	// 		}
	//
	// 		// for _, class := range classes.Classes {
	// 		// 	aff += float64(len(class.Deps.Classes))
	// 		// 	eff += float64(len(class.DepsBy.Classes))
	// 		// }
	//
	// 		var stability float64
	// 		if eff+aff == 0 {
	// 			stability = 0
	// 		} else {
	// 			stability = eff / (eff + aff)
	// 		}
	//
	// 		fmt.Println(aff)
	// 		fmt.Println(eff)
	// 		fmt.Println(stability)
	//
	// 	default:
	// 		fmt.Printf("Нераспознанная команда: %s", args[0])
	// 	}
	//
	// }}

	return shell
}

func (s *Shell) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(`>>> `)
		line, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		tokens := strings.FieldsFunc(string(line), func(r rune) bool {
			return r == '=' || r == ' '
		})
		if len(tokens) == 0 {
			continue
		}

		command := tokens[0]
		exec, has := s.Execs[command]
		if !has {
			fmt.Printf("connamd %s not found\n", command)
			continue
		}

		exec.Execute(&Context{
			Args:  tokens[1:],
			Flags: exec.Flags,
			Exec:  exec,
		})
	}
}
