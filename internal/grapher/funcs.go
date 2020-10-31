package grapher

import (
	"fmt"
	"strings"

	"github.com/i582/phpstats/internal/grapher/unl"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func (g *Grapher) FuncDeps(f *symbols.Function) string {
	var res string

	res += graphHeader

	funcUml := uml.GetUmlForFunction(f)
	umlGraph := g.outputWithColor("   "+funcUml, g.getColorForFunction(f), defaultColor)

	definitions, connections := g.funcDeps(f)

	definitions[umlGraph] = struct{}{}

	deps := f.Deps()
	depsBy := f.DepsBy()

	shovedDef := map[string]struct{}{}

	for _, dep := range deps.Classes {
		name := utils.NameToIdentifier(dep.Name)

		var depSubGraph string

		for def := range definitions {
			if strings.Contains(def, name) {
				depSubGraph += def
				shovedDef[def] = struct{}{}
			}
		}

		res += g.subGraphWrapper(depSubGraph, utils.NormalizeSlashes(dep.Name))
	}

	for _, dep := range depsBy.Classes {
		name := utils.NameToIdentifier(dep.Name)

		var depSubGraph string

		for def := range definitions {
			if strings.Contains(def, name) {
				depSubGraph += def
				shovedDef[def] = struct{}{}
			}
		}

		res += g.subGraphWrapper(depSubGraph, utils.NormalizeSlashes(dep.Name))
	}

	if f.Class != nil {
		name := utils.NameToIdentifier(f.Class.Name)

		var depSubGraph string

		for def := range definitions {
			if strings.Contains(def, name) {
				depSubGraph += def
				shovedDef[def] = struct{}{}
			}
		}

		res += g.subGraphWrapperColor(depSubGraph, utils.NormalizeSlashes(f.Class.Name), "#bbbbbb")
	}

	for def := range definitions {
		if _, ok := shovedDef[def]; ok {
			continue
		}

		res += def
	}

	for con := range connections {
		res += con
	}

	return g.graphWrapper(res, utils.NameToIdentifier(f.Name.String()))
}

func (g *Grapher) funcDeps(f *symbols.Function) (map[string]struct{}, map[string]struct{}) {
	definitions := make(map[string]struct{}, f.Called.Len()+f.CalledBy.Len())
	connections := make(map[string]struct{}, f.Called.Len()+f.CalledBy.Len())

	for _, called := range f.Called.Funcs {
		str := fmt.Sprintf("   %s -> %s\n", utils.NameToIdentifier(f.Name.String()), utils.NameToIdentifier(called.Name.String()))
		connections[str] = struct{}{}

		funcUml := uml.GetUmlForFunction(called)
		colorFuncUml := g.outputWithColor("   "+funcUml, g.getColorForFunction(called), defaultColor)
		definitions[colorFuncUml] = struct{}{}
	}

	for _, calledBy := range f.CalledBy.Funcs {
		str := fmt.Sprintf("   %s -> %s\n", utils.NameToIdentifier(calledBy.Name.String()), utils.NameToIdentifier(f.Name.String()))
		connections[str] = struct{}{}

		funcUml := uml.GetUmlForFunction(calledBy)
		colorFuncUml := g.outputWithColor("   "+funcUml, g.getColorForFunction(calledBy), defaultColor)
		definitions[colorFuncUml] = struct{}{}
	}

	return definitions, connections
}

func (g *Grapher) getColorForFunction(c *symbols.Function) string {
	if c.Class != nil {
		return "#bbbbbb"
	}

	return defaultColor
}
