package templates

import (
	"fmt"
	"strings"

	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func TemplateClassNode(c *symbols.Class) *graph.Node {
	name := utils.NameToIdentifier(c.Name)

	nameParts := strings.Split(c.Name, `\`)
	className := nameParts[len(nameParts)-1]
	nsName := strings.Join(nameParts[0:len(nameParts)-1], `\\`)
	if nsName == "" {
		nsName = "global scope"
	}

	tp := c.Type()

	label := "(" + tp + ")\\n" + nsName + "\\n" + className + "\\n(links: " + fmt.Sprint(c.Deps.Len()) + ")"

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
