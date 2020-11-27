package shell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/gookit/color"

	flags2 "github.com/i582/phpstats/internal/shell/flags"
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
			os.Exit(0)
		},
	})

	return shell
}

func (e Executors) GetSuggests(commands string) []prompt.Suggest {
	parts := strings.Fields(commands)
	if len(parts) == 0 {
		parts = []string{commands}
	}
	if strings.HasSuffix(commands, " ") {
		parts = append(parts, "")
	}

	return e.getSuggests(parts)
}

func (e Executors) getSuggests(commands []string) []prompt.Suggest {
	var suggests []prompt.Suggest
	mainCommand := commands[0]

	if len(commands) == 1 {
		executors := make([]*Executor, 0, len(e))
		for _, executor := range e {
			executors = append(executors, executor)
		}
		sort.Slice(executors, func(i, j int) bool {
			return executors[i].Name < executors[j].Name
		})

		for _, executor := range executors {
			if strings.HasPrefix(executor.Name, mainCommand) {
				suggests = append(suggests, prompt.Suggest{
					Text:        executor.Name,
					Description: executor.Help,
				})
			}
		}

		return suggests
	}

	command, found := e[mainCommand]
	if found {
		return command.getSuggests(commands[1:])
	}

	return suggests
}

func (e Executor) getSuggests(commands []string) []prompt.Suggest {
	mainCommand := commands[0]

	if len(commands) == 1 {
		var suggests []prompt.Suggest

		executors := make([]*Executor, 0, len(e.SubExecs))
		for _, executor := range e.SubExecs {
			executors = append(executors, executor)
		}
		sort.Slice(executors, func(i, j int) bool {
			return executors[i].Name < executors[j].Name
		})

		for _, executor := range executors {
			if strings.HasPrefix(executor.Name, mainCommand) {
				suggests = append(suggests, prompt.Suggest{
					Text:        executor.Name,
					Description: executor.Help,
				})
			}
		}

		if len(executors) == 0 {
			if e.Flags != nil {
				flags := make([]*flags2.Flag, 0, len(e.Flags.Flags))
				for _, flag := range e.Flags.Flags {
					flags = append(flags, flag)
				}
				sort.Slice(flags, func(i, j int) bool {
					return flags[i].Name < flags[j].Name
				})

				for _, flag := range flags {
					if strings.HasPrefix(flag.Name, mainCommand) {
						suggests = append(suggests, prompt.Suggest{
							Text:        flag.Name,
							Description: flag.Help,
						})
					}
				}
			}
		}

		return suggests
	}

	command, found := e.SubExecs[mainCommand]
	if found {
		return command.getSuggests(commands[1:])
	} else {
		return e.getSuggests(commands[1:])
	}

	return nil
}

func (s *Shell) completer(d prompt.Document) []prompt.Suggest {
	return s.Execs.GetSuggests(d.CurrentLine())
}

func (s *Shell) executor(in string) {
	line := strings.TrimSpace(in)
	if line == "" {
		return
	}

	tokens := strings.FieldsFunc(line, func(r rune) bool {
		return r == '=' || r == ' '
	})
	if len(tokens) == 0 {
		return
	}

	command := tokens[0]

	e, has := s.Execs[command]
	if !has {
		s.Error(fmt.Sprintf("command %s not found", command))
		return
	}

	e.Execute(&Context{
		Args:  tokens[1:],
		Flags: e.Flags,
		Exec:  e,
	})
}

func (s *Shell) Run() {
	// reader := bufio.NewReader(os.Stdin)

	fmt.Println("Entering interactive mode (type \"help\" for commands)")

	p := prompt.New(
		s.executor,
		s.completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionSuggestionBGColor(prompt.Yellow),
		prompt.OptionSuggestionTextColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),

		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSelectedDescriptionTextColor(prompt.White),
	)
	p.Run()
	//
	// for s.Active {
	// 	color.Yellow.Print(`>>> `)
	// 	ln, _, err := reader.ReadLine()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	line := string(ln)
	// 	line = strings.TrimSpace(line)
	// 	if line == "" {
	// 		continue
	// 	}
	//
	// 	tokens := strings.FieldsFunc(line, func(r rune) bool {
	// 		return r == '=' || r == ' '
	// 	})
	// 	if len(tokens) == 0 {
	// 		continue
	// 	}
	//
	// 	command := tokens[0]
	//
	// 	e, has := s.Execs[command]
	// 	if !has {
	// 		s.Error(fmt.Sprintf("command %s not found", command))
	// 		continue
	// 	}
	//
	// 	e.Execute(&Context{
	// 		Args:  tokens[1:],
	// 		Flags: e.Flags,
	// 		Exec:  e,
	// 	})
	// }
}
