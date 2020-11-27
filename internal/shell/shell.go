package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gookit/color"
)

type Shell struct {
	Execs Executors

	Active bool
}

func (s *Shell) Error(msg string) {
	color.Red.Printf("Error: %v\n", msg)
}

func (s *Shell) AddExecutor(exec *Executor) {
	if s.Execs == nil {
		s.Execs = Executors{}
	}

	s.Execs[exec.Name] = exec
}

func (s *Shell) GetExecutor(name string) (*Executor, bool) {
	if s.Execs == nil {
		return nil, false
	}

	exc, ok := s.Execs[name]
	return exc, ok
}

func NewShell() *Shell {
	shell := &Shell{
		Active: true,
	}

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

	shell.AddExecutor(&Executor{
		Name:    "exit",
		Aliases: []string{"quit"},
		Help:    "exit the program",
		Func: func(c *Context) {
			shell.Active = false
		},
	})

	return shell
}

func (s *Shell) Run() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Entering interactive mode (type \"help\" for commands)")

	for s.Active {
		color.Yellow.Print(`>>> `)
		ln, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		line := string(ln)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		tokens := strings.FieldsFunc(line, func(r rune) bool {
			return r == '=' || r == ' '
		})
		if len(tokens) == 0 {
			continue
		}

		command := tokens[0]

		e, has := s.Execs[command]
		if !has {
			s.Error(fmt.Sprintf("command %s not found", command))
			continue
		}

		e.Execute(&Context{
			Args:  tokens[1:],
			Flags: e.Flags,
			Exec:  e,
		})
	}
}
