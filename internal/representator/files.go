package representator

import (
	"github.com/alexeyco/simpletable"
	"github.com/gookit/color"
	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

func GetTableFilesRepr(f []*symbols.File, offset int64) string {
	if f == nil {
		return ""
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("#")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Name")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Root}}::green\n{{inclusions}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Block}}::green\n{{inclusions}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Count}}::green\n{{required by}}::green")},
		},
	}

	for index, file := range f {
		data := fileToData(file)

		name := data.Name
		name = splitText(name)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: color.Gray.Sprint(int64(index+1) + offset)},
			{Text: name},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.CountRequiredRoot)},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.CountRequiredBlock)},
			{Align: simpletable.AlignRight, Text: ColorOutputIntZeroableValue(data.CountRequiredBy)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table.String()
}
