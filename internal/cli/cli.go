package cli

import (
	"fmt"
	"log"
	"os"

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

	log.SetFlags(0)
	if err := run(); err != nil {
		log.Printf("%s: run %q error: %+v", tool.Name, subcmd, err)
		return
	}

	s := shell.NewShell()

	s.AddExecutor(commands.Info())
	s.AddExecutor(commands.List())
	s.AddExecutor(commands.Graph())
	s.AddExecutor(commands.Brief())

	s.Run()
}
