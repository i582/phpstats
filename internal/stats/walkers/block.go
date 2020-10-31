package walkers

import (
	"github.com/VKCOM/noverify/src/ir"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/VKCOM/noverify/src/solver"

	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

type blockChecker struct {
	linter.BlockCheckerDefaults
	Ctx  *linter.BlockContext
	Root *rootChecker
}

func (b *blockChecker) BeforeEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.FunctionStmt:
		funcName := n.FunctionName.Value

		funcInfo, ok := meta.Info.GetFunction(`\` + funcName)
		if !ok {
			return
		}

		fun, ok := GlobalCtx.Funcs.Get(symbols.NewFuncKey(funcInfo.Name))
		if !ok {
			return
		}

		b.Root.CurFile.AddFunc(fun)
	}
}

func (b *blockChecker) AfterEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.FunctionCallExpr:
		funcName, ok := solver.GetFuncName(b.Ctx.ClassParseState(), n.Function)
		if !ok {
			return
		}

		b.handleFunc(funcName)

	case *ir.MethodCallExpr:
		method, ok := n.Method.(*ir.Identifier)
		if !ok {
			return
		}
		methodName := method.Value
		classType := solver.ExprType(b.Ctx.Scope(), b.Ctx.ClassParseState(), n.Variable)

		b.handleMethod(methodName, classType)

	case *ir.StaticCallExpr:
		method, ok := n.Call.(*ir.Identifier)
		if !ok {
			return
		}
		methodName := method.Value
		className, ok := solver.GetClassName(b.Ctx.ClassParseState(), n.Class)
		if !ok {
			return
		}
		_, ok = meta.Info.GetClassOrTrait(className)
		if !ok {
			return
		}

		classType := meta.NewTypesMap(className)

		b.handleMethod(methodName, classType)

	case *ir.ImportExpr:
		filename, ok := utils.ResolveRequirePath(b.Ctx.ClassParseState(), ProjectRoot, n.Expr)
		if !ok {
			return
		}

		requiredFile, ok := GlobalCtx.Files.Get(filename)
		if !ok {
			return
		}

		b.Root.CurFile.AddRequiredFile(requiredFile)
		requiredFile.AddRequiredByFile(b.Root.CurFile)

	case *ir.NewExpr:
		curClass, ok := b.Root.GetCurrentClass()
		if !ok {
			return
		}

		classNameNode, ok := n.Class.(*ir.Name)
		if !ok {
			return
		}
		className := classNameNode.Value

		class, ok := GlobalCtx.Classes.Get(className)
		if !ok {
			return
		}

		class.AddDepsBy(curClass)
		curClass.AddDeps(class)

	case *ir.ClassConstFetchExpr:
		classNameNode, ok := n.Class.(*ir.Name)
		if !ok {
			return
		}

		curClass, ok := b.Root.GetCurrentClass()
		if !ok {
			return
		}

		constClassName, ok := solver.GetClassName(b.Root.Ctx.ClassParseState(), classNameNode)
		if !ok {
			return
		}

		class, ok := GlobalCtx.Classes.Get(constClassName)
		if !ok {
			return
		}

		curClass.AddDeps(class)
		class.AddDepsBy(curClass)

	case *ir.SimpleVar:
		curClass, ok := b.Root.GetCurrentClass()
		if !ok {
			return
		}
		curMethod, ok := b.Root.GetCurrentFunc()
		if !ok {
			return
		}

		name := n.Name
		curClass.Fields.AddMethodAccess(symbols.NewFieldKey(name, curClass.Name), curMethod)

	default:
		return
	}
}

func (b *blockChecker) handleMethod(name string, classType meta.TypesMap) {
	var calledMethodInfo solver.FindMethodResult

	found := classType.Find(func(typ string) bool {
		var ok bool
		calledMethodInfo, ok = solver.FindMethod(typ, name)
		if !ok {
			return false
		}
		return true
	})

	if !found {
		return
	}

	calledName := calledMethodInfo.Info.Name
	calledFuncKey := symbols.FuncKey{
		Name:      calledName,
		ClassName: calledMethodInfo.ImplName(),
	}
	calledFunPos := calledMethodInfo.Info.Pos

	calledClass, ok := GlobalCtx.Classes.Get(calledMethodInfo.ImplName())
	if !ok {
		return
	}

	calledFunc, found := GlobalCtx.Funcs.Get(calledFuncKey)
	if !found {
		calledFunc = symbols.NewMethod(calledFuncKey, calledFunPos, calledClass)
	}

	b.handleCalled(calledFunc)
}

func (b *blockChecker) handleFunc(name string) {
	calledFuncInfo, ok := meta.Info.GetFunction(name)
	if !ok {
		return
	}

	calledName := name
	calledFuncKey := symbols.FuncKey{
		Name: calledName,
	}
	calledFunPos := calledFuncInfo.Pos

	calledFunc, found := GlobalCtx.Funcs.Get(calledFuncKey)
	if !found {
		calledFunc = symbols.NewFunction(calledFuncKey, calledFunPos)
	}

	b.handleCalled(calledFunc)
}

func (b *blockChecker) handleCalled(calledFunc *symbols.Function) {
	curFunc, ok := b.Root.GetCurrentFunc()
	if !ok {
		return
	}

	curFunc.AddCalled(calledFunc)
	calledFunc.AddCalledBy(curFunc)
}
