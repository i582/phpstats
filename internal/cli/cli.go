package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/i582/phpstats/internal/server"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/commands"
)

type PhplinterTool struct {
	Name    string
	Collect func() error
	Process func() error
}

func errorFunc(format string, args ...interface{}) func() error {
	return func() error {
		return fmt.Errorf(format, args...)
	}
}

var MainShell = shell.NewShell()

func RunPhplinterTool(tool *PhplinterTool) {
	subcmd := ""
	if len(os.Args) > 1 {
		subcmd = os.Args[1]
		// Remove sub command from os.Args.
		os.Args = append(os.Args[:1], os.Args[2:]...)
		os.Args[0] = tool.Name + "/" + subcmd
	}

	run := errorFunc("sub-command %q not found", subcmd)
	switch subcmd {
	case "collect":
		if tool.Collect == nil {
			run = errorFunc("collect command is nil")
		} else {
			run = tool.Collect
		}
	case "process":
		if tool.Process == nil {
			run = errorFunc("process command is nil")
		} else {
			run = tool.Process
		}
	}

	server.RunServer()

	log.SetFlags(0)
	if err := run(); err != nil {
		log.Printf("%MainShell: run %q error: %+v", tool.Name, subcmd, err)
		return
	}

	MainShell.AddExecutor(commands.Info())
	MainShell.AddExecutor(commands.List())
	MainShell.AddExecutor(commands.Graph())
	MainShell.AddExecutor(commands.Brief())
	MainShell.AddExecutor(commands.Top())

	MainShell.Run()
}
