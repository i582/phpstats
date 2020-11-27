package representator

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"

	"github.com/i582/phpstats/internal/stats/symbols"
)

func GetTableFilesRepr(f []*symbols.File, offset int64) string {
	if f == nil {
		return ""
	}

	w := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"#", "Name", "Root\ninclusions", "Block\ninclusions", "Count\nrequired by"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)

	for index, file := range f {
		data := fileToData(file)

		name := data.Name
		name = SplitText(name)

		line := []string{fmt.Sprint(int64(index+1) + offset), name, fmt.Sprint(data.CountRequiredRoot), fmt.Sprint(data.CountRequiredBlock), fmt.Sprint(data.CountRequiredBy)}
		table.Append(line)
	}

	table.Render()
	return w.String()
}
