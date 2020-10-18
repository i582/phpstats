package stats

import (
	"log"
	"os"

	"github.com/VKCOM/noverify/src/cmd"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"

	"github.com/i582/phpstats/internal/shell/flags"
)

func CollectMain() error {
	linter.RegisterBlockChecker(func(ctx *linter.BlockContext) linter.BlockChecker {
		if meta.IsIndexingComplete() {
			return &blockChecker{
				ctx:  ctx,
				root: ctx.RootState()["vklints-root"].(*rootChecker),
			}
		}

		return &blockIndexer{
			ctx:  ctx,
			root: ctx.RootState()["vklints-root"].(*rootIndexer),
		}
	})

	linter.RegisterRootCheckerWithCacher(GlobalCtx, func(ctx *linter.RootContext) linter.RootChecker {
		if meta.IsIndexingComplete() {
			checker := &rootChecker{
				ctx: ctx,
			}
			ctx.State()["vklints-root"] = checker
			return checker
		}

		indexer := &rootIndexer{
			ctx:  ctx,
			meta: NewFileMeta(),
		}
		ctx.State()["vklints-root"] = indexer
		return indexer
	})

	fs, args := flags.ParseFlags(os.Args, flags.NewFlags(&flags.Flag{
		Name:      "--project-path",
		WithValue: true,
	}))

	os.Args = args

	if len(os.Args) < 2 {
		log.Fatalf("Error: too few arguments")
	}

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

	return nil
}
