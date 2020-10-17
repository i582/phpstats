package stats

import (
	"log"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/VKCOM/noverify/src/solver"
)

type rootIndexer struct {
	linter.RootCheckerDefaults

	ctx  *linter.RootContext
	meta FileMeta
}

func (r *rootIndexer) BeforeEnterFile() {
	curFileName := r.ctx.Filename()
	curFile := NewFile(curFileName)

	r.meta.Files.Add(curFile)
}

func (r *rootIndexer) AfterLeaveFile() {
	GlobalCtx.updateMeta(&r.meta)
}

func (r *rootIndexer) AfterEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.ClassStmt:
		curFileName := r.ctx.Filename()

		className, ok := solver.GetClassName(r.ctx.ClassParseState(), &ir.Name{
			Value: n.ClassName.Value,
		})
		if !ok {
			return
		}

		var isAbstract bool
		for _, modifier := range n.Modifiers {
			if modifier.Value == "abstract" {
				isAbstract = true
			}
		}

		curFile, ok := r.meta.Files.Get(curFileName)
		if !ok {
			log.Fatalf("file not found")
		}

		class := NewClass(className, curFile)
		class.IsAbstract = isAbstract

		r.meta.Classes.Add(class)

	case *ir.InterfaceStmt:
		curFileName := r.ctx.Filename()

		ifaceName, ok := solver.GetClassName(r.ctx.ClassParseState(), &ir.Name{
			Value: n.InterfaceName.Value,
		})
		if !ok {
			return
		}

		curFile, ok := r.meta.Files.Get(curFileName)
		if !ok {
			log.Fatalf("file not found")
		}

		iface := NewInterface(ifaceName, curFile)
		r.meta.Classes.Add(iface)

	case *ir.ClassMethodStmt:
		currentClassName := r.ctx.ClassParseState().CurrentClass
		methodName := n.MethodName.Value
		pos := r.getElementPos(n)

		fn := NewFunctionInfo(NewMethodKey(methodName, currentClassName), pos)
		r.meta.Funcs.Add(fn)

	case *ir.ClassConstListStmt:
		for _, c := range n.Consts {
			r.meta.Constants.Add(NewConstant(c.(*ir.ConstantStmt).ConstantName.Value, r.ctx.ClassParseState().CurrentClass))
		}
	}
}

func (r *rootIndexer) getElementPos(n ir.Node) meta.ElementPosition {
	pos := ir.GetPosition(n)

	return meta.ElementPosition{
		Filename:  r.ctx.ClassParseState().CurrentFile,
		Character: int32(0),
		Line:      int32(pos.StartLine),
		EndLine:   int32(pos.EndLine),
		Length:    int32(pos.EndPos - pos.StartPos),
	}
}
