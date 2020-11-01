package walkers

import (
	"log"
	"os"

	"github.com/VKCOM/noverify/src/cmd"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/cheggaaa/pb/v3"

	"github.com/i582/phpstats/internal/stats/filemeta"
)

func Collect() error {
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

	if GlobalCtx.ProjectRoot == "" {
		GlobalCtx.ProjectRoot = os.Args[len(os.Args)-1]
	}

	if _, err := os.Stat(GlobalCtx.ProjectRoot); os.IsNotExist(err) {
		log.Fatalf("Error: invalid project path: %v", err)
	}

	meta.OnIndexingComplete(func() {
		GlobalCtx.BarLinting = pb.StartNew(GlobalCtx.Files.Len())
	})

	_, _ = cmd.Run(&cmd.MainConfig{
		BeforeReport: func(*linter.Report) bool {
			return false
		},
	})

	GlobalCtx.BarLinting.Finish()
	return nil
}
