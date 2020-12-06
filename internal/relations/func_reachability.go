package relations

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

type ReachabilityExcludedMap map[*symbols.Function]*symbols.Function

type ReachabilityFunctionJsonResult struct {
	ParentFunction    string     `json:"parentFunctions"`
	ChildFunction     string     `json:"childFunctions"`
	ExcludedFunctions []string   `json:"excluded"`
	Count             int64      `json:"count"`
	Offset            int64      `json:"offset"`
	Paths             [][]string `json:"paths"`
}

type ReachabilityFunctionResult struct {
	ParentFunction *symbols.Function
	ChildFunction  *symbols.Function
	Reachable      bool

	Paths             [][]*symbols.Function
	ExcludedFunctions ReachabilityExcludedMap

	PrintPaths  bool
	PrintCount  int64
	PrintOffset int64
}

func NewReachabilityFunctionResult() *ReachabilityFunctionResult {
	return &ReachabilityFunctionResult{}
}

func functionsPathsToStringPaths(paths [][]*symbols.Function) [][]string {
	var res [][]string

	for _, path := range paths {
		tempPath := make([]string, 0, len(path))

		for _, function := range path {
			tempPath = append(tempPath, function.Name.String())
		}

		res = append(res, tempPath)
	}

	return res
}

func excludedMapToString(excluded ReachabilityExcludedMap) []string {
	var res []string

	for _, function := range excluded {
		res = append(res, function.Name.String())
	}

	return res
}

func (r *ReachabilityFunctionResult) Json() ([]byte, error) {
	return json.MarshalIndent(&ReachabilityFunctionJsonResult{
		ParentFunction:    r.ParentFunction.Name.String(),
		ChildFunction:     r.ChildFunction.Name.String(),
		ExcludedFunctions: excludedMapToString(r.ExcludedFunctions),
		Paths:             functionsPathsToStringPaths(r.Paths),
		Count:             int64(len(r.Paths)),
		Offset:            r.PrintOffset,
	}, "", "\t")
}

func (r *ReachabilityFunctionResult) String() string {
	var res string

	res += cfmt.Sprintf("Is function {{%s}}::green reachable from function {{%s}}::yellow: %t", r.ChildFunction.Name, r.ParentFunction.Name, r.Reachable)
	if r.Reachable {
		res += fmt.Sprintf(" (%d paths)\n", len(r.Paths))
	} else {
		res += "\n"
	}

	if r.Reachable && r.PrintPaths {
		sort.Slice(r.Paths, func(i, j int) bool {
			return len(r.Paths[i]) < len(r.Paths[j])
		})

		neededPaths := r.Paths
		if r.PrintCount+r.PrintOffset < int64(len(neededPaths)) {
			neededPaths = neededPaths[:r.PrintCount+r.PrintOffset]
		}

		if r.PrintOffset < int64(len(neededPaths)) {
			neededPaths = neededPaths[r.PrintOffset:]
		}

		res += fmt.Sprintf("Showing %d paths out of %d starting from %d:\n\n", len(neededPaths), len(r.Paths), r.PrintOffset+1)

		for _, path := range neededPaths {
			res += "    " + stringCallstack(path) + "\n"
		}
	}

	return res
}

func GetReachabilityFunction(parentFunction *symbols.Function, childFunction *symbols.Function, excludedFunctions ReachabilityExcludedMap, maxDepth int64) *ReachabilityFunctionResult {
	rel := NewReachabilityFunctionResult()
	rel.ParentFunction = parentFunction
	rel.ChildFunction = childFunction
	rel.ExcludedFunctions = excludedFunctions

	rel.Reachable, rel.Paths = calledInCallstack(parentFunction, childFunction, nil, map[*symbols.Function]struct{}{}, 0, maxDepth)

	return rel
}
