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

func (g *Grapher) Lcom4(c *symbols.Class) string {
	var res string

	res += graphHeader

	main := g.lcom4(c)

	res += g.subGraphWrapper(main, "Lack of Cohesion in Methods 4 (LCOM4) graph for "+utils.NormalizeSlashes(c.Name))

	return g.graphWrapper(res, utils.NameToIdentifier(c.Name))
}

func (g *Grapher) lcom4(c *symbols.Class) string {
	var res string

	showed := map[string]struct{}{}

	for _, method := range c.Methods.Funcs {
		methodUml := uml.GetUmlForFunction(method)
		res += fmt.Sprintf("  %s", methodUml)
	}

	for _, method := range c.Methods.Funcs {
		for _, called := range method.Called.Funcs {
			if _, ok := c.Methods.Get(called.Name); ok && method != called {
				str := fmt.Sprintf("   %s -> %s\n", utils.NameToIdentifier(method.Name.String()), utils.NameToIdentifier(called.Name.String()))

				if _, ok := showed[str]; ok {
					continue
				}
				showed[str] = struct{}{}

				res += str
			}
		}
	}

	for _, field := range c.Fields.Fields {
		functions := make([]*symbols.Function, 0, len(field.Used))

		for used := range field.Used {
			functions = append(functions, used)
		}

		for i := 0; i < len(functions)-1; i++ {
			for j := i + 1; j < len(functions); j++ {
				str := fmt.Sprintf("   %s -> %s\n", utils.NameToIdentifier(functions[i].Name.String()), utils.NameToIdentifier(functions[j].Name.String()))

				if _, ok := showed[str]; ok {
					continue
				}
				showed[str] = struct{}{}

				res += str
			}
		}
	}

	return res
}
