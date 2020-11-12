package commands

import (
	"fmt"

	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/stats/walkers"
)

func Relation() *shell.Executor {
	relationFuncsExecutor := &shell.Executor{
		Name:    "funcs",
		Aliases: []string{"methods"},
		Help:    "show relation between functions or methods",
		Flags: flags.NewFlags(
			&flags.Flag{
				Name:      "--parent",
				Help:      "parent function",
				WithValue: true,
			},
			&flags.Flag{
				Name:      "--child",
				Help:      "children function",
				WithValue: true,
			},
		),
		Func: func(c *shell.Context) {
			parentName := c.GetFlagValue("--parent")
			childName := c.GetFlagValue("--child")

			parentFun, err := walkers.GlobalCtx.GetFunction(parentName)
			if err != nil {
				c.Error(err)
				return
			}

			childFun, err := walkers.GlobalCtx.GetFunction(childName)
			if err != nil {
				c.Error(err)
				return
			}

			_, callstacks := CalledInCallstack(parentFun, childFun, nil, map[*symbols.Function]struct{}{})

			fmt.Print("Reachability: ")

			if len(callstacks) == 0 {
				fmt.Println("false")
				return
			}

			fmt.Print("true\n\n")
			fmt.Println("Callstacks:")

			for _, callstack := range callstacks {
				fmt.Print("[")
				for i, f := range callstack {
					fmt.Print(f.Name)
					if i != len(callstack)-1 {
						fmt.Print(" -> ")
					}
				}
				fmt.Println("]")
			}
		},
	}

	relationExecutor := &shell.Executor{
		Name:  "relation",
		Help:  "shows relation of",
		Flags: flags.NewFlags(),
		Func: func(c *shell.Context) {
			c.ShowHelpPage()
		},
	}

	relationExecutor.AddExecutor(relationFuncsExecutor)

	return relationExecutor
}

func CalledInCallstack(parent, child *symbols.Function, callstack []*symbols.Function, visited map[*symbols.Function]struct{}) (bool, [][]*symbols.Function) {
	if parent.Called.Len() == 0 {
		return false, nil
	}

	if callstack == nil {
		callstack = []*symbols.Function{parent}
	}

	if parent == child {
		return true, [][]*symbols.Function{callstack}
	}

	var callstacks [][]*symbols.Function

	for _, called := range parent.Called.Funcs {
		newCallstack := copyCallstack(callstack)
		newVisited := copyVisited(visited)

		newCallstack = append(newCallstack, called)

		if _, ok := newVisited[called]; ok {
			continue
		}
		if called == parent {
			continue
		}

		newVisited[called] = struct{}{}

		if called == child {
			callstacks = append(callstacks, newCallstack)
			continue
		}

		call, callstack := CalledInCallstack(called, child, newCallstack, newVisited)
		if call {
			callstacks = append(callstacks, callstack...)
		}
	}

	return len(callstacks) != 0, callstacks
}

func copyCallstack(callstack []*symbols.Function) []*symbols.Function {
	tmp := make([]*symbols.Function, len(callstack))
	copy(tmp, callstack)
	return tmp
}

func copyVisited(visited map[*symbols.Function]struct{}) map[*symbols.Function]struct{} {
	targetMap := make(map[*symbols.Function]struct{})

	for key, value := range visited {
		targetMap[key] = value
	}

	return targetMap
}
