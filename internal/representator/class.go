package representator

import (
	"encoding/json"
	"fmt"

	"github.com/i582/phpstats/internal/stats/metrics"
	"github.com/i582/phpstats/internal/stats/symbols"
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

	implements *symbols.Classes
	extends    *symbols.Classes

	fields    *symbols.Fields
	methods   *symbols.Functions
	constants *symbols.Constants
}

func classToData(c *symbols.Class) *ClassData {
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

func GetShortClassRepr(c *symbols.Class) string {
	if c == nil {
		return ""
	}

	data := classToData(c)

	return fmt.Sprintf("%s %s", data.Type, data.Name)
}

func GetStringClassRepr(c *symbols.Class) string {
	if c == nil {
		return ""
	}

	data := classToData(c)

	var res string

	res += fmt.Sprintf("%s %s\n", data.Type, data.Name)
	res += fmt.Sprintf("  File:                          %s\n", data.File)
	res += fmt.Sprintf("  Afferent coupling:             %.2f\n", data.Afferent)
	res += fmt.Sprintf("  Efferent coupling:             %.2f\n", data.Efferent)
	res += fmt.Sprintf("  Instability:                   %.2f\n", data.Instability)
	res += fmt.Sprintf("  Lack of Cohesion in Methods:   %.2f\n", data.Lcom)
	res += fmt.Sprintf("  Lack of Cohesion in Methods 4: %d\n", data.Lcom4)
	res += fmt.Sprintf("  Count class dependencies:      %d\n", data.CountDeps)
	res += fmt.Sprintf("  Count dependent classes:       %d\n", data.CountDepsBy)

	return res
}

func GetJsonClassRepr(c *symbols.Class) (string, error) {
	data := classToData(c)

	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetPrettifyJsonClassesRepr(c []*symbols.Class) (string, error) {
	data := make([]*ClassData, 0, len(c))

	for _, class := range c {
		data = append(data, classToData(class))
	}

	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetJsonClassReprWithFlag(c *symbols.Class) (string, error) {
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
