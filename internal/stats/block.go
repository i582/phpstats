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
			b.root.CurFunc = nil
			return
		}

		fun, ok := GlobalCtx.Funcs.Get(NewFuncKey(funcInfo.Name))
		if !ok {
			return
		}
		b.root.CurFunc = fun

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
		filename, ok := utils.ResolveRequirePath(b.ctx.ClassParseState(), "C:\\projects\\vkcom", n.Expr)
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
		classNameNode, ok := n.Class.(*ir.Name)
		if !ok {
			return
		}
		className := classNameNode.Value

		if b.ctx.ClassParseState().CurrentClass == "" {
			return
		}

		curClassName := b.ctx.ClassParseState().CurrentClass

		class, ok := GlobalCtx.Classes.Get(className)
		if !ok {
			return
		}

		curClass, ok := GlobalCtx.Classes.Get(curClassName)
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

		constClassName := classNameNode.Value

		constClassName, ok = solver.GetClassName(b.root.ctx.ClassParseState(), classNameNode)
		if !ok {
			return
		}

		class, ok := GlobalCtx.Classes.Get(constClassName)
		if !ok {
			return
		}

		if b.root.CurClass != nil {
			b.root.CurClass.AddDeps(class)
			class.AddDepsBy(b.root.CurClass)
		}

	case *ir.SimpleVar:
		if b.root.CurClass == nil {
			return
		}
		if b.root.CurMethod == nil {
			return
		}

		name := n.Name
		b.root.CurClass.Fields.AddMethodAccess(NewFieldKey(name, b.root.CurClass.Name), b.root.CurMethod)

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
	var curFunc *Function

	if b.root.CurFunc != nil {
		curFunc = b.root.CurFunc
	} else if b.root.CurMethod != nil {
		curFunc = b.root.CurMethod
	}

	if curFunc != nil {
		// добавляем, что текущая функция вызывает функцию
		curFunc.AddCalled(calledFunc)
		// выставляем, что вызываемая функция вызывается из текущей
		calledFunc.AddCalledBy(curFunc)
	}

	if b.root.CurClass != nil && calledFunc.Class != nil {
		b.root.CurClass.AddDeps(calledFunc.Class)
		calledFunc.Class.AddDepsBy(b.root.CurClass)
	}

	calledFunc.AddUse()
}
