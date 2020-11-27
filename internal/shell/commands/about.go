package commands

import (
	"fmt"

	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
)

func About() *shell.Executor {
	aboutExecutor := &shell.Executor{
		Name:  "about",
		Help:  "shows information about phpstats",
		Flags: flags.NewFlags(),
		Func: func(c *shell.Context) {
			fmt.Print(`About PHPStats v0.0.3

PHPStats is a utility for collecting project statistics and building 
dependency graphs for PHP, that allows you to find places in the code 
that can be improved.

It tries to be fast, ~150k LOC/s (lines of code per second) on Core i5 
with SSD with ~3500Mb/s for reading.

This tool is written in Go and uses NoVerify (https://github.com/VKCOM/noverify).

Author: Petr Makhnev (tg: @petr_makhnev)

MIT (c) 2020
`)
		},
	}

	return aboutExecutor
}
