package cli

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gookit/color"

	"github.com/i582/phpstats/internal/config"
	"github.com/i582/phpstats/internal/server"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/commands"
	"github.com/i582/phpstats/internal/stats/walkers"
	"github.com/i582/phpstats/internal/utils"

	"github.com/urfave/cli/v2"
)

// Run launches the entire application.
func Run() {
	log.SetFlags(0)

	if len(os.Args) > 1 {
		subcmd := os.Args[1]
		// Remove sub command from os.Args.
		os.Args = append(os.Args[:1], os.Args[2:]...)
		os.Args[0] = "phpstats/" + subcmd
	}

	MainShell := shell.NewShell()

	MainShell.AddExecutor(commands.Info())
	MainShell.AddExecutor(commands.List())
	MainShell.AddExecutor(commands.Graph())
	MainShell.AddExecutor(commands.Brief())
	MainShell.AddExecutor(commands.About())
	MainShell.AddExecutor(commands.Top())
	MainShell.AddExecutor(commands.Metrics())
	MainShell.AddExecutor(commands.Relation())

	var cacheDir string
	var configPath string
	var disableCache bool
	var port int64

	app := &cli.App{
		Name:  "collect",
		Usage: "data collection",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "port",
				Usage:       "port used by the server.",
				Value:       8080,
				Destination: &port,
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
				Value:       "./phpstats.yml",
			},
		},
		Action: func(c *cli.Context) error {
			if len(os.Args) == 1 {
				commands.About().Execute(&shell.Context{})
				fmt.Printf("\nUsage\n\t$ phpstats collect [--config-path <dir>] [--disable-cache] [--port <value>] [--project-path <dir>] [--cache-dir <dir>] <analyze-dir>\n\n")
				return fmt.Errorf("empty")
			}

			cfg, errOpen, errDecode := config.OpenConfig(configPath)

			switch {
			case errDecode != nil:
				color.Red.Printf("Config error: %v", errDecode)
				return errDecode
			case errOpen == nil && cfg != nil && cfg.CacheDir == "":
				cfg.CacheDir = utils.DefaultCacheDir()
			case errOpen == nil && cfg.CacheDir != "":
			default:
				cfg = &config.Config{
					Port:         port,
					CacheDir:     cacheDir,
					DisableCache: disableCache,
					ProjectPath:  walkers.GlobalCtx.ProjectRoot,
					Exclude:      nil,
					Packages:     nil,
					Extensions:   nil,
				}
			}

			cfg.AddPackagesToContext(walkers.GlobalCtx.Packages)
			server.RunServer(port)

			// Normalize flags for NoVerify
			exe := os.Args[0]

			countArgs := c.NArg()
			var analyzeDirs []string
			if countArgs > 0 {
				analyzeDirs = c.Args().Slice()
			}

			cfgCli := cfg.ToCliArgs()
			os.Args = []string{exe}
			os.Args = append(os.Args, cfgCli...)
			os.Args = append(os.Args, analyzeDirs...)

			if c.NArg() > 1 {
				log.Fatalf(color.Red.Sprintf("Error: too many arguments"))
			}

			if cfg.Exclude != nil {
				excludeRegexp, err := regexp.Compile(strings.Join(cfg.Exclude, "|"))
				if err != nil {
					log.Fatalf(color.Red.Sprintf("Error converting exclude to regexp: %v", err))
					return err
				}
				walkers.GlobalCtx.ExcludeRegexp = excludeRegexp
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
