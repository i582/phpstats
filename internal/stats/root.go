package stats

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/VKCOM/noverify/src/solver"

	"phpstats/internal/utils"
)

type rootChecker struct {
	linter.RootCheckerDefaults

	ctx *linter.RootContext

	CurFile   *File
	CurClass  *Class
	CurMethod *Function
	CurFunc   *Function
}

func (r *rootChecker) BeforeEnterFile() {
	filename := r.ctx.Filename()
	r.CurFile = GlobalCtx.Files.GetOrCreate(filename)
}

func (r *rootChecker) AfterEnterNode(n ir.Node) {
	if !meta.IsIndexingComplete() {
		return
	}

	switch n := n.(type) {
	case *ir.ImportExpr:
		filename, ok := utils.ResolveRequirePath(r.ctx.ClassParseState(), "C:\\projects\\vkcom", n.Expr)
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

		r.CurClass = class
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

		r.CurClass = iface
		r.CurFile.AddClass(iface)

	case *ir.ClassMethodStmt:
		if r.CurClass == nil {
			return
		}

		methodName := n.MethodName.Value

		method, ok := GlobalCtx.Funcs.Get(NewMethodKey(methodName, r.CurClass.Name))
		if !ok {
			return
		}

		r.CurMethod = method

	case *ir.ClassConstFetchExpr:
		classNameNode, ok := n.Class.(*ir.Name)
		if !ok {
			return
		}

		constClassName := classNameNode.Value

		constClassName, ok = solver.GetClassName(r.ctx.ClassParseState(), classNameNode)
		if !ok {
			return
		}

		class, ok := GlobalCtx.Classes.Get(constClassName)
		if !ok {
			return
		}

		if r.CurClass != nil {
			r.CurClass.AddDeps(class)
			class.AddDepsBy(r.CurClass)
		}

	case *ir.ClassConstListStmt:
		for _, c := range n.Consts {
			constant, ok := GlobalCtx.Constants.Get(*NewConstant(c.(*ir.ConstantStmt).ConstantName.Value, r.CurClass.Name))
			if !ok {
				continue
			}

			r.CurClass.Constants.Add(constant)
		}
	case *ir.PropertyListStmt:
		if r.CurClass == nil {
			return
		}

		for _, prop := range n.Properties {
			prop := prop.(*ir.PropertyStmt)

			r.CurClass.Fields.Add(NewField(prop.Variable.Name, r.CurClass.Name))
		}
	}
}

func (r *rootChecker) BeforeLeaveNode(n ir.Node) {
	switch n.(type) {
	case *ir.ClassStmt:
		r.CurClass = nil
	case *ir.ClassMethodStmt:
		r.CurMethod = nil
	}
}
