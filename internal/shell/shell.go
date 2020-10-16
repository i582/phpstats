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
		e, has := s.Execs[command]
		if !has {
			fmt.Printf("connamd %s not found\n", command)
			continue
		}

		e.Execute(&Context{
			Args:  tokens[1:],
			Flags: e.Flags,
			Exec:  e,
		})
	}
}
