package representator

import (
	"encoding/json"
	"fmt"

	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/metrics"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

type NamespaceData struct {
	Name     string `json:"name"`
	FullName string `json:"fullName"`

	Files           int64 `json:"files"`
	Classes         int64 `json:"classes"`
	AbstractClasses int64 `json:"abstract-classes"`

	OwnClasses int64 `json:"own-classes"`

	Afferent    float64 `json:"aff"`
	Efferent    float64 `json:"eff"`
	Instability float64 `json:"instab"`

	Abstractness float64

	Childs int64 `json:"childs"`
}

func NamespaceToData(n *symbols.Namespace) *NamespaceData {
	if n == nil {
		return nil
	}

	aff, eff, instab := metrics.AfferentEfferentStabilityOfNamespace(n)
	abstractness := metrics.AbstractnessOfNamespace(n)

	abstractClasses, allClasses := n.CountAbstractAndAllClasses()

	return &NamespaceData{
		Name:            n.Name,
		FullName:        n.FullName,
		Files:           int64(n.Files.Len()),
		Classes:         allClasses,
		AbstractClasses: abstractClasses,
		OwnClasses:      int64(n.Classes.Len()),
		Afferent:        aff,
		Efferent:        eff,
		Instability:     instab,
		Abstractness:    abstractness,
		Childs:          int64(n.Childs.Len()),
	}
}

func GetShortStringNamespaceRepr(f *symbols.Namespace) string {
	if f == nil {
		return ""
	}

	data := NamespaceToData(f)

	return fmt.Sprintf("Namespace %s", data.FullName)
}

func GetStringNamespaceRepr(n *symbols.Namespace) string {
	if n == nil {
		return ""
	}

	data := NamespaceToData(n)

	var res string

	res += fmt.Sprintf("Namespace %s\n", data.FullName)

	res += cfmt.Sprintf("  {{Files}}::green:        %s\n", ColorOutputIntZeroableValue(data.Files))
	res += cfmt.Sprintf("  {{Classes}}::green:      %s\n", ColorOutputIntZeroableValue(data.Classes))
	res += cfmt.Sprintf("    {{Abstract}}::green:   %s %s\n", ColorOutputIntZeroableValue(data.AbstractClasses), ColorOutputFloatZeroablePercentValue(utils.Percent(data.AbstractClasses, data.Classes)))
	res += cfmt.Sprintf("    {{Own}}::green:        %s %s\n", ColorOutputIntZeroableValue(data.OwnClasses), ColorOutputFloatZeroablePercentValue(utils.Percent(data.OwnClasses, data.Classes)))

	res += cfmt.Sprintf("  {{Afferent}}::green:     %s\n", ColorOutputFloatZeroableValue(data.Afferent))
	res += cfmt.Sprintf("  {{Efferent}}::green:     %s\n", ColorOutputFloatZeroableValue(data.Efferent))
	res += cfmt.Sprintf("  {{Instability}}::green:  %s\n", ColorOutputFloatZeroableValue(data.Instability))
	res += cfmt.Sprintf("  {{Abstractness}}::green: %s\n", ColorOutputFloatZeroableValue(data.Abstractness))
	res += cfmt.Sprintf("  {{Childs}}::green:       %s\n", ColorOutputIntZeroableValue(data.Childs))

	return res
}

func GetJsonNamespaceRepr(f *symbols.Namespace) (string, error) {
	data := NamespaceToData(f)

	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetPrettifyJsonNamespacesRepr(n []*symbols.Namespace) (string, error) {
	data := make([]*NamespaceData, 0, len(n))

	for _, namespace := range n {
		data = append(data, NamespaceToData(namespace))
	}

	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetJsonNamespaceReprWithFlag(f *symbols.Namespace) (string, error) {
	type Response struct {
		Data  *NamespaceData `json:"data"`
		Found bool           `json:"found"`
	}
	var resp Response

	resp.Data = NamespaceToData(f)
	resp.Found = f != nil

	res, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
