package grapher

import (
	"fmt"

	"github.com/i582/phpstats/internal/grapher/unl"
	"github.com/i582/phpstats/internal/stats"
)

func (g *Grapher) FileDeps(f *stats.File, maxRecursion int64, root, block bool) string {
	var res string

	res += graphHeader

	main := g.fileDepsRecursive(f, 0, maxRecursion, root, block, visitedMap{})

	res += main

	return g.graphWrapper(res, f.UniqueId())
}

func (g *Grapher) fileDepsRecursive(f *stats.File, levelRecursion, maxRecursion int64, root, block bool, visited visitedMap) string {
	var res string

	classUml := uml.GetUmlForFile(f)
	umlGraph := g.outputWithColor("   "+classUml, defaultColor, defaultColor)

	if _, ok := visited[umlGraph]; !ok {
		res += umlGraph
		visited[umlGraph] = struct{}{}
	}

	if levelRecursion > maxRecursion {
		return res
	}

	if root {
		for _, file := range f.RequiredRoot.Files {
			str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", file.UniqueId(), f.UniqueId())
			if _, ok := visited[str]; !ok {
				res += str
				visited[str] = struct{}{}
			}

			res += g.fileDepsRecursive(file, levelRecursion+1, maxRecursion, block, root, visited)
		}
	}

	if block {
		for _, file := range f.RequiredBlock.Files {
			str := fmt.Sprintf("   \"%s\" -> \"%s\" [style=\"dashed\"]\n", file.UniqueId(), f.UniqueId())
			if _, ok := visited[str]; !ok {
				res += str
				visited[str] = struct{}{}
			}

			res += g.fileDepsRecursive(file, levelRecursion+1, maxRecursion, block, root, visited)
		}
	}

	return res
}
