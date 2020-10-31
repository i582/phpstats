package walkers

import (
	"github.com/VKCOM/noverify/src/linter"
)

type BlockIndexer struct {
	linter.BlockCheckerDefaults

	Ctx  *linter.BlockContext
	Root *RootIndexer
}
