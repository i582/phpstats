package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/gookit/color"

	"github.com/i582/phpstats/internal/server"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/commands"
	"github.com/i582/phpstats/internal/stats/walkers"
	"github.com/i582/phpstats/internal/utils"

	"github.com/urfave/cli/v2"
)

var MainShell = shell.NewShell()

func Run() {
	log.SetFlags(0)

	if len(os.Args) > 1 {
		subcmd := os.Args[1]
		// Remove sub command from os.Args.
		os.Args = append(os.Args[:1], os.Args[2:]...)
		os.Args[0] = "phpstats/" + subcmd
	}

	MainShell.AddExecutor(commands.Info())
	MainShell.AddExecutor(commands.List())
	MainShell.AddExecutor(commands.Graph())
	MainShell.AddExecutor(commands.Brief())
	MainShell.AddExecutor(commands.About())
	MainShell.AddExecutor(commands.Top())
	MainShell.AddExecutor(commands.Metrics())

	var cacheDir string
	var port int64
	app := &cli.App{
		Name:  "collect",
		Usage: "data collection",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:  "port",
				Usage: "port used by the server.",
				Value: 8080,
			},
			&cli.StringFlag{
				Name:        "cache-dir",
				Usage:       "custom directory for cache storage.",
				Value:       utils.DefaultCacheDir(),
				Destination: &cacheDir,
			},
			&cli.StringFlag{
				Name:        "project-path",
				Usage:       "path to the project relative to which all imports are allowed.",
				Destination: &walkers.GlobalCtx.ProjectRoot,
			},
		},
		Action: func(c *cli.Context) error {
			if len(os.Args) == 1 {
				commands.About().Execute(&shell.Context{})
				fmt.Printf("\nUsage\n\t$ phpstats collect [--port <value>] [--project-path <dir>] [--cache-dir <dir>] <analyze-dir>\n")
				return fmt.Errorf("empty")
			}

			server.RunServer(port)

			// Normalize flags for NoVerify
			exe := os.Args[0]
			path := os.Args[len(os.Args)-1]
			os.Args = []string{exe, "-cache-dir", cacheDir, path}

			if c.NArg() > 1 {
				log.Fatalf(color.Red.Sprintf("Error: too many arguments"))
			}

			if c.NArg() < 1 {
				log.Fatalf(color.Red.Sprintf("Error: too few arguments"))
			}

			err := walkers.Collect()
			if err != nil {
				return err
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return
	}

	MainShell.Run()
}
