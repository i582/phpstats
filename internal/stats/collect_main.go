package stats

import (
	"github.com/VKCOM/noverify/src/cmd"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
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

	_, _ = cmd.Run(&cmd.MainConfig{
		BeforeReport: func(*linter.Report) bool {
			return false
		},
	})

	return nil
}
