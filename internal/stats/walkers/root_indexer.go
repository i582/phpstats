package walkers

import (
	"log"
	"strings"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/ir/irutil"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/VKCOM/noverify/src/php/parser/freefloating"
	"github.com/VKCOM/noverify/src/solver"

	"github.com/i582/phpstats/internal/stats/filemeta"
	"github.com/i582/phpstats/internal/stats/symbols"
)

type rootIndexer struct {
	linter.RootCheckerDefaults

	Ctx  *linter.RootContext
	Meta filemeta.FileMeta
}

// BeforeEnterFile describes the processing logic before entering the file.
func (r *rootIndexer) BeforeEnterFile() {
	curFileName := r.Ctx.Filename()
	curFile := symbols.NewFile(curFileName)

	r.Meta.Files.Add(curFile)
}

// AfterLeaveFile describes the processing logic after leaving the file.
func (r *rootIndexer) AfterLeaveFile() {
	GlobalCtx.UpdateMeta(&r.Meta, "")
}

// AfterEnterNode describes the processing logic after entering the node.
func (r *rootIndexer) AfterEnterNode(n ir.Node) {
	if ffs := n.GetFreeFloating(); ffs != nil {
		for _, cs := range *ffs {
			for _, c := range cs {
				r.handleComment(c)
			}
		}
	}

	switch n := n.(type) {
	case *ir.ClassStmt:
		r.handleClass(n)
	case *ir.InterfaceStmt:
		r.handleInterface(n)
	case *ir.FunctionStmt:
		r.handleFunction(n)
	}
}

func (r *rootIndexer) handleComment(c freefloating.String) {
	if c.StringType != freefloating.CommentType {
		return
	}

	lines := strings.Count(c.Value, "\n")
	r.Meta.CountCommentLine += int64(lines + 1)
}

func (r *rootIndexer) inVendor() bool {
	curFileName := r.Ctx.Filename()
	return strings.Contains(curFileName, "vendor") || strings.Contains(curFileName, "phpstorm-stubs")
}

func (r *rootIndexer) handleFunction(n *ir.FunctionStmt) {
	namespace := r.Ctx.ClassParseState().Namespace
	funcName := n.FunctionName.Value

	if namespace != "" {
		funcName = namespace + `\` + funcName
	} else {
		funcName = `\` + funcName
	}

	pos := r.getElementPos(n)

	cc := r.calculateCyclomaticComplexity(&ir.StmtList{
		Stmts: n.Stmts,
	})
	cmn := r.calculateCountMagicNumbers(&ir.StmtList{
		Stmts: n.Stmts,
	})

	fn := symbols.NewFunction(symbols.NewFuncKey(funcName), pos)
	fn.CyclomaticComplexity = cc
	fn.CountMagicNumbers = cmn
	r.Meta.Funcs.Add(fn)
}

func (r *rootIndexer) handleInterface(n *ir.InterfaceStmt) {
	curFileName := r.Ctx.Filename()

	ifaceName, ok := solver.GetClassName(r.Ctx.ClassParseState(), &ir.Name{
		Value: n.InterfaceName.Value,
	})
	if !ok {
		return
	}

	curFile, ok := r.Meta.Files.Get(curFileName)
	if !ok {
		return
	}

	iface := symbols.NewInterface(ifaceName, curFile)
	iface.Vendor = r.inVendor()
	r.Meta.Classes.Add(iface)

	for _, n := range n.Stmts {
		r.handleClassInterfaceMethodsConstants(iface, n)
	}
}

func (r *rootIndexer) handleClass(n *ir.ClassStmt) {
	curFileName := r.Ctx.Filename()

	className, ok := solver.GetClassName(r.Ctx.ClassParseState(), &ir.Name{
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

	curFile, ok := r.Meta.Files.Get(curFileName)
	if !ok {
		log.Fatalf("file not found")
	}

	class := symbols.NewClass(className, curFile)
	class.IsAbstract = isAbstract
	class.Vendor = r.inVendor()

	r.Meta.Classes.Add(class)

	for _, n := range n.Stmts {
		r.handleClassInterfaceMethodsConstants(class, n)
	}
}

func (r *rootIndexer) handleClassInterfaceMethodsConstants(class *symbols.Class, n ir.Node) {
	switch n := n.(type) {
	case *ir.ClassMethodStmt:
		methodName := n.MethodName.Value
		pos := r.getElementPos(n)

		var cc int64
		var cmn int64
		if n, ok := n.Stmt.(*ir.StmtList); ok {
			cc = r.calculateCyclomaticComplexity(n)
			cmn = r.calculateCountMagicNumbers(n)
		}

		fn := symbols.NewFunction(symbols.NewMethodKey(methodName, class.Name), pos)
		fn.CyclomaticComplexity = cc
		fn.CountMagicNumbers = cmn
		r.Meta.Funcs.Add(fn)

	case *ir.ClassConstListStmt:
		for _, c := range n.Consts {
			r.Meta.Constants.Add(symbols.NewConstant(c.(*ir.ConstantStmt).ConstantName.Value, class.Name))
		}
	}
}

func (r *rootIndexer) calculateCyclomaticComplexity(stmts *ir.StmtList) int64 {
	var complexity int64
	irutil.Inspect(stmts, func(n ir.Node) bool {
		switch n.(type) {
		case *ir.IfStmt, *ir.ForStmt, *ir.WhileStmt, *ir.ForeachStmt,
			*ir.CaseStmt, *ir.DefaultStmt, *ir.ContinueStmt, *ir.BreakStmt,
			*ir.GotoStmt, *ir.CatchStmt, *ir.TernaryExpr, *ir.CoalesceExpr,
			*ir.BooleanOrExpr, *ir.BooleanAndExpr:
			complexity++
		}
		return true
	})
	return complexity
}

func (r *rootIndexer) calculateCountMagicNumbers(stmts *ir.StmtList) int64 {
	var count int64
	irutil.Inspect(stmts, func(n ir.Node) bool {
		switch n := n.(type) {
		case *ir.Lnumber:
			if n.Value == "0" || n.Value == "1" {
				return true
			}
			count++
		case *ir.Dnumber:
			if n.Value == "0.0" || n.Value == "1.0" {
				return true
			}
			count++
		case *ir.ArrayExpr:
			return false
		case *ir.ArrayItemExpr:
			return false
		case *ir.ModExpr:
			return false
		case *ir.ArrayDimFetchExpr:
			return false
		}
		return true
	})
	return count
}

func (r *rootIndexer) getElementPos(n ir.Node) meta.ElementPosition {
	pos := ir.GetPosition(n)

	return meta.ElementPosition{
		Filename:  r.Ctx.ClassParseState().CurrentFile,
		Character: int32(0),
		Line:      int32(pos.StartLine),
		EndLine:   int32(pos.EndLine),
		Length:    int32(pos.EndPos - pos.StartPos),
	}
}
