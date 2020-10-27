package grapher

import (
	"fmt"
	"strings"

	"github.com/i582/phpstats/internal/stats"
)

const defaultColor = "#D5D5D5"
const graphHeader = `
	size="5,5"
	node[shape=record,style=filled,fillcolor="#cccccc"]
	edge[arrowtail=empty]
`
const subGraphHeader = `
	label="Vendor";
	fillcolor="#eeeeee";
	style=filled;
`

type visitedMap map[string]struct{}

type Grapher struct{}

func NewGrapher() *Grapher {
	return &Grapher{}
}

func (g *Grapher) ClassDeps(c *stats.Class, maxRecursion int64) string {
	var res string

	res += graphHeader

	main, sub := g.classDepsRecursive(c, 0, maxRecursion, visitedMap{})

	res += g.subGraphWrapper(sub)
	res += main

	return g.graphWrapper(res)
}

func (g *Grapher) classDepsRecursive(c *stats.Class, levelRecursion, maxRecursion int64, visited visitedMap) (string, string) {
	var res string
	var sub string

	classUml := getUmlForClassWithFilter(c, func(m *stats.Function) bool {
		return m.Deps().Len() != 0
	}, func(f *stats.Field) bool {
		return len(f.Used) != 0
	})

	umlGraph := g.outputWithColor("   "+classUml, g.getColorForClass(c), defaultColor)

	if _, ok := visited[c.Name]; !ok {
		if c.Vendor {
			sub += umlGraph
		} else {
			res += umlGraph
		}
		visited[c.Name] = struct{}{}
	}

	if levelRecursion > maxRecursion {
		return res, sub
	}

	for _, implement := range c.Implements.Classes {
		str := fmt.Sprintf("   %s -> %s\n", transformName(c.Name), transformName(implement.Name))

		if _, ok := visited[str]; ok {
			continue
		}
		visited[str] = struct{}{}

		res += g.outputImplement(str)

		mn, sb := g.classDepsRecursive(implement, levelRecursion+1, maxRecursion, visited)
		res += mn
		sub += sb
	}

	for _, field := range c.Fields.Fields {
		for caller := range field.Used {
			for _, class := range caller.Deps().Classes {
				str := fmt.Sprintf("   %s:fields -> %s\n", transformName(c.Name), transformName(class.Name))

				if _, ok := visited[str]; ok {
					continue
				}
				visited[str] = struct{}{}

				mn, sb := g.classDepsRecursive(class, levelRecursion+1, maxRecursion, visited)
				res += mn
				sub += sb

				res += str
			}
		}
	}

	for _, method := range c.Methods.Funcs {
		deps := method.Deps()
		if deps.Len() == 0 {
			continue
		}

		for _, class := range deps.Classes {
			str := fmt.Sprintf("   %s:methods -> %s\n", transformName(c.Name), transformName(class.Name))

			if _, ok := visited[str]; ok {
				continue
			}
			visited[str] = struct{}{}

			mn, sb := g.classDepsRecursive(class, levelRecursion+1, maxRecursion, visited)
			res += mn
			sub += sb

			res += str
		}
	}

	return res, sub
}

func (g *Grapher) graphWrapper(str string) string {
	var res string
	res += "digraph test{\n"
	res += str
	res += "}\n"
	return res
}

func (g *Grapher) subGraphWrapper(str string) string {
	var res string
	res += "\tsubgraph cluster_vendor {\n"
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

func (g *Grapher) outputImplement(str string) string {
	var res string
	res += "\tedge [style=\"dashed\"];"
	res += str
	res += "\tedge [style=\"solid\"];"
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

func transformName(name string) string {
	return strings.ReplaceAll(name, `\`, `_`)
}
