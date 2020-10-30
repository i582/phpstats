package representator

import (
	"encoding/json"
	"fmt"

	"github.com/i582/phpstats/internal/stats"
	"github.com/i582/phpstats/internal/stats/metrics"
)

type ClassData struct {
	Name string `json:"name"`
	File string `json:"file"`
	Type string `json:"type"`

	Afferent    float64 `json:"aff"`
	Efferent    float64 `json:"eff"`
	Instability float64 `json:"instab"`

	Lcom  float64 `json:"lcom"`
	Lcom4 int64   `json:"lcom4"`

	CountDeps   int64 `json:"countDeps"`
	CountDepsBy int64 `json:"countDepsBy"`

	implements *stats.Classes
	extends    *stats.Classes

	fields    *stats.Fields
	methods   *stats.Functions
	constants *stats.Constants
}

func classToData(c *stats.Class) *ClassData {
	if c == nil {
		return nil
	}

	var tp string
	if c.IsInterface {
		tp = "Interface"
	} else {
		if c.IsAbstract {
			tp = "Abstract"
		} else {
			tp = "Class"
		}
	}

	aff, eff, instab := metrics.AfferentEfferentStabilityOfClass(c)
	lcom, _ := metrics.LackOfCohesionInMethodsOfCLass(c)
	lcom4 := metrics.Lcom4(c)

	return &ClassData{
		Name:        c.Name,
		File:        c.File.Path,
		Type:        tp,
		Afferent:    aff,
		Efferent:    eff,
		Instability: instab,
		Lcom:        lcom,
		Lcom4:       lcom4,
		CountDeps:   int64(c.Deps.Len()),
		CountDepsBy: int64(c.DepsBy.Len()),

		implements: c.Implements,
		extends:    c.Extends,

		fields:    c.Fields,
		methods:   c.Methods,
		constants: c.Constants,
	}
}

func GetShortClassRepr(c *stats.Class) string {
	if c == nil {
		return ""
	}

	data := classToData(c)

	return fmt.Sprintf("%s %s", data.Type, data.Name)
}

func GetClassRepr(c *stats.Class) string {
	if c == nil {
		return ""
	}

	data := classToData(c)

	var res string

	res += fmt.Sprintf("%s %s\n", data.Type, data.Name)
	res += fmt.Sprintf("  File:        %s\n", data.File)
	res += fmt.Sprintf("  Afferent:    %.2f\n", data.Afferent)
	res += fmt.Sprintf("  Efferent:    %.2f\n", data.Efferent)
	res += fmt.Sprintf("  Instability: %.2f\n", data.Instability)
	res += fmt.Sprintf("  LCOM:        %.2f\n", data.Lcom)
	res += fmt.Sprintf("  LCOM4:       %d\n", data.Lcom4)
	res += fmt.Sprintf("  Deps:        %d\n", data.CountDeps)
	res += fmt.Sprintf("  DepsBy:      %d\n", data.CountDepsBy)

	return res
}

func GetJsonClassRepr(c *stats.Class) (string, error) {
	data := classToData(c)

	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetJsonClassReprWithFlag(c *stats.Class) (string, error) {
	type Response struct {
		Data  *ClassData `json:"data"`
		Found bool       `json:"found"`
	}
	var resp Response

	resp.Data = classToData(c)
	resp.Found = c != nil

	res, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
