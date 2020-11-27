package templates

import (
	"fmt"

	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func TemplateFileNode(f *symbols.File) *graph.Node {
	name := utils.NameToIdentifier(f.Path)

	allInclude := f.RequiredBlock.Len() + f.RequiredRoot.Len()
	if allInclude == 0 {
		allInclude = 1
	}
	blockPercent := float64(f.RequiredBlock.Len()) / float64(allInclude) * 100
	rootPercent := float64(f.RequiredRoot.Len()) / float64(allInclude) * 100
	label := "file\\n" + f.Name + fmt.Sprintf("\\n(block: %d [%.2f%%], root: %d [%.2f%%])", f.RequiredBlock.Len(), blockPercent, f.RequiredRoot.Len(), rootPercent)

	return &graph.Node{
		Name: name,
		Styles: graph.NodeStyles{
			Label:     label,
			Shape:     "rect",
			FillColor: DefaultFillColor,
			EdgeColor: DefaultOutlineColor,
			Style:     "filled",
			FontSize:  12,
		},
	}
}
