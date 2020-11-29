package representator

import (
	"encoding/json"
	"fmt"

	"github.com/i582/cfmt"

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

func ClassToData(c *symbols.Class) *ClassData {
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

	aff, eff, instab := metrics.AfferentEfferentInstabilityOfClass(c)
	lcom, _ := metrics.LackOfCohesionInMethods(c)
	lcom4 := metrics.LackOfCohesionInMethods4(c)

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

func GetStringClassRepr(c *symbols.Class) string {
	if c == nil {
		return ""
	}

	data := ClassToData(c)

	var res string

	res += fmt.Sprintf("%s %s\n", data.Type, data.Name)
	res += cfmt.Sprintf("  {{Afferent coupling}}::green:             %s\n", ColorOutputFloatZeroableValue(data.Afferent))
	res += cfmt.Sprintf("  {{Efferent coupling}}::green:             %s\n", ColorOutputFloatZeroableValue(data.Efferent))
	res += cfmt.Sprintf("  {{Instability}}::green:                   %s\n", ColorOutputFloatZeroableValue(data.Instability))
	res += cfmt.Sprintf("  {{Lack of Cohesion in Methods}}::green:   %s\n", ColorOutputFloatZeroableValue(data.Lcom))
	res += cfmt.Sprintf("  {{Lack of Cohesion in Methods 4}}::green: %s\n", ColorOutputIntZeroableValue(data.Lcom4))
	res += cfmt.Sprintf("  {{Count class dependencies}}::green:      %s\n", ColorOutputIntZeroableValue(data.CountDeps))
	res += cfmt.Sprintf("  {{Count dependent classes}}::green:       %s\n", ColorOutputIntZeroableValue(data.CountDepsBy))

	return res
}

func GetJsonClassRepr(c *symbols.Class) (string, error) {
	data := ClassToData(c)

	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetPrettifyJsonClassesRepr(c []*symbols.Class) (string, error) {
	data := make([]*ClassData, 0, len(c))

	for _, class := range c {
		data = append(data, ClassToData(class))
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

	resp.Data = ClassToData(c)
	resp.Found = c != nil

	res, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
