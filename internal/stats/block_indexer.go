package stats

import (
	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
)

type blockIndexer struct {
	linter.BlockCheckerDefaults

	ctx  *linter.BlockContext
	root *rootIndexer
}

func (b *blockIndexer) BeforeEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.FunctionStmt:
		funcName := n.FunctionName.Value
		pos := b.root.getElementPos(n)

		fn := NewFunctionInfo(NewFuncKey(funcName), pos)
		b.root.meta.Funcs.Add(fn)
	}
}
