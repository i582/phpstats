package walkers

import (
	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
)

type blockIndexer struct {
	linter.BlockCheckerDefaults
	Ctx  *linter.BlockContext
	Root *rootIndexer
}

// AfterEnterNode describes the processing logic after entering the node.
func (b *blockIndexer) AfterEnterNode(n ir.Node) {
	switch n.(type) {
	case *ir.ClosureExpr:
		b.Root.Meta.CountAnonymousFunctions++
	}
}
