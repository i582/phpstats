package main

import (
	"fmt"
	"log"

	"github.com/i582/phpstats/internal/cli"
	"github.com/i582/phpstats/internal/stats"
)

func main() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)

	cli.RunPhplinterTool(&cli.PhplinterTool{
		Name:    "stats",
		Collect: stats.CollectMain,
		Process: nil,
	})

	fmt.Println(1)
}
