package templates

import (
	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func splitText(text string) string {
	if len(text) > 20 {
		indexOfSlash := len(text) / 2
		for i := indexOfSlash; i >= 0; i-- {
			if text[i] == '\\' {
				indexOfSlash = i + 1
				break
			}
		}

		text = text[:indexOfSlash] + " \\n" + text[indexOfSlash:]
	}

	return text
}

func TemplateNamespaceNode(c *symbols.Namespace) *graph.Node {
	name := utils.NameToIdentifier(c.FullName)
	label := utils.NormalizeSlashes(c.FullName)

	label = splitText(label)

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
