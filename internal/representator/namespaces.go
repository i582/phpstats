package representator

import (
	"github.com/alexeyco/simpletable"
	"github.com/gookit/color"
	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

func GetTableNamespacesRepr(n []*symbols.Namespace, offset int64) string {
	if n == nil {
		return ""
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("#")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Name")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Files")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{All}}::green\n{{classes}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Own}}::green\n{{classes}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Aff}}::green\n{{coup}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Eff}}::green\n{{coup}}::green")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Instab")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Abstract")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Childs")},
		},
	}

	for index, namespace := range n {
		data := NamespaceToData(namespace)

		name := data.FullName
		name = splitText(name)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: color.Gray.Sprint(int64(index+1) + offset)},
			{Text: name},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.Files)},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.Classes)},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.OwnClasses)},
			{Align: simpletable.AlignRight, Text: ColorOutputFloatZeroableValue(data.Afferent)},
			{Align: simpletable.AlignRight, Text: ColorOutputFloatZeroableValue(data.Efferent)},
			{Align: simpletable.AlignRight, Text: ColorOutputFloatZeroableValue(data.Instability)},
			{Align: simpletable.AlignRight, Text: ColorOutputFloatZeroableValue(data.Abstractness)},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.Childs)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table.String()
}
