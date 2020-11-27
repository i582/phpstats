package representator

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"

	"github.com/i582/phpstats/internal/stats/symbols"
)

func GetTableNamespacesRepr(n []*symbols.Namespace, offset int64) string {
	if n == nil {
		return ""
	}

	w := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"#", "Name", "Files", "Classes", "Aff\ncoup", "Eff\ncoup", "Instab", "Childs"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)

	for index, namespace := range n {
		data := namespaceToData(namespace)

		name := data.FullName
		name = SplitText(name)

		line := []string{fmt.Sprint(int64(index+1) + offset), name, fmt.Sprint(data.Files), fmt.Sprint(data.Classes),
			fmt.Sprintf("%.2f", data.Aff), fmt.Sprintf("%.2f", data.Eff), fmt.Sprintf("%.2f", data.Instab), fmt.Sprint(data.Childs)}
		table.Append(line)
	}

	table.Render()
	return w.String()
}
