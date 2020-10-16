package main

import (
	"phpstats/internal/cli"
	"phpstats/internal/shell"
	"phpstats/internal/shell/commands"
	"phpstats/internal/stats"
)

func main() {
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
