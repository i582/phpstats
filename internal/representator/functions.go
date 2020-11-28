package representator

import (
	"github.com/alexeyco/simpletable"
	"github.com/gookit/color"
	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

func GetTableFunctionsRepr(f []*symbols.Function, offset int64) string {
	if f == nil {
		return ""
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("#")},
			{Align: simpletable.AlignCenter, Text: color.Green.Sprint("Name")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Number}}::green\n{{of uses}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Deps}}::green\n{{classes}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Classes}}::green\n{{depends}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Called}}::green\n{{funcs}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Called}}::green\n{{by funcs}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Cyclo}}::green\n{{compl}}::green")},
			{Align: simpletable.AlignCenter, Text: cfmt.Sprint("{{Magic}}::green\n{{nums}}::green")},
		},
	}

	for index, fun := range f {
		data := funcToData(fun)

		name := data.Name
		name = splitText(name)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: color.Gray.Sprint(int64(index+1) + offset)},
			{Text: name},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.UsesCount)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.CountDeps)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.CountDepsBy)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.CountCalled)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.CountCalledBy)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.CyclomaticComplexity)},
			{Align: simpletable.AlignRight, Text: colorOutputIntZeroableValue(data.CountMagicNumbers)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table.String()
}
