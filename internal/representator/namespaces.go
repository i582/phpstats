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
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Classes")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Aff}}::green\n{{coup}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Eff}}::green\n{{coup}}::green")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Instab")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Childs")},
		},
	}

	for index, namespace := range n {
		data := namespaceToData(namespace)

		name := data.FullName
		name = splitText(name)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: color.Gray.Sprint(int64(index+1) + offset)},
			{Text: name},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.Files)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.Classes)},
			{Align: simpletable.AlignRight, Text: colorOutputFloatZeroableValue(data.Aff)},
			{Align: simpletable.AlignRight, Text: colorOutputFloatZeroableValue(data.Eff)},
			{Align: simpletable.AlignRight, Text: colorOutputFloatZeroableValue(data.Instab)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.Childs)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table.String()
}
