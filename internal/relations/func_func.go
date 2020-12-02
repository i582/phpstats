package relations

import (
	"fmt"

	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

type Func2FuncRelation struct {
	TargetFunction  *symbols.Function
	RelatedFunction *symbols.Function

	BelongsToSameClass bool
	BelongClass        *symbols.Class

	TargetUsedInRelated bool
	RelatedUsedInTarget bool

	TargetReachableFromRelated      bool
	TargetReachableFromRelatedPaths [][]*symbols.Function

	RelatedReachableFromTarget      bool
	RelatedReachableFromTargetPaths [][]*symbols.Function
}

func NewFunc2FuncRelation() *Func2FuncRelation {
	return &Func2FuncRelation{}
}

func (r *Func2FuncRelation) String() string {
	var res string

	res += cfmt.Sprintf("Function {{%s}}::green connection with function {{%s}}::yellow.\n\n", r.TargetFunction.Name, r.RelatedFunction.Name)

	res += fmt.Sprintf("    Functions belong to the same class:%*s      %t", len(r.RelatedFunction.Name.String()+r.TargetFunction.Name.String()), "", r.BelongsToSameClass)

	if r.BelongsToSameClass {
		res += cfmt.Sprintf(" (%s)\n", r.BelongClass.Name)
	} else {
		res += cfmt.Sprintf("\n")
	}

	res += cfmt.Sprintf("    Function {{%s}}::yellow is used in function {{%s}}::green:          %t\n", r.RelatedFunction.Name, r.TargetFunction.Name, r.RelatedUsedInTarget)
	res += cfmt.Sprintf("    Function {{%s}}::green is used in function {{%s}}::yellow:          %t\n", r.TargetFunction.Name, r.RelatedFunction.Name, r.RelatedUsedInTarget)

	res += cfmt.Sprintf("    Is function {{%s}}::yellow reachable from function {{%s}}::green:   %t\n", r.RelatedFunction.Name, r.TargetFunction.Name, r.RelatedReachableFromTarget)
	if r.RelatedReachableFromTarget {
		res += cfmt.Sprintf("    The function is reachable by the following calls:\n")
		for _, path := range r.RelatedReachableFromTargetPaths {
			res += "        " + stringCallstack(path)
		}
		res += fmt.Sprintln()
	}

	res += cfmt.Sprintf("    Is function {{%s}}::green reachable from function {{%s}}::yellow:   %t\n", r.TargetFunction.Name, r.RelatedFunction.Name, r.TargetReachableFromRelated)
	if r.TargetReachableFromRelated {
		res += cfmt.Sprintf("    The function is reachable by the following calls:\n")
		for _, path := range r.TargetReachableFromRelatedPaths {
			res += "        " + stringCallstack(path)
		}
	}

	return res
}

func stringCallstack(callstack []*symbols.Function) string {
	var res string
	res += fmt.Sprint("[")
	for i, f := range callstack {
		res += fmt.Sprint(f.Name)
		if i != len(callstack)-1 {
			res += fmt.Sprint(" -> ")
		}
	}
	res += fmt.Sprintln("]")
	return res
}

func GetFunc2FuncRelation(targetFunction *symbols.Function, relatedFunction *symbols.Function) *Func2FuncRelation {
	rel := NewFunc2FuncRelation()
	rel.TargetFunction = targetFunction
	rel.RelatedFunction = relatedFunction

	if rel.TargetFunction.Class != nil && rel.TargetFunction.Class == rel.RelatedFunction.Class {
		rel.BelongsToSameClass = true
		rel.BelongClass = rel.TargetFunction.Class
	}

	for _, calledInTarget := range rel.TargetFunction.Called.Funcs {
		if calledInTarget == relatedFunction {
			rel.RelatedUsedInTarget = true
			break
		}
	}

	for _, calledInRelated := range rel.RelatedFunction.Called.Funcs {
		if calledInRelated == relatedFunction {
			rel.TargetUsedInRelated = true
			break
		}
	}

	rel.RelatedReachableFromTarget, rel.RelatedReachableFromTargetPaths = calledInCallstack(targetFunction, relatedFunction, nil, map[*symbols.Function]struct{}{})
	rel.TargetReachableFromRelated, rel.TargetReachableFromRelatedPaths = calledInCallstack(relatedFunction, targetFunction, nil, map[*symbols.Function]struct{}{})

	return rel
}

func calledInCallstack(parent, child *symbols.Function, callstack []*symbols.Function, visited map[*symbols.Function]struct{}) (bool, [][]*symbols.Function) {
	if parent.Called.Len() == 0 {
		return false, nil
	}

	if callstack == nil {
		callstack = []*symbols.Function{parent}
	}

	if parent == child {
		return true, [][]*symbols.Function{callstack}
	}

	var callstacks [][]*symbols.Function

	for _, called := range parent.Called.Funcs {
		newCallstack := copyCallstack(callstack)
		newVisited := copyVisited(visited)

		newCallstack = append(newCallstack, called)

		if _, ok := newVisited[called]; ok {
			continue
		}
		if called == parent {
			continue
		}

		newVisited[called] = struct{}{}

		if called == child {
			callstacks = append(callstacks, newCallstack)
			continue
		}

		call, callstack := calledInCallstack(called, child, newCallstack, newVisited)
		if call {
			callstacks = append(callstacks, callstack...)
		}
	}

	return len(callstacks) != 0, callstacks
}

func copyCallstack(callstack []*symbols.Function) []*symbols.Function {
	tmp := make([]*symbols.Function, len(callstack))
	copy(tmp, callstack)
	return tmp
}

func copyVisited(visited map[*symbols.Function]struct{}) map[*symbols.Function]struct{} {
	targetMap := make(map[*symbols.Function]struct{})

	for key, value := range visited {
		targetMap[key] = value
	}

	return targetMap
}
