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

	s := shell.NewShell()
	s.Run()

	// // create new shell.
	// // by default, new shell includes 'exit', 'help' and 'clear' commands.
	// shell := ishell.New()
	//
	// // display welcome info.
	// shell.Println("Sample Interactive Shell")
	//
	// // register a function for "greet" command.
	// shell.AddCmd(&ishell.Cmd{
	// 	Name: "greet",
	// 	Help: "greet user",
	// 	Func: func(c *ishell.Context) {
	// 		c.Println("Hello", strings.Join(c.Args, " "))
	// 	},
	// })
	//
	// // run shell
	// shell.Run()
}
