package representator

import (
	"bytes"
	"fmt"

	"github.com/i582/phpstats/internal/stats/symbols"

	"github.com/olekukonko/tablewriter"
)

func GetTableClassesRepr(c []*symbols.Class, offset int64) string {
	if c == nil {
		return ""
	}

	w := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"#", "Name", "Aff\ncoup", "Eff\ncoup", "Instab", "LCOM", "LCOM 4", "Class\ndeps", "Classes\ndepends"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)

	for index, class := range c {
		data := classToData(class)

		var lcom string
		if data.Lcom == -1 {
			lcom = "undef"
		} else {
			lcom = fmt.Sprintf("%.2f", data.Lcom)
		}

		name := data.Name
		name = SplitText(name)

		line := []string{fmt.Sprint(int64(index+1) + offset), name, fmt.Sprint(data.Afferent), fmt.Sprint(data.Efferent),
			fmt.Sprintf("%.2f", data.Instability), lcom, fmt.Sprint(data.Lcom4), fmt.Sprint(data.CountDeps), fmt.Sprint(data.CountDepsBy)}
		table.Append(line)
	}

	table.Render()
	return w.String()
}
