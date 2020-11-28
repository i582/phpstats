package walkers

import (
	"bytes"
	"path/filepath"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
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
	case *ir.FunctionStmt:
		r.handleFunction(n)
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

func (r *rootChecker) handleFunction(n *ir.FunctionStmt) {
	funcName, ok := solver.GetFuncName(r.Ctx.ClassParseState(), &ir.Name{
		Value: n.FunctionName.Value,
	})
	if !ok {
		return
	}

	funcInfo, ok := meta.Info.GetFunction(funcName)
	if !ok {
		return
	}

	fun, ok := GlobalCtx.Functions.Get(symbols.NewFuncKey(funcInfo.Name))
	if !ok {
		return
	}

	r.CurFile.AddFunc(fun)
	GlobalCtx.Namespaces.AddFunctionToNamespace(r.Ctx.ClassParseState().Namespace, fun)
}

func (r *rootChecker) handlePropertyList(n *ir.PropertyListStmt) bool {
	curClass, ok := r.getCurrentClass()
	if !ok {
		return true
	}

	for _, prop := range n.Properties {
		prop := prop.(*ir.PropertyStmt)

		curClass.Fields.Add(symbols.NewField(prop.Variable.Name, curClass))
	}
	return false
}

func (r *rootChecker) handleClassConstList(n *ir.ClassConstListStmt) {
	curClass, ok := r.getCurrentClass()
	if !ok {
		return
	}

	for _, c := range n.Consts {
		constant := symbols.NewConstant(c.(*ir.ConstantStmt).ConstantName.Value, curClass)

		if _, found := GlobalCtx.Constants.Get(*constant); !found {
			GlobalCtx.Constants.Add(constant)
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

	for _, stmt := range n.Stmts {
		r.handleClassInterfaceMethodsConstants(iface, stmt)
	}
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

			ifaceName, ok := solver.GetClassName(r.Ctx.ClassParseState(), &ir.Name{
				Value: ifaceName,
			})

			iface, ok := GlobalCtx.Classes.Get(ifaceName)
			if !ok {
				continue
			}

			class.AddImplements(iface)
			class.AddDeps(iface)
			iface.AddDepsBy(class)
		}
	}

	if n.Extends != nil {
		className, ok := solver.GetClassName(r.Ctx.ClassParseState(), &ir.Name{
			Value: n.Extends.ClassName.Value,
		})
		if ok {
			extend, ok := GlobalCtx.Classes.Get(className)
			if ok {
				class.AddExtends(extend)
				class.AddDeps(extend)
				extend.AddDepsBy(class)
			}
		}
	}

	GlobalCtx.Namespaces.AddClassToNamespace(r.Ctx.ClassParseState().Namespace, class)

	for _, stmt := range n.Stmts {
		r.handleClassInterfaceMethodsConstants(class, stmt)
	}
}

func (r *rootChecker) handleClassInterfaceMethodsConstants(class *symbols.Class, n ir.Node) {
	switch n := n.(type) {
	case *ir.ClassMethodStmt:
		methodName := n.MethodName.Value

		method, found := GlobalCtx.Functions.Get(symbols.NewMethodKey(methodName, class.Name))
		if !found {
			return
		}

		class.AddMethod(method)

	case *ir.ClassConstListStmt:
		for _, c := range n.Consts {
			constantStmt := c.(*ir.ConstantStmt)
			constantName := constantStmt.ConstantName.Value

			constant, found := GlobalCtx.Constants.Get(symbols.NewConstantKey(constantName, class))
			if !found {
				return
			}

			class.Constants.Add(constant)

			b := &blockChecker{
				BlockCheckerDefaults: linter.BlockCheckerDefaults{},
				Ctx:                  &linter.BlockContext{},
				Root:                 r,
			}

			constantStmt.Expr.Walk(b)
		}
	}
}

func (r *rootChecker) handleImport(n *ir.ImportExpr) {
	curFileDir := filepath.Dir(r.CurFile.Path)
	filename, ok := utils.ResolveRequirePath(r.Ctx.ClassParseState(), curFileDir, n.Expr)
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
	if n.NamespaceName == nil {
		return
	}

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
