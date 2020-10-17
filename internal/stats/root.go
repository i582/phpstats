package stats

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/solver"

	"phpstats/internal/utils"
)

type rootChecker struct {
	linter.RootCheckerDefaults

	ctx *linter.RootContext

	CurFile *File
}

func (r *rootChecker) BeforeEnterFile() {
	filename := r.ctx.Filename()
	r.CurFile = GlobalCtx.Files.GetOrCreate(filename)
}

func (r *rootChecker) AfterEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.ImportExpr:
		filename, ok := utils.ResolveRequirePath(r.ctx.ClassParseState(), ProjectRoot, n.Expr)
		if !ok {
			return
		}

		requiredFile, ok := GlobalCtx.Files.Get(filename)
		if !ok {
			curFileName := r.ctx.ClassParseState().CurrentFile
			dir, _ := filepath.Split(curFileName)
			if !strings.HasSuffix(dir, `www\`) {
				if _, err := os.Stat(filename); err == nil {
					// log.Printf("%s not found", filename)
				}
			}
			return
		}

		r.CurFile.AddRequiredRootFile(requiredFile)
		requiredFile.AddRequiredByFile(r.CurFile)

	case *ir.ClassStmt:
		className, ok := solver.GetClassName(r.ctx.ClassParseState(), &ir.Name{
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
			className, ok := solver.GetClassName(r.ctx.ClassParseState(), &ir.Name{
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

	case *ir.InterfaceStmt:
		ifaceName, ok := solver.GetClassName(r.ctx.ClassParseState(), &ir.Name{
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

	case *ir.ClassConstListStmt:
		curClass, ok := r.GetCurrentClass()
		if !ok {
			return
		}

		for _, c := range n.Consts {
			constant, ok := GlobalCtx.Constants.Get(*NewConstant(c.(*ir.ConstantStmt).ConstantName.Value, curClass.Name))
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

			curClass.Fields.Add(NewField(prop.Variable.Name, curClass.Name))
		}
	}
}

func (r *rootChecker) GetCurrentFunc() (*Function, bool) {
	if r.ctx.ClassParseState().CurrentFunction == "" {
		return nil, false
	}

	class, ok := r.GetCurrentClass()
	if ok {
		method, ok := class.Methods.Get(NewMethodKey(r.ctx.ClassParseState().CurrentFunction, class.Name))
		if !ok {
			return nil, false
		}

		return method, true
	}

	funcName, ok := solver.GetFuncName(r.ctx.ClassParseState(), &ir.Name{
		Value: r.ctx.ClassParseState().CurrentFunction,
	})
	if !ok {
		return nil, false
	}

	fn, ok := GlobalCtx.Funcs.Get(NewFuncKey(funcName))
	if !ok {
		return nil, false
	}

	return fn, true
}

func (r *rootChecker) GetCurrentClass() (*Class, bool) {
	if r.ctx.ClassParseState().CurrentClass == "" {
		return nil, false
	}

	className, ok := solver.GetClassName(r.ctx.ClassParseState(), &ir.Name{
		Value: r.ctx.ClassParseState().CurrentClass,
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
