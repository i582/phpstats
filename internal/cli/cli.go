package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/gookit/color"

	"github.com/i582/phpstats/internal/config"
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

	var cacheDir string
	var configPath string
	var disableCache bool
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
			&cli.BoolFlag{
				Name:        "disable-cache",
				Usage:       "",
				Destination: &disableCache,
			},
			&cli.StringFlag{
				Name:        "config-path",
				Usage:       "path to the config.",
				Destination: &configPath,
			},
		},
		Action: func(c *cli.Context) error {
			if len(os.Args) == 1 {
				commands.About().Execute(&shell.Context{})
				fmt.Printf("\nUsage\n\t$ phpstats collect [--port <value>] [--project-path <dir>] [--cache-dir <dir>] <analyze-dir>\n\n")
				return fmt.Errorf("empty")
			}

			cfg, err := config.OpenConfig(configPath)
			if err == nil {
				if cfg.CacheDir == "" {
					cfg.CacheDir = utils.DefaultCacheDir()
				}
			} else {
				cfg = &config.Config{
					Port:         port,
					CacheDir:     cacheDir,
					DisableCache: disableCache,
					ProjectPath:  walkers.GlobalCtx.ProjectRoot,
					Exclude:      nil,
					Groups:       nil,
					Extensions:   nil,
				}
			}

			server.RunServer(port)

			// Normalize flags for NoVerify
			exe := os.Args[0]
			path := os.Args[len(os.Args)-1]

			cfgCli := cfg.ToCliArgs()
			os.Args = []string{exe}
			os.Args = append(os.Args, cfgCli...)
			os.Args = append(os.Args, path)

			fmt.Print(os.Args)

			if c.NArg() > 1 {
				log.Fatalf(color.Red.Sprintf("Error: too many arguments"))
			}

			if c.NArg() < 1 {
				log.Fatalf(color.Red.Sprintf("Error: too few arguments"))
			}

			err = walkers.Collect()
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
