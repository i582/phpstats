package templates

import (
	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func TemplateFunctionNode(c *symbols.Function) *graph.Node {
	name := utils.NameToIdentifier(c.Name.String())
	label := c.Name.Name

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
