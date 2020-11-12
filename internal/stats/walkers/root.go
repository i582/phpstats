package walkers

import (
	"bytes"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/solver"

	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

type rootChecker struct {
	linter.RootCheckerDefaults

	Ctx *linter.RootContext

	CurFile *symbols.File
}

// BeforeEnterFile describes the processing logic before entering the file.
func (r *rootChecker) BeforeEnterFile() {
	filename := r.Ctx.Filename()

	var ok bool
	r.CurFile, ok = GlobalCtx.Files.Get(filename)
	if !ok {
		r.CurFile = symbols.NewFile("")
		return
	}
	// hack, yet
	r.CurFile.CountLines = int64(bytes.Count(r.Ctx.FileContents(), []byte("\n")) + 1)

	GlobalCtx.BarLinting.Increment()
}

// AfterEnterNode describes the processing logic after entering the node.
func (r *rootChecker) AfterEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.NamespaceStmt:
		r.handleNamespace(n)
	case *ir.ImportExpr:
		r.handleImport(n)
	case *ir.ClassStmt:
		r.handleClass(n)
	case *ir.InterfaceStmt:
		r.handleInterface(n)
	case *ir.ClassConstListStmt:
		r.handleClassConstList(n)
	case *ir.PropertyListStmt:
		r.handlePropertyList(n)
	}
}

func (r *rootChecker) handlePropertyList(n *ir.PropertyListStmt) bool {
	curClass, ok := r.getCurrentClass()
	if !ok {
		return true
	}

	for _, prop := range n.Properties {
		prop := prop.(*ir.PropertyStmt)

		curClass.Fields.Add(symbols.NewField(prop.Variable.Name, curClass.Name))
	}
	return false
}

func (r *rootChecker) handleClassConstList(n *ir.ClassConstListStmt) {
	curClass, ok := r.getCurrentClass()
	if !ok {
		return
	}

	for _, c := range n.Consts {
		constant, ok := GlobalCtx.Constants.Get(*symbols.NewConstant(c.(*ir.ConstantStmt).ConstantName.Value, curClass.Name))
		if !ok {
			continue
		}

		curClass.Constants.Add(constant)
	}
	return
}

func (r *rootChecker) handleInterface(n *ir.InterfaceStmt) {
	ifaceName, ok := solver.GetClassName(r.Ctx.ClassParseState(), &ir.Name{
		Value: n.InterfaceName.Value,
	})
	if !ok {
		return
	}

	iface, ok := GlobalCtx.Classes.Get(ifaceName)
	if !ok {
		return
	}

	r.CurFile.AddClass(iface)
	GlobalCtx.Namespaces.AddClassToNamespace(r.Ctx.ClassParseState().Namespace, iface)
}

func (r *rootChecker) handleClass(n *ir.ClassStmt) {
	className, ok := solver.GetClassName(r.Ctx.ClassParseState(), &ir.Name{
		Value: n.ClassName.Value,
	})
	if !ok {
		return
	}

	class, ok := GlobalCtx.Classes.Get(className)
	if !ok {
		return
	}

	r.CurFile.AddClass(class)

	if n.Implements != nil {
		for _, implement := range n.Implements.InterfaceNames {
			implement := implement.(*ir.Name)
			ifaceName := implement.Value

			iface, ok := GlobalCtx.Classes.Get(ifaceName)
			if !ok {
				return
			}

			class.AddDeps(iface)
			iface.AddDepsBy(class)
		}
	}

	if n.Extends != nil {
		className, ok := solver.GetClassName(r.Ctx.ClassParseState(), &ir.Name{
			Value: n.Extends.ClassName.Value,
		})
		if !ok {
			return
		}

		extend, ok := GlobalCtx.Classes.Get(className)
		if !ok {
			return
		}

		class.AddExtends(extend)
		class.AddDeps(extend)
		extend.AddDepsBy(class)
	}

	GlobalCtx.Namespaces.AddClassToNamespace(r.Ctx.ClassParseState().Namespace, class)
	return
}

func (r *rootChecker) handleImport(n *ir.ImportExpr) {
	filename, ok := utils.ResolveRequirePath(r.Ctx.ClassParseState(), GlobalCtx.ProjectRoot, n.Expr)
	if !ok {
		return
	}

	requiredFile, ok := GlobalCtx.Files.Get(filename)
	if !ok {
		return
	}

	r.CurFile.AddRequiredRootFile(requiredFile)
	requiredFile.AddRequiredByFile(r.CurFile)
}

func (r *rootChecker) handleNamespace(n *ir.NamespaceStmt) {
	nsName := n.NamespaceName.Value

	GlobalCtx.Namespaces.CreateNamespace(nsName)
	GlobalCtx.Namespaces.AddFileToNamespace(nsName, r.CurFile)
}

func (r *rootChecker) getCurrentFunc() (*symbols.Function, bool) {
	if r.Ctx.ClassParseState().CurrentFunction == "" {
		return nil, false
	}

	class, ok := r.getCurrentClass()
	if ok {
		method, ok := class.Methods.Get(symbols.NewMethodKey(r.Ctx.ClassParseState().CurrentFunction, class.Name))
		if !ok {
			return nil, false
		}

		return method, true
	}

	funcName, ok := solver.GetFuncName(r.Ctx.ClassParseState(), &ir.Name{
		Value: r.Ctx.ClassParseState().CurrentFunction,
	})
	if !ok {
		return nil, false
	}

	fn, ok := GlobalCtx.Functions.Get(symbols.NewFuncKey(funcName))
	if !ok {
		return nil, false
	}

	return fn, true
}

func (r *rootChecker) getCurrentClass() (*symbols.Class, bool) {
	if r.Ctx.ClassParseState().CurrentClass == "" {
		return nil, false
	}

	className, ok := solver.GetClassName(r.Ctx.ClassParseState(), &ir.Name{
		Value: r.Ctx.ClassParseState().CurrentClass,
	})
	if !ok {
		return nil, false
	}

	class, ok := GlobalCtx.Classes.Get(className)
	if !ok {
		return nil, false
	}

	return class, true
}
