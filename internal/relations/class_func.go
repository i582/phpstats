package relations

import (
	"fmt"

	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

type Class2FuncRelation struct {
	TargetClass     *symbols.Class
	RelatedFunction *symbols.Function

	BelongsToClass      bool
	ClassUsedInFunction bool

	FunctionUsedInClass bool
	MethodsWhereUsed    *symbols.Functions

	UsedMethods   *symbols.Functions
	UsedFields    *symbols.Fields
	UsedConstants *symbols.Constants
}

func NewClass2FuncRelation() *Class2FuncRelation {
	return &Class2FuncRelation{
		MethodsWhereUsed: symbols.NewFunctions(),
		UsedMethods:      symbols.NewFunctions(),
		UsedFields:       symbols.NewFields(),
		UsedConstants:    symbols.NewConstants(),
	}
}

func (r *Class2FuncRelation) String() string {
	var res string

	res += cfmt.Sprintf("Class {{%s}}::green connection with function {{%s}}::yellow.\n\n", r.TargetClass.Name, r.RelatedFunction.Name)

	res += cfmt.Sprintf("    Class {{%s}}::green contains function {{%s}}::yellow:   %t\n", r.TargetClass.Name, r.RelatedFunction.Name, r.BelongsToClass)

	res += cfmt.Sprintf("    Method {{%s}}::yellow uses class {{%s}}::green:         %t\n", r.RelatedFunction.Name, r.TargetClass.Name, r.ClassUsedInFunction)
	if r.ClassUsedInFunction {
		res += cfmt.Sprintf("       Uses:\n")
		for _, usedFunction := range r.UsedMethods.Funcs {
			res += cfmt.Sprintf("        method   {{%s}}::green\n", usedFunction.Name)
		}
		for _, usedField := range r.UsedFields.Fields {
			res += cfmt.Sprintf("        field    {{%s}}::green\n", usedField)
		}
		for _, usedConstant := range r.UsedConstants.Constants {
			res += cfmt.Sprintf("        constant {{%s}}::green\n", usedConstant)
		}
		res += fmt.Sprintln()
	}

	res += cfmt.Sprintf("    Class {{%s}}::green uses function {{%s}}::yellow:       %t\n", r.TargetClass.Name, r.RelatedFunction.Name, r.FunctionUsedInClass)
	if r.FunctionUsedInClass {
		res += cfmt.Sprintf("    Uses in the following methods:\n")
		for _, methodWhereUsed := range r.MethodsWhereUsed.Funcs {
			res += cfmt.Sprintf("        {{%s}}::green\n", methodWhereUsed.Name)
		}
	}

	return res
}

// relation classes --target TargetClass --related RelatedClass
func GetClass2FuncRelation(targetClass *symbols.Class, relatedFunction *symbols.Function) *Class2FuncRelation {
	rel := NewClass2FuncRelation()
	rel.TargetClass = targetClass
	rel.RelatedFunction = relatedFunction

	for _, method := range targetClass.Methods.Funcs {
		if method == relatedFunction {
			rel.BelongsToClass = true
		}

		for _, calledFunction := range method.Called.Funcs {
			if calledFunction == relatedFunction {
				rel.FunctionUsedInClass = true
				rel.MethodsWhereUsed.Add(method)
			}
		}
	}

	for _, calledFunction := range relatedFunction.Called.Funcs {
		if calledFunction.Class != targetClass {
			continue
		}

		rel.UsedMethods.Add(calledFunction)
		rel.ClassUsedInFunction = true
	}

	for _, usedField := range relatedFunction.UsedFields.Fields {
		if usedField.Class != targetClass {
			continue
		}

		rel.UsedFields.Add(usedField)
		rel.ClassUsedInFunction = true
	}

	for _, usedConstant := range relatedFunction.UsedConstants.Constants {
		if usedConstant.Class != targetClass {
			continue
		}

		rel.UsedConstants.Add(usedConstant)
		rel.ClassUsedInFunction = true
	}

	return rel
}
