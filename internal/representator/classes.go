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
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Instab}}::green\n{{ility}}::green")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("LCOM")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{LCOM}}::green\n{{4}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Class}}::green\n{{deps}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Classes}}::green\n{{depends}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Fully}}::green\n{{typed}}::green\n{{methods}}::green")},
		},
	}

	for index, class := range c {
		data := ClassToData(class)

		var lcom string
		if data.Lcom == -1 {
			lcom = color.Gray.Sprint("undef")
		} else {
			lcom = ColorOutputFloatZeroableValue(data.Lcom)
		}

		name := data.Name
		name = splitText(name)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: color.Gray.Sprint(int64(index+1) + offset)},
			{Text: name},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(int64(data.Afferent))},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(int64(data.Efferent))},
			{Align: simpletable.AlignRight, Text: ColorOutputFloatZeroableValue(data.Instability)},
			{Align: simpletable.AlignRight, Text: lcom},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.Lcom4)},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.CountDeps)},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.CountDepsBy)},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.CountFullyTypedMethods) + color.Gray.Sprintf("(%d)", data.methods.Len())},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table.String()
}
