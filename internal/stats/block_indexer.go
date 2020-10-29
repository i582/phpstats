package stats

import (
	"github.com/VKCOM/noverify/src/linter"
)

type blockIndexer struct {
	linter.BlockCheckerDefaults

	ctx  *linter.BlockContext
	root *rootIndexer
}
