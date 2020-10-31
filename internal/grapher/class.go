package grapher

import (
	"fmt"

	uml "github.com/i582/phpstats/internal/grapher/unl"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func (g *Grapher) ClassDeps(c *symbols.Class, maxRecursion int64) string {
	var res string

	res += graphHeader

	main, sub := g.classDepsRecursive(c, 0, maxRecursion, visitedMap{})

	res += g.subGraphVendorWrapper(sub)
	res += main

	return g.graphWrapper(res, utils.NameToIdentifier(c.Name))
}

func (g *Grapher) classDepsRecursive(c *symbols.Class, levelRecursion, maxRecursion int64, visited visitedMap) (string, string) {
	var res string
	var sub string

	classUml := uml.GetUmlForClassWithFilter(c, func(m *symbols.Function) bool {
		return m.Deps().Len() != 0
	}, func(f *symbols.Field) bool {
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
		str := fmt.Sprintf("   %s -> %s\n", utils.NameToIdentifier(c.Name), utils.NameToIdentifier(implement.Name))

		if _, ok := visited[str]; ok {
			continue
		}
		visited[str] = struct{}{}

		res += g.outputClassImplement(str)

		mn, sb := g.classDepsRecursive(implement, levelRecursion+1, maxRecursion, visited)
		res += mn
		sub += sb
	}

	for _, field := range c.Fields.Fields {
		for caller := range field.Used {
			for _, class := range caller.Deps().Classes {
				str := fmt.Sprintf("   %s:fields -> %s\n", utils.NameToIdentifier(c.Name), utils.NameToIdentifier(class.Name))

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
			str := fmt.Sprintf("   %s:methods -> %s\n", utils.NameToIdentifier(c.Name), utils.NameToIdentifier(class.Name))

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

func (g *Grapher) outputClassImplement(str string) string {
	var res string
	res += "\tedge [style=\"dashed\"];"
	res += str
	res += "\tedge [style=\"solid\"];"
	return res
}
