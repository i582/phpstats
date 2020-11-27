package walkers

import (
	"path/filepath"

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

func (b *blockChecker) EnterNode(n ir.Node) bool {
	b.AfterEnterNode(n)
	return true
}

func (b *blockChecker) LeaveNode(ir.Node) {}

// AfterEnterNode describes the processing logic after entering the node.
func (b *blockChecker) AfterEnterNode(n ir.Node) {
	switch n := n.(type) {
	case *ir.Argument:
		n.Expr.Walk(b)
	case *ir.FunctionStmt:
		b.handleFunction(n)
	case *ir.FunctionCallExpr:
		b.handleFunctionCall(n)
	case *ir.MethodCallExpr:
		b.handleMethodCall(n)
	case *ir.StaticCallExpr:
		b.handleStaticMethodCall(n)
	case *ir.ImportExpr:
		b.handleImport(n)
	case *ir.NewExpr:
		b.handleNew(n)
	case *ir.ClassConstFetchExpr:
		b.handleClassConstFetch(n)
	case *ir.ConstFetchExpr:
		b.handleConstFetch(n)
	case *ir.StaticPropertyFetchExpr:
		b.handleStaticPropertyFetch(n)
	case *ir.SimpleVar:
		b.handleSimpleVar(n)
	case *ir.PropertyFetchExpr:
		b.handlePropertyFetch(n)
	case *ir.Assign:
		b.handleAssign(n)
	}
}

func (b *blockChecker) handleAssign(a *ir.Assign) {
	switch n := a.Variable.(type) {
	case *ir.PropertyFetchExpr:
		b.handlePropertyFetch(n)
	case *ir.StaticPropertyFetchExpr:
		b.handleStaticPropertyFetch(n)
	}
}

func (b *blockChecker) handleStaticPropertyFetch(n *ir.StaticPropertyFetchExpr) {
	curMethod, ok := b.Root.getCurrentFunc()
	if !ok {
		return
	}

	propNameNode, ok := n.Property.(*ir.SimpleVar)
	if !ok {
		return
	}
	propName := propNameNode.Name

	className, ok := solver.GetClassName(b.Root.Ctx.ClassParseState(), n.Class)
	if !ok {
		return
	}

	p, ok := solver.FindProperty(className, "$"+propName)
	if !ok {
		return
	}

	class, ok := GlobalCtx.Classes.Get(p.ImplName())
	if !ok {
		return
	}

	class.Fields.AddMethodAccess(symbols.NewFieldKey(propName, class.Name), class, curMethod)
}

func (b *blockChecker) handlePropertyFetch(n *ir.PropertyFetchExpr) {
	curMethod, ok := b.Root.getCurrentFunc()
	if !ok {
		return
	}

	propNameNode, ok := n.Property.(*ir.Identifier)
	if !ok {
		return
	}
	propName := propNameNode.Value

	tp := solver.ExprType(b.Ctx.Scope(), b.Root.Ctx.ClassParseState(), n.Variable)
	tp.Iterate(func(typ string) {
		p, ok := solver.FindProperty(typ, propName)
		if !ok {
			return
		}

		class, ok := GlobalCtx.Classes.Get(p.ImplName())
		if !ok {
			return
		}

		class.Fields.AddMethodAccess(symbols.NewFieldKey(propName, class.Name), class, curMethod)
	})
}

func (b *blockChecker) handleSimpleVar(*ir.SimpleVar) {}

func (b *blockChecker) handleConstFetch(n *ir.ConstFetchExpr) {
	curClass, ok := b.Root.getCurrentClass()
	if !ok {
		return
	}
	curMethod, ok := b.Root.getCurrentFunc()
	if !ok {
		return
	}
	constantName := n.Constant.Value
	var constantKey symbols.Constant
	if utils.IsEmbeddedConstant(constantName) {
		constantKey = symbols.NewConstantKey(constantName, nil)
	} else if utils.IsSuperGlobal(constantName) {
		constantKey = symbols.NewConstantKey(constantName, nil)
	} else {
		constantKey = symbols.NewConstantKey(constantName, curClass)
	}

	curClass.UsedConstants.AddMethodAccess(constantKey, curMethod)
}

func (b *blockChecker) handleClassConstFetch(n *ir.ClassConstFetchExpr) {
	constClassName, ok := solver.GetClassName(b.Root.Ctx.ClassParseState(), n.Class)
	if !ok {
		return
	}

	curClass, ok := b.Root.getCurrentClass()
	if !ok {
		return
	}

	class, ok := GlobalCtx.Classes.Get(constClassName)
	if !ok {
		return
	}

	constantName := n.ConstantName.Value

	curMethod, ok := b.Root.getCurrentFunc()
	if ok {
		class.Constants.AddMethodAccess(symbols.NewConstantKey(constantName, class), curMethod)
	}

	curClass.AddDeps(class)
	class.AddDepsBy(curClass)
}

func (b *blockChecker) handleNew(n *ir.NewExpr) {
	curClass, ok := b.Root.getCurrentClass()
	if !ok {
		return
	}

	className, ok := solver.GetClassName(b.Ctx.ClassParseState(), n.Class)
	if !ok {
		return
	}

	class, ok := GlobalCtx.Classes.Get(className)
	if !ok {
		return
	}

	class.AddDepsBy(curClass)
	curClass.AddDeps(class)
}

func (b *blockChecker) handleImport(n *ir.ImportExpr) {
	curFileDir := filepath.Dir(b.Root.CurFile.Path)
	filename, ok := utils.ResolveRequirePath(b.Ctx.ClassParseState(), curFileDir, n.Expr)
	if !ok {
		return
	}

	requiredFile, ok := GlobalCtx.Files.Get(filename)
	if !ok {
		return
	}

	b.Root.CurFile.AddRequiredFile(requiredFile)
	requiredFile.AddRequiredByFile(b.Root.CurFile)
}

func (b *blockChecker) handleStaticMethodCall(n *ir.StaticCallExpr) {
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
}

func (b *blockChecker) handleMethodCall(n *ir.MethodCallExpr) {
	method, ok := n.Method.(*ir.Identifier)
	if !ok {
		return
	}
	methodName := method.Value
	classType := solver.ExprType(b.Ctx.Scope(), b.Ctx.ClassParseState(), n.Variable)

	b.handleMethod(methodName, classType)

	for _, nn := range n.Args {
		nn.Walk(b)
	}
}

func (b *blockChecker) handleFunctionCall(n *ir.FunctionCallExpr) {
	name, ok := solver.GetFuncName(b.Ctx.ClassParseState(), n.Function)
	if !ok {
		return
	}

	calledFuncInfo, ok := meta.Info.GetFunction(name)
	if !ok {
		return
	}

	calledName := name
	calledFuncKey := symbols.FuncKey{
		Name: calledName,
	}
	calledFunPos := calledFuncInfo.Pos

	calledFunc, found := GlobalCtx.Functions.Get(calledFuncKey)
	if !found {
		calledFunc = symbols.NewFunction(calledFuncKey, calledFunPos)
	}

	b.handleCalled(calledFunc)

	for _, nn := range n.Args {
		nn.Walk(b)
	}
}

func (b *blockChecker) handleFunction(n *ir.FunctionStmt) {
	funcName := n.FunctionName.Value

	funcInfo, ok := meta.Info.GetFunction(`\` + funcName)
	if !ok {
		return
	}

	fun, ok := GlobalCtx.Functions.Get(symbols.NewFuncKey(funcInfo.Name))
	if !ok {
		return
	}

	b.Root.CurFile.AddFunc(fun)
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

	calledFunc, found := GlobalCtx.Functions.Get(calledFuncKey)
	if !found {
		calledFunc = symbols.NewMethod(calledFuncKey, calledFunPos, calledClass)
	}

	b.handleCalled(calledFunc)
}

func (b *blockChecker) handleCalled(calledFunc *symbols.Function) {
	curFunc, ok := b.Root.getCurrentFunc()
	if !ok {
		return
	}

	curFunc.AddCalled(calledFunc)
	calledFunc.AddCalledBy(curFunc)
}
