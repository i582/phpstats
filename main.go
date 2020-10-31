package main

import (
	"log"

	"github.com/i582/phpstats/internal/cli"
	"github.com/i582/phpstats/internal/stats/walkers"
)

func main() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)

	cli.RunPhplinterTool(&cli.PhplinterTool{
		Name:    "stats",
		Collect: walkers.CollectMain,
		Process: nil,
	})
}
