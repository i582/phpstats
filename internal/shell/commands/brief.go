package commands

import (
	"fmt"

	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
)

func Brief() *shell.Executor {
	briefExecutor := &shell.Executor{
		Name:  "brief",
		Help:  "shows general information",
		Flags: flags.NewFlags(),
		Func: func(c *shell.Context) {
			var countMethods int64
			var countFuncs int64

			for _, fn := range walkers.GlobalCtx.Functions.Funcs {
				if fn.IsMethod() {
					countMethods++
					continue
				}

				if !fn.IsMethod() {
					countFuncs++
				}

				if fn.IsEmbeddedFunc() {
					continue
				}
			}

			var countLines int64
			for _, file := range walkers.GlobalCtx.Files.Files {
				countLines += file.CountLines
			}

			fmt.Printf("General project statistics\n")
			fmt.Printf("Classes:       %d\n", walkers.GlobalCtx.Classes.Len())
			fmt.Printf("  Methods:     %d\n", countMethods)
			fmt.Printf("  Constants:   %d\n", walkers.GlobalCtx.Constants.Len())
			fmt.Printf("Functions:     %d\n", countFuncs)
			fmt.Printf("Files:         %d\n", walkers.GlobalCtx.Files.Len())
			fmt.Printf("Lines of code: %d\n", countLines)

			fmt.Println()
		},
	}

	return briefExecutor
}
