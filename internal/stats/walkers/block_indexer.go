package walkers

import (
	"github.com/VKCOM/noverify/src/linter"
)

type blockIndexer struct {
	linter.BlockCheckerDefaults

	Ctx  *linter.BlockContext
	Root *rootIndexer
}
