package representator

import (
	"encoding/json"
	"fmt"

	"github.com/i582/phpstats/internal/stats"
	"github.com/i582/phpstats/internal/stats/metrics"
)

type NamespaceData struct {
	Name     string `json:"name"`
	FullName string `json:"fullName"`

	Files   int64 `json:"files"`
	Classes int64 `json:"classes"`

	Aff    float64 `json:"aff"`
	Eff    float64 `json:"eff"`
	Instab float64 `json:"instab"`

	Childs int64 `json:"childs"`
}

func namespaceToData(n *stats.Namespace) *NamespaceData {
	if n == nil {
		return nil
	}

	aff, eff, instab := metrics.AfferentEfferentStabilityOfNamespace(n)

	return &NamespaceData{
		Name:     n.Name,
		FullName: n.FullName,
		Files:    int64(n.Files.Len()),
		Classes:  int64(n.Classes.Len()),
		Aff:      aff,
		Eff:      eff,
		Instab:   instab,
		Childs:   int64(n.Childs.Len()),
	}
}

func GetShortStringNamespaceRepr(f *stats.Namespace) string {
	if f == nil {
		return ""
	}

	data := namespaceToData(f)

	return fmt.Sprintf("Namespace %s", data.FullName)
}

func GetStringNamespaceRepr(n *stats.Namespace) string {
	if n == nil {
		return ""
	}

	data := namespaceToData(n)

	var res string

	res += fmt.Sprintf("Namespace %s\n", data.FullName)

	res += fmt.Sprintf("  Files:       %d\n", data.Files)
	res += fmt.Sprintf("  Classes:     %d\n", data.Classes)
	res += fmt.Sprintf("  Afferent:    %f\n", data.Aff)
	res += fmt.Sprintf("  Efferent:    %f\n", data.Eff)
	res += fmt.Sprintf("  Instability: %f\n", data.Instab)
	res += fmt.Sprintf("  Childs:      %d\n", data.Childs)

	return res
}

func GetJsonNamespaceRepr(f *stats.Namespace) (string, error) {
	data := namespaceToData(f)

	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetJsonNamespaceReprWithFlag(f *stats.Namespace) (string, error) {
	type Response struct {
		Data  *NamespaceData `json:"data"`
		Found bool           `json:"found"`
	}
	var resp Response

	resp.Data = namespaceToData(f)
	resp.Found = f != nil

	res, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
