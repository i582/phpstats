package grapher

import (
	"fmt"

	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/grapher/templates"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func ColorizeFileGraph(classGraph *graph.Graph) {
	maxConnections := 0
	minConnections := 10000

	for _, node := range classGraph.Nodes {
		connections := len(node.Edges)

		if connections < minConnections {
			minConnections = connections
		}

		if connections > maxConnections {
			maxConnections = connections
		}
	}

	if maxConnections == 0 {
		maxConnections = 1
	}

	maxScale := 1.8
	minScale := 0.8
	diffScale := maxScale - minScale

	for _, node := range classGraph.Nodes {
		scaleForNode := minScale + diffScale*(float64(len(node.Edges))/float64(maxConnections))

		if scaleForNode > 1 {
			fmt.Print()
		}

		node.Scale(scaleForNode)
		templates.ColorizeByScale(node, scaleForNode)
	}
}

func (g *Grapher) FileDeps(f *symbols.File, maxRecursion int64, root, block bool) string {
	graphName := "GraphFor_" + utils.NameToIdentifier(f.Path)
	classGraph := &graph.Graph{
		Name:       graphName,
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label:      "File " + utils.NormalizeSlashes(f.Name) + " included files.",
			Padding:    2.0,
			NodeMargin: 3,
		},
		NodeStyle: graph.NodeStyles{},
		EdgeStyle: graph.EdgeStyles{},
	}

	if !root && !block {
		root = true
		block = true
	}

	g.fileDepsRecursive(classGraph, f, 0, maxRecursion, root, block, map[string]struct{}{})

	ColorizeFileGraph(classGraph)

	return classGraph.String()
}

func (g *Grapher) fileDepsRecursive(classGraph *graph.Graph, f *symbols.File, levelRecursion, maxRecursion int64, root, block bool, visitedFiles map[string]struct{}) {
	fileNode := templates.TemplateFileNode(f)
	fileNode, _ = classGraph.AddNode(fileNode)

	if levelRecursion > maxRecursion {
		return
	}

	if root {
		for _, rootFile := range f.RequiredRoot.Files {
			rootFileNode := templates.TemplateFileNode(rootFile)
			rootFileNode, _ = classGraph.AddNode(rootFileNode)

			classGraph.AddEdgeByNode(fileNode, rootFileNode, templates.TemplateRootFileEdgeStyle())

			if _, found := visitedFiles[rootFile.Path+f.Path]; !found {
				visitedFiles[rootFile.Path+f.Path] = struct{}{}
				g.fileDepsRecursive(classGraph, rootFile, levelRecursion+1, maxRecursion, root, block, visitedFiles)
			}
		}
	}

	if block {
		for _, blockFile := range f.RequiredBlock.Files {
			blockFileNode := templates.TemplateFileNode(blockFile)
			blockFileNode, _ = classGraph.AddNode(blockFileNode)

			classGraph.AddEdgeByNode(fileNode, blockFileNode, templates.TemplateBlockFileEdgeStyle())

			if _, found := visitedFiles[blockFile.Path+f.Path]; !found {
				visitedFiles[blockFile.Path+f.Path] = struct{}{}
				g.fileDepsRecursive(classGraph, blockFile, levelRecursion+1, maxRecursion, root, block, visitedFiles)
			}
		}
	}
}
