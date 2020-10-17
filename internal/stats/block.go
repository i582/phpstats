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

type blockChecker struct {
	linter.BlockCheckerDefaults
	ctx  *linter.BlockContext
	root *rootChecker
}

func (b *blockChecker) BeforeEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.FunctionStmt:
		funcName := n.FunctionName.Value

		funcInfo, ok := meta.Info.GetFunction(`\` + funcName)
		if !ok {
			return
		}

		fun, ok := GlobalCtx.Funcs.Get(NewFuncKey(funcInfo.Name))
		if !ok {
			return
		}

		// добавляем текущую функцию в текущий файл
		b.root.CurFile.AddFunc(fun)
	}
}

func (b *blockChecker) AfterEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.FunctionCallExpr:
		funcName, ok := solver.GetFuncName(b.ctx.ClassParseState(), n.Function)
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
		classType := solver.ExprType(b.ctx.Scope(), b.ctx.ClassParseState(), n.Variable)

		b.handleMethod(methodName, classType)

	case *ir.StaticCallExpr:
		method, ok := n.Call.(*ir.Identifier)
		if !ok {
			return
		}
		methodName := method.Value
		className, ok := solver.GetClassName(b.ctx.ClassParseState(), n.Class)
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
		filename, ok := utils.ResolveRequirePath(b.ctx.ClassParseState(), ProjectRoot, n.Expr)
		if !ok {
			return
		}

		requiredFile, ok := GlobalCtx.Files.Get(filename)
		if !ok {
			curFileName := b.root.ctx.ClassParseState().CurrentFile
			dir, _ := filepath.Split(curFileName)
			if !strings.HasSuffix(dir, `www\`) {
				if _, err := os.Stat(filename); err == nil {
					// log.Printf("%s not found", filename)
				}
			}
			return
		}

		b.root.CurFile.AddRequiredFile(requiredFile)
		requiredFile.AddRequiredByFile(b.root.CurFile)

	case *ir.NewExpr:
		curClass, ok := b.root.GetCurrentClass()
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
		curClass, ok := b.root.GetCurrentClass()
		if !ok {
			return
		}

		classNameNode, ok := n.Class.(*ir.Name)
		if !ok {
			return
		}

		constClassName := classNameNode.Value

		constClassName, ok = solver.GetClassName(b.root.ctx.ClassParseState(), classNameNode)
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
		curClass, ok := b.root.GetCurrentClass()
		if !ok {
			return
		}
		curMethod, ok := b.root.GetCurrentFunc()
		if !ok {
			return
		}

		name := n.Name
		curClass.Fields.AddMethodAccess(NewFieldKey(name, curClass.Name), curMethod)

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

	// вызываемая функция
	calledName := calledMethodInfo.Info.Name
	calledFuncKey := FuncKey{
		Name:      calledName,
		ClassName: calledMethodInfo.ImplName(),
	}
	// позиция вызываемой функции
	calledFunPos := calledMethodInfo.Info.Pos

	calledClass, ok := GlobalCtx.Classes.Get(calledMethodInfo.ImplName())
	if !ok {
		return
	}

	calledFunc := GlobalCtx.Funcs.GetOrCreateMethod(calledFuncKey, calledFunPos, calledClass)

	b.handleCalled(calledFunc)
}

func (b *blockChecker) handleFunc(name string) {
	calledFuncInfo, ok := meta.Info.GetFunction(name)
	if !ok {
		return
	}

	// вызываемая функция
	calledName := name
	calledFuncKey := FuncKey{
		Name: calledName,
	}
	// позиция вызываемой функции
	calledFunPos := calledFuncInfo.Pos
	calledFunc := GlobalCtx.Funcs.GetOrCreateFunction(calledFuncKey, calledFunPos)

	b.handleCalled(calledFunc)
}

func (b *blockChecker) handleCalled(calledFunc *Function) {
	curFunc, ok := b.root.GetCurrentFunc()
	if !ok {
		return
	}

	if curFunc != nil {
		// добавляем, что текущая функция вызывает функцию
		curFunc.AddCalled(calledFunc)
		// выставляем, что вызываемая функция вызывается из текущей
		calledFunc.AddCalledBy(curFunc)
	}

	curClass, ok := b.root.GetCurrentClass()
	if ok && calledFunc.Class != nil {
		curClass.AddDeps(calledFunc.Class)
		calledFunc.Class.AddDepsBy(curClass)
	}

	calledFunc.AddUse()
}
