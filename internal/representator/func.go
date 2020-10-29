package representator

import (
	"encoding/json"
	"fmt"

	"github.com/i582/phpstats/internal/stats"
)

type FunctionData struct {
	Name stats.FuncKey `json:"name"`
	Type string        `json:"type"`

	UsesCount int64 `json:"uses-count"`

	CountCalled   int64 `json:"count-called"`
	CountCalledBy int64 `json:"count-called-by"`

	CountDeps   int64 `json:"count-deps"`
	CountDepsBy int64 `json:"count-deps-by"`
}

func funcToData(f *stats.Function) *FunctionData {
	var tp string
	if f.Name.IsMethod() {
		tp = "Method"
	} else {
		tp = "Function"
	}

	return &FunctionData{
		Name:          f.Name,
		Type:          tp,
		UsesCount:     f.UsesCount,
		CountCalled:   int64(f.Called.Len()),
		CountCalledBy: int64(f.CalledBy.Len()),
		CountDeps:     f.CountDeps(),
		CountDepsBy:   f.CountDepsBy(),
	}
}

func GetShortFunctionRepr(f *stats.Function) string {
	data := funcToData(f)

	return fmt.Sprintf("%s %s", data.Type, data.Name)
}

func GetFunctionRepr(f *stats.Function) string {
	data := funcToData(f)

	var res string

	res += fmt.Sprintf("%s %s\n", data.Type, data.Name)
	if data.Name.IsMethod() {
		res += fmt.Sprintf("  Class:               %s\n", data.Name.ClassName)
	}
	res += fmt.Sprintf("  Number of uses:      %d\n", data.UsesCount)
	res += fmt.Sprintf("  Depends of classes:  %d\n", data.CountDeps)
	res += fmt.Sprintf("  Classes depends:     %d\n", data.CountDepsBy)
	res += fmt.Sprintf("  Called functions:    %d\n", data.CountCalled)
	res += fmt.Sprintf("  Called by functions: %d\n", data.CountCalledBy)

	return res
}

func GetJsonFunctionRepr(f *stats.Function) (string, error) {
	data := funcToData(f)

	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
