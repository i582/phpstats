package representator

import (
	"encoding/json"
	"fmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

type FunctionData struct {
	Name string `json:"name"`
	Type string `json:"type"`

	Class string `json:"className"`

	UsesCount int64 `json:"usesCount"`

	CountCalled   int64 `json:"countCalled"`
	CountCalledBy int64 `json:"countCalledBy"`

	CountDeps   int64 `json:"countDeps"`
	CountDepsBy int64 `json:"countDepsBy"`

	CyclomaticComplexity int64 `json:"cc"`
}

func funcToData(f *symbols.Function) *FunctionData {
	if f == nil {
		return nil
	}

	var tp string
	if f.Name.IsMethod() {
		tp = "Method"
	} else {
		tp = "Function"
	}

	return &FunctionData{
		Name:                 f.Name.String(),
		Type:                 tp,
		Class:                f.Name.ClassName,
		UsesCount:            f.UsesCount,
		CountCalled:          int64(f.Called.Len()),
		CountCalledBy:        int64(f.CalledBy.Len()),
		CountDeps:            f.CountDeps(),
		CountDepsBy:          f.CountDepsBy(),
		CyclomaticComplexity: f.CyclomaticComplexity,
	}
}

func GetShortStringFunctionRepr(f *symbols.Function) string {
	if f == nil {
		return ""
	}

	data := funcToData(f)

	return fmt.Sprintf("%s %s", data.Type, data.Name)
}

func GetStringFunctionRepr(f *symbols.Function) string {
	if f == nil {
		return ""
	}

	data := funcToData(f)

	var res string

	res += fmt.Sprintf("%s %s\n", data.Type, data.Name)
	if data.Class != "" {
		res += fmt.Sprintf("  Class:                 %s\n", data.Class)
	}
	res += fmt.Sprintf("  Number of uses:        %d\n", data.UsesCount)
	res += fmt.Sprintf("  Depends of classes:    %d\n", data.CountDeps)
	res += fmt.Sprintf("  Classes depends:       %d\n", data.CountDepsBy)
	res += fmt.Sprintf("  Called functions:      %d\n", data.CountCalled)
	res += fmt.Sprintf("  Called by functions:   %d\n", data.CountCalledBy)
	res += fmt.Sprintf("  Cyclomatic complexity: %d (>15 hard to understand, >30 extremply complex)\n", data.CyclomaticComplexity)

	return res
}

func GetJsonFunctionRepr(f *symbols.Function) (string, error) {
	data := funcToData(f)

	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetJsonFunctionReprWithFlag(f *symbols.Function) (string, error) {
	type Response struct {
		Data  *FunctionData `json:"data"`
		Found bool          `json:"found"`
	}
	var resp Response

	resp.Data = funcToData(f)
	resp.Found = f != nil

	res, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
