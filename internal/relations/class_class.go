package relations

import (
	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

type Class2ClassRelation struct {
	TargetClass  *symbols.Class
	RelatedClass *symbols.Class

	UsedRelatedMethods   *symbols.Functions
	UsedRelatedFields    *symbols.Fields
	UsedRelatedConstants *symbols.Constants

	WhereRelatedUsedMethods   map[*symbols.Function]*symbols.Function
	WhereRelatedUsedFields    map[*symbols.Field]*symbols.Function
	WhereRelatedUsedConstants map[*symbols.Constant]*symbols.Function

	UsedTargetMethods   *symbols.Functions
	UsedTargetFields    *symbols.Fields
	UsedTargetConstants *symbols.Constants

	WhereTargetUsedMethods   map[*symbols.Function]*symbols.Function
	WhereTargetUsedFields    map[*symbols.Field]*symbols.Function
	WhereTargetUsedConstants map[*symbols.Constant]*symbols.Function

	IsTargetImplements bool
	IsTargetExtends    bool

	IsRelatedImplements bool
	IsRelatedExtends    bool
}

func NewClass2ClassRelation() *Class2ClassRelation {
	return &Class2ClassRelation{
		UsedRelatedMethods:        symbols.NewFunctions(),
		UsedRelatedFields:         symbols.NewFields(),
		UsedRelatedConstants:      symbols.NewConstants(),
		WhereRelatedUsedMethods:   map[*symbols.Function]*symbols.Function{},
		WhereRelatedUsedFields:    map[*symbols.Field]*symbols.Function{},
		WhereRelatedUsedConstants: map[*symbols.Constant]*symbols.Function{},
		UsedTargetMethods:         symbols.NewFunctions(),
		UsedTargetFields:          symbols.NewFields(),
		UsedTargetConstants:       symbols.NewConstants(),
		WhereTargetUsedMethods:    map[*symbols.Function]*symbols.Function{},
		WhereTargetUsedFields:     map[*symbols.Field]*symbols.Function{},
		WhereTargetUsedConstants:  map[*symbols.Constant]*symbols.Function{},
	}
}

func (r *Class2ClassRelation) String() string {
	var res string

	res += cfmt.Sprintf("Class {{%s}}::green connection with class {{%s}}::yellow.\n\n", r.TargetClass.Name, r.RelatedClass.Name)

	res += cfmt.Sprintf("    Class {{%s}}::green extends class {{%s}}::yellow:         %t\n", r.TargetClass.Name, r.RelatedClass.Name, r.IsTargetExtends)
	res += cfmt.Sprintf("    Class {{%s}}::green implements interface {{%s}}::yellow:  %t\n", r.TargetClass.Name, r.RelatedClass.Name, r.IsTargetImplements)

	for _, method := range r.UsedRelatedMethods.Funcs {
		res += cfmt.Sprintf("    Class {{%s}}::green uses method {{%s}}::yellow in method {{%s}}::green.\n", r.TargetClass.Name, method.Name, r.WhereRelatedUsedMethods[method].Name)
	}
	for _, field := range r.UsedRelatedFields.Fields {
		res += cfmt.Sprintf("    Class {{%s}}::green uses field {{%s}}::yellow in method {{%s}}::green.\n", r.TargetClass.Name, field, r.WhereRelatedUsedFields[field].Name)
	}
	for _, constant := range r.UsedRelatedConstants.Constants {
		res += cfmt.Sprintf("    Class {{%s}}::green uses constant {{%s}}::yellow in method {{%s}}::green.\n", r.TargetClass.Name, constant, r.WhereRelatedUsedConstants[constant].Name)
	}

	res += cfmt.Sprintln()

	res += cfmt.Sprintf("    Class {{%s}}::green extends class {{%s}}::yellow:         %t\n", r.RelatedClass.Name, r.TargetClass.Name, r.IsRelatedExtends)
	res += cfmt.Sprintf("    Class {{%s}}::green implements interface {{%s}}::yellow:  %t\n", r.RelatedClass.Name, r.TargetClass.Name, r.IsRelatedImplements)

	for _, method := range r.UsedTargetMethods.Funcs {
		res += cfmt.Sprintf("    Class {{%s}}::green uses method {{%s}}::yellow in method {{%s}}::green.\n", r.RelatedClass.Name, method.Name, r.WhereTargetUsedMethods[method].Name)
	}
	for _, field := range r.UsedTargetFields.Fields {
		res += cfmt.Sprintf("    Class {{%s}}::green uses field {{%s}}::yellow in method {{%s}}::green.\n", r.RelatedClass.Name, field, r.WhereTargetUsedFields[field].Name)
	}
	for _, constant := range r.UsedTargetConstants.Constants {
		res += cfmt.Sprintf("    Class {{%s}}::green uses constant {{%s}}::yellow in method {{%s}}::green.\n", r.RelatedClass.Name, constant, r.WhereTargetUsedConstants[constant].Name)
	}

	return res
}

func GetClass2ClassRelation(targetClass, relatedClass *symbols.Class) *Class2ClassRelation {
	if targetClass == relatedClass {
		return nil
	}
	rel := NewClass2ClassRelation()
	rel.TargetClass = targetClass
	rel.RelatedClass = relatedClass

	for _, method := range targetClass.Methods.Funcs {
		for _, calledFunction := range method.Called.Funcs {
			if calledFunction.Class != relatedClass {
				continue
			}

			rel.UsedRelatedMethods.Add(calledFunction)
			rel.WhereRelatedUsedMethods[calledFunction] = method
		}

		for _, usedField := range method.UsedFields.Fields {
			if usedField.Class != relatedClass {
				continue
			}

			rel.UsedRelatedFields.Add(usedField)
			rel.WhereRelatedUsedFields[usedField] = method
		}

		for _, usedConstant := range method.UsedConstants.Constants {
			if usedConstant.Class != relatedClass {
				continue
			}

			rel.UsedRelatedConstants.Add(usedConstant)
			rel.WhereRelatedUsedConstants[usedConstant] = method
		}
	}

	for _, method := range relatedClass.Methods.Funcs {
		for _, calledFunction := range method.Called.Funcs {
			if calledFunction.Class != targetClass {
				continue
			}

			rel.UsedTargetMethods.Add(calledFunction)
			rel.WhereTargetUsedMethods[calledFunction] = method
		}

		for _, usedField := range method.UsedFields.Fields {
			if usedField.Class != targetClass {
				continue
			}

			rel.UsedTargetFields.Add(usedField)
			rel.WhereTargetUsedFields[usedField] = method
		}

		for _, usedConstant := range method.UsedConstants.Constants {
			if usedConstant.Class != targetClass {
				continue
			}

			rel.UsedTargetConstants.Add(usedConstant)
			rel.WhereTargetUsedConstants[usedConstant] = method
		}
	}

	for _, implementedIface := range targetClass.Implements.Classes {
		if relatedClass == implementedIface {
			rel.IsTargetImplements = true
		}
	}

	for _, extendedClass := range targetClass.Extends.Classes {
		if relatedClass == extendedClass {
			rel.IsTargetExtends = true
		}
	}

	for _, implementedIface := range relatedClass.Implements.Classes {
		if targetClass == implementedIface {
			rel.IsRelatedImplements = true
		}
	}

	for _, extendedClass := range relatedClass.Extends.Classes {
		if targetClass == extendedClass {
			rel.IsRelatedExtends = true
		}
	}

	return rel
}
