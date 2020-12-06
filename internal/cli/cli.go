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

	MainShell := shell.NewShell()

	MainShell.AddExecutor(commands.Info())
	MainShell.AddExecutor(commands.List())
	MainShell.AddExecutor(commands.Graph())
	MainShell.AddExecutor(commands.Brief())
	MainShell.AddExecutor(commands.About())
	MainShell.AddExecutor(commands.Metrics())
	MainShell.AddExecutor(commands.Relation())

	var cacheDir string
	var configPath string
	var disableCache bool
	var port int64

	app := &cli.App{
		Name:        "phpstats",
		Version:     "v0.4.0",
		Usage:       "Tool for collecting statistics and building dependency graphs for PHP",
		Description: "Tool for collecting statistics and building dependency graphs for PHP",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Launches the configuration wizard to create a configuration file",
				Action: func(c *cli.Context) error {
					config.ConfigureConfig()
					os.Exit(0)
					return nil
				},
			},
			{
				Name:  "collect",
				Usage: "Starts collecting information and starts an interactive shell",
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
					cfg, errOpen, errDecode := config.OpenConfig(configPath)
					if errDecode != nil {
						return fmt.Errorf("config: %v", errDecode)
					}
					if errOpen != nil {
						color.Yellow.Printf("Warning: config file '%s' not open (the default configuration is used)\n", configPath)

						cfg = &config.Config{
							ProjectName:  "Untitled",
							Port:         port,
							CacheDir:     cacheDir,
							DisableCache: disableCache,
							ProjectPath:  walkers.GlobalCtx.ProjectRoot,
							Exclude:      nil,
							Packages:     nil,
							Extensions:   nil,
						}
					}

					if cfg.CacheDir == "" {
						cfg.CacheDir = utils.DefaultCacheDir()
					}

					walkers.GlobalCtx.ProjectName = cfg.ProjectName
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
						return fmt.Errorf("too many arguments")
					}

					if cfg.Exclude != nil {
						excludeRegexp, err := regexp.Compile(strings.Join(cfg.Exclude, "|"))
						if err != nil {
							return fmt.Errorf("converting exclude to regexp: %v", err)
						}
						walkers.GlobalCtx.ExcludeRegexp = excludeRegexp
					}

					err := walkers.Collect()
					if err != nil {
						return fmt.Errorf("collect: %v", err)
					}

					MainShell.Run()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(color.Red.Sprintf("Error: %v", err))
	}
}
