package main

import (
	"log"

	"github.com/i582/phpstats/internal/cli"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/commands"
	"github.com/i582/phpstats/internal/stats"
)

func main() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)

	cli.RunPhplinterTool(&cli.PhplinterTool{
		Name:    "stats",
		Collect: stats.CollectMain,
		Process: nil,
	})

	s := shell.NewShell()

	s.AddExecutor(commands.Info())
	s.AddExecutor(commands.List())
	s.AddExecutor(commands.Graph())

	s.Run()
}
