package stats

import (
	"log"
	"strings"

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

func (r *rootIndexer) inVendor() bool {
	curFileName := r.ctx.Filename()
	return strings.Contains(curFileName, "vendor") || strings.Contains(curFileName, "phpstorm-stubs")
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
		class.Vendor = r.inVendor()

		r.meta.Classes.Add(class)

		for _, n := range n.Stmts {
			r.handleClassInterfaceMethodsConstants(class, n)
		}

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

		if ifaceName == "\\Exception" {
			log.Print()
		}

		iface := NewInterface(ifaceName, curFile)
		iface.Vendor = r.inVendor()
		r.meta.Classes.Add(iface)

		for _, n := range n.Stmts {
			r.handleClassInterfaceMethodsConstants(iface, n)
		}

	case *ir.FunctionStmt:
		funcName := n.FunctionName.Value
		pos := r.getElementPos(n)

		fn := NewFunctionInfo(NewFuncKey(funcName), pos)
		r.meta.Funcs.Add(fn)
	}
}

func (r *rootIndexer) handleClassInterfaceMethodsConstants(class *Class, n ir.Node) {
	switch n := n.(type) {
	case *ir.ClassMethodStmt:
		methodName := n.MethodName.Value
		pos := r.getElementPos(n)

		fn := NewFunctionInfo(NewMethodKey(methodName, class.Name), pos)
		r.meta.Funcs.Add(fn)

	case *ir.ClassConstListStmt:
		for _, c := range n.Consts {
			r.meta.Constants.Add(NewConstant(c.(*ir.ConstantStmt).ConstantName.Value, class.Name))
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
