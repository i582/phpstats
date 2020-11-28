package representator

import (
	"github.com/alexeyco/simpletable"
	"github.com/gookit/color"
	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

func GetTableClassesRepr(c []*symbols.Class, offset int64) string {
	if c == nil {
		return ""
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("#")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Name")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Aff}}::green\n{{coup}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Eff}}::green\n{{coup}}::green")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Instab")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("LCOM")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("LCOM 4")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Class}}::green\n{{deps}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Classes}}::green\n{{depends}}::green")},
		},
	}

	for index, class := range c {
		data := ClassToData(class)

		var lcom string
		if data.Lcom == -1 {
			lcom = color.Gray.Sprint("undef")
		} else {
			lcom = colorOutputFloatZeroableValue(data.Lcom)
		}

		name := data.Name
		name = splitText(name)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: color.Gray.Sprint(int64(index+1) + offset)},
			{Text: name},
			{Align: simpletable.AlignRight, Text: colorOutputFloatZeroableValue(data.Afferent)},
			{Align: simpletable.AlignRight, Text: colorOutputFloatZeroableValue(data.Efferent)},
			{Align: simpletable.AlignRight, Text: colorOutputFloatZeroableValue(data.Instability)},
			{Align: simpletable.AlignRight, Text: lcom},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.Lcom4)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.CountDeps)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.CountDepsBy)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table.String()
}
