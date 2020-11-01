package walkers

import (
	"log"
	"os"

	"github.com/VKCOM/noverify/src/cmd"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/cheggaaa/pb/v3"

	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/filemeta"
	"github.com/i582/phpstats/internal/utils"
)

var BarLinting *pb.ProgressBar

func CollectMain() error {
	linter.RegisterBlockChecker(func(ctx *linter.BlockContext) linter.BlockChecker {
		if meta.IsIndexingComplete() {
			return &blockChecker{
				Ctx:  ctx,
				Root: ctx.RootState()["vklints-root"].(*rootChecker),
			}
		}

		return &blockIndexer{}
	})

	linter.RegisterRootCheckerWithCacher(GlobalCtx, func(ctx *linter.RootContext) linter.RootChecker {
		if meta.IsIndexingComplete() {
			checker := &rootChecker{
				Ctx: ctx,
			}
			ctx.State()["vklints-root"] = checker
			return checker
		}

		indexer := &rootIndexer{
			Ctx:  ctx,
			Meta: filemeta.NewFileMeta(),
		}
		ctx.State()["vklints-root"] = indexer
		return indexer
	})

	meta.OnIndexingComplete(func() {
		BarLinting = pb.StartNew(GlobalCtx.Files.Len())
	})

	fs, args := flags.ParseFlags(os.Args, flags.NewFlags(&flags.Flag{
		Name:      "--project-path",
		WithValue: true,
	}, &flags.Flag{
		Name:      "--cache-dir",
		WithValue: true,
	}, &flags.Flag{
		Name: "--server",
	}))

	os.Args = args

	if len(os.Args) < 2 {
		log.Fatalf("Error: too few arguments")
	}

	var cacheDir string
	if f, ok := fs.Get("--cache-dir"); ok {
		cacheDir = f.Value
	} else {
		cacheDir = utils.DefaultCacheDir()
	}

	argstmp := []string{os.Args[0]}
	argstmp = append(argstmp, "-cache-dir", cacheDir)
	argstmp = append(argstmp, os.Args[1:]...)
	os.Args = argstmp

	if flag, ok := fs.Get("--project-path"); ok {
		ProjectRoot = flag.Value
	} else if len(os.Args) > 0 {
		ProjectRoot = os.Args[len(os.Args)-1]
	}

	if _, err := os.Stat(ProjectRoot); os.IsNotExist(err) {
		log.Fatalf("Error: invalid project path: %v", err)
	}

	_, _ = cmd.Run(&cmd.MainConfig{
		BeforeReport: func(*linter.Report) bool {
			return false
		},
	})

	BarLinting.Finish()
	return nil
}
