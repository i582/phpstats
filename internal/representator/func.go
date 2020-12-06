package representator

import (
	"encoding/json"
	"fmt"

	"github.com/i582/cfmt"

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
	CountMagicNumbers    int64 `json:"cmn"`

	FullyTypes bool `json:"fully_typed"`
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
		CountMagicNumbers:    f.CountMagicNumbers,
		FullyTypes:           f.FullyTyped,
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
		res += cfmt.Sprintf("  {{Class}}::green:                 %s\n", data.Class)
	}
	res += cfmt.Sprintf("  {{Number of uses}}::green:        %s\n", ColorOutputIntZeroableValue(data.UsesCount))
	res += cfmt.Sprintf("  {{Depends of classes}}::green:    %s\n", ColorOutputIntZeroableValue(data.CountDeps))
	res += cfmt.Sprintf("  {{Classes depends}}::green:       %s\n", ColorOutputIntZeroableValue(data.CountDepsBy))
	res += cfmt.Sprintf("  {{Called functions}}::green:      %s\n", ColorOutputIntZeroableValue(data.CountCalled))
	res += cfmt.Sprintf("  {{Called by functions}}::green:   %s\n", ColorOutputIntZeroableValue(data.CountCalledBy))
	res += cfmt.Sprintf("  {{Cyclomatic complexity}}::green: %s {{(>15 hard to understand, >30 extremely complex)}}::gray\n", ColorOutputIntZeroableValue(data.CyclomaticComplexity))
	res += cfmt.Sprintf("  {{Count magic numbers}}::green:   %s\n", ColorOutputIntZeroableValue(data.CountMagicNumbers))
	res += cfmt.Sprintf("  {{Fully typed}}::green:           %s\n", ColorOutputBoolZeroableValue(data.FullyTypes))

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

func GetPrettifyJsonFunctionsRepr(f []*symbols.Function) (string, error) {
	data := make([]*FunctionData, 0, len(f))

	for _, function := range f {
		data = append(data, funcToData(function))
	}

	res, err := json.MarshalIndent(data, "", "\t")
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
