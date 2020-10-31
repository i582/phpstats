package walkers

import (
	"bytes"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/solver"

	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

type RootChecker struct {
	linter.RootCheckerDefaults

	Ctx *linter.RootContext

	CurFile *symbols.File
}

func (r *RootChecker) BeforeEnterFile() {
	filename := r.Ctx.Filename()

	var ok bool
	r.CurFile, ok = GlobalCtx.Files.Get(filename)
	if !ok {
		return
	}
	// hack, yet
	r.CurFile.CountLines = int64(bytes.Count(r.Ctx.FileContents(), []byte("\n")) + 1)

	BarLinting.Increment()
}

func (r *RootChecker) AfterEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.NamespaceStmt:
		nsName := n.NamespaceName.Value

		GlobalCtx.Namespaces.CreateNamespace(nsName)
		GlobalCtx.Namespaces.AddFileToNamespace(nsName, r.CurFile)

	case *ir.ImportExpr:
		filename, ok := utils.ResolveRequirePath(r.Ctx.ClassParseState(), ProjectRoot, n.Expr)
		if !ok {
			return
		}

		requiredFile, ok := GlobalCtx.Files.Get(filename)
		if !ok {
			return
		}

		r.CurFile.AddRequiredRootFile(requiredFile)
		requiredFile.AddRequiredByFile(r.CurFile)

	case *ir.ClassStmt:
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

	case *ir.InterfaceStmt:
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

	case *ir.ClassConstListStmt:
		curClass, ok := r.GetCurrentClass()
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

	case *ir.PropertyListStmt:
		curClass, ok := r.GetCurrentClass()
		if !ok {
			return
		}

		for _, prop := range n.Properties {
			prop := prop.(*ir.PropertyStmt)

			curClass.Fields.Add(symbols.NewField(prop.Variable.Name, curClass.Name))
		}
	}
}

func (r *RootChecker) GetCurrentFunc() (*symbols.Function, bool) {
	if r.Ctx.ClassParseState().CurrentFunction == "" {
		return nil, false
	}

	class, ok := r.GetCurrentClass()
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

	fn, ok := GlobalCtx.Funcs.Get(symbols.NewFuncKey(funcName))
	if !ok {
		return nil, false
	}

	return fn, true
}

func (r *RootChecker) GetCurrentClass() (*symbols.Class, bool) {
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
