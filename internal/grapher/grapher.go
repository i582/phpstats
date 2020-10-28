package grapher

import (
	"fmt"

	"github.com/i582/phpstats/internal/stats"
	"github.com/i582/phpstats/internal/utils"
)

const defaultColor = "#D5D5D5"
const graphHeader = `
	size="5,5"
	node[shape=record,style=filled,fillcolor="#cccccc"]
	edge[arrowtail=empty]
`
const subGraphHeader = `
	style=filled;
`

type visitedMap map[string]struct{}

type Grapher struct{}

func NewGrapher() *Grapher {
	return &Grapher{}
}

func (g *Grapher) graphWrapper(str string, name string) string {
	var res string
	res += "digraph " + name + "{\n"
	res += str
	res += "}\n"
	return res
}

func (g *Grapher) subGraphVendorWrapper(str string) string {
	return g.subGraphWrapper(str, "vendor")
}

func (g *Grapher) subGraphWrapper(str string, name string) string {
	return g.subGraphWrapperColor(str, name, "#eeeeee")
}

func (g *Grapher) subGraphWrapperColor(str string, name string, color string) string {
	var res string
	res += "\tsubgraph cluster_" + utils.ClassNameNormalize(name) + "{\n"
	res += "\tlabel=\"" + name + "\";\n"
	res += "\tfillcolor=\"" + color + "\";\n"
	res += "\t" + subGraphHeader + "\n"
	res += "\t" + str
	res += "\t}\n"
	return res
}

func (g *Grapher) outputWithColor(str string, newColor, oldColor string) string {
	var res string

	if newColor == oldColor {
		return str
	}

	res += fmt.Sprintf("\tnode[shape=record,style=filled,fillcolor=\"%s\"]\n", newColor)
	res += str
	res += fmt.Sprintf("\tnode[shape=record,style=filled,fillcolor=\"%s\"]\n", oldColor)
	return res
}

func (g *Grapher) getColorForClass(c *stats.Class) string {
	if c.IsInterface {
		return "#bbbbbb"
	}

	if c.IsAbstract {
		return "#888888"
	}

	return defaultColor
}
