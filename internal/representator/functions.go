package representator

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"

	"github.com/i582/phpstats/internal/stats/symbols"
)

func GetTableFunctionsRepr(f []*symbols.Function, offset int64) string {
	if f == nil {
		return ""
	}

	w := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"#", "Name", "Number\nof uses", "Deps\nclasses", "Classes\ndepends",
		"Called\nfuncs", "Called\nby funcs", "Cyclo\ncompl", "Magic\nnums"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)

	for index, fun := range f {
		data := funcToData(fun)

		name := data.Name
		name = SplitText(name)

		line := []string{fmt.Sprint(int64(index+1) + offset), name, fmt.Sprint(data.UsesCount), fmt.Sprint(data.CountDeps),
			fmt.Sprint(data.CountDepsBy), fmt.Sprint(data.CountCalled), fmt.Sprint(data.CountCalledBy), fmt.Sprint(data.CyclomaticComplexity), fmt.Sprint(data.CountMagicNumbers)}
		table.Append(line)
	}

	table.Render()
	return w.String()
}
