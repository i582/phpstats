package grapher

import (
	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/grapher/templates"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func (g *Grapher) NewFuncDeps(f *symbols.Function, maxRecursion int64) string {
	graphName := "GraphFor_" + utils.NameToIdentifier(f.Name.String())
	functionGraph := &graph.Graph{
		Name:       graphName,
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label:      "Function " + utils.NormalizeSlashes(f.Name.String()) + " dependencies",
			Padding:    2.0,
			NodeMargin: 1.5,
		},
		NodeStyle: graph.NodeStyles{},
		EdgeStyle: templates.TemplateFunctionConnectionEdgeStyle(),
	}

	g.funcDepsRecursive(functionGraph, f, 0, maxRecursion)

	funcNode, found := functionGraph.GetNodeInSubgraphs(utils.NameToIdentifier(f.Name.String()))
	if found {
		funcNode.Styles.FillColor = templates.FillColorLevel3
		funcNode.Styles.EdgeColor = templates.OutlineColorLevel3
		funcNode.Scale(1.7)
	}

	return functionGraph.String()
}

func (g *Grapher) funcDepsRecursive(functionGraph *graph.Graph, f *symbols.Function, levelRecursion, maxRecursion int64) {
	mainFunctionSubGraph := g.createSubGraphForFunctionClass(f, functionGraph)
	mainFuncNode, _ := mainFunctionSubGraph.AddNode(templates.TemplateFunctionNode(f))

	if levelRecursion > maxRecursion {
		return
	}

	for _, function := range f.Called.Funcs {
		subGraph := g.createSubGraphForFunctionClass(function, functionGraph)
		subGraph.AddNode(templates.TemplateFunctionNode(function))

		g.funcDepsRecursive(functionGraph, function, levelRecursion+1, maxRecursion)
	}

	for _, function := range f.CalledBy.Funcs {
		subGraph := g.createSubGraphForFunctionClass(function, functionGraph)
		subGraph.AddNode(templates.TemplateFunctionNode(function))

		g.funcDepsRecursive(functionGraph, function, levelRecursion+1, maxRecursion)
	}

	for _, field := range f.UsedFields.Fields {
		subGraph := g.createSubGraphForFieldClass(field, functionGraph)
		subGraph.AddNode(templates.TemplateFieldNode(field))
	}

	for _, function := range f.Called.Funcs {
		funcNode, found := functionGraph.GetNodeInSubgraphs(utils.NameToIdentifier(function.Name.String()))
		if !found {
			continue
		}

		functionGraph.AddEdgeByNode(mainFuncNode, funcNode, graph.EdgeStyles{Color: templates.OutlineColorLevel2})
	}

	for _, function := range f.CalledBy.Funcs {
		funcNode, found := functionGraph.GetNodeInSubgraphs(utils.NameToIdentifier(function.Name.String()))
		if !found {
			continue
		}

		functionGraph.AddEdgeByNode(funcNode, mainFuncNode, graph.EdgeStyles{})
	}

	for _, field := range f.UsedFields.Fields {
		fieldNode, found := functionGraph.GetNodeInSubgraphs(utils.NameToIdentifier(field.String()))
		if !found {
			continue
		}

		functionGraph.AddEdgeByNode(mainFuncNode, fieldNode, graph.EdgeStyles{})
	}
}

func (g *Grapher) createSubGraphForFunctionClass(function *symbols.Function, functionGraph *graph.Graph) *graph.Graph {
	return g.createSubGraphForClass(function.Class, functionGraph)
}

func (g *Grapher) createSubGraphForFieldClass(field *symbols.Field, functionGraph *graph.Graph) *graph.Graph {
	return g.createSubGraphForClass(field.Class, functionGraph)
}

func (g *Grapher) createSubGraphForClass(class *symbols.Class, functionGraph *graph.Graph) *graph.Graph {
	var subGraph *graph.Graph
	var found bool

	if class != nil {
		subGraphName := utils.NameToIdentifier(class.Name)
		subGraph, found = functionGraph.GetSubGraph(subGraphName)
		if !found {
			subGraph = functionGraph.AddSubGraph(&graph.Graph{
				Name:       subGraphName,
				IsSubgraph: false,
				GraphStyle: graph.Styles{
					Label:       utils.NormalizeSlashes(class.Name),
					BorderColor: templates.DefaultOutlineColor,
					FontColor:   templates.DefaultOutlineColor,
				},
			})
		}
	} else {
		subGraph, found = functionGraph.GetSubGraph("globalScope")
		if !found {
			subGraph = functionGraph.AddSubGraph(&graph.Graph{
				Name:       "globalScope",
				IsSubgraph: false,
				GraphStyle: graph.Styles{
					Label:       "Global Scope",
					BorderColor: templates.DefaultOutlineColor,
					FontColor:   templates.DefaultOutlineColor,
				},
			})
		}
	}
	return subGraph
}
