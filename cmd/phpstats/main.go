package main

import (
	"phpstats/internal/cli"
	"phpstats/internal/shell"
	"phpstats/internal/stats"
)

func main() {
	cli.RunPhplinterTool(&cli.PhplinterTool{
		Name:    "stats",
		Collect: stats.CollectMain,
		Process: nil,
	})

	shell.Init()
}
