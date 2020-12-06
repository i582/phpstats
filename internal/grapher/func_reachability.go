package grapher

import (
	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/grapher/templates"
	"github.com/i582/phpstats/internal/relations"
	"github.com/i582/phpstats/internal/utils"
)

func (g *Grapher) FunctionReachability(rel *relations.ReachabilityFunctionResult) string {
	graphName := "GraphFor_" + utils.NameToIdentifier(rel.ChildFunction.Name.String()+rel.ParentFunction.Name.String())
	functionGraph := &graph.Graph{
		Name:       graphName,
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label:      "Function " + utils.NormalizeSlashes(rel.ChildFunction.Name.String()) + " reachability from " + utils.NormalizeSlashes(rel.ParentFunction.Name.String()),
			Padding:    2.0,
			NodeMargin: 1.5,
		},
		NodeStyle: graph.NodeStyles{},
		EdgeStyle: templates.TemplateFunctionConnectionEdgeStyle(),
	}

	g.functionReachability(functionGraph, rel)
	ColorizeClassGraph(functionGraph)

	parentFunctionNode, _ := functionGraph.GetNode(templates.TemplateFunctionNode(rel.ParentFunction).Name)
	parentFunctionNode.Styles.FillColor = templates.FillColorLevel4
	parentFunctionNode.Styles.EdgeColor = templates.OutlineColorLevel4
	parentFunctionNode.ForceScale(2.3)

	childFunctionNode, _ := functionGraph.GetNode(templates.TemplateFunctionNode(rel.ChildFunction).Name)
	childFunctionNode.Styles.FillColor = templates.FillColorLevel4
	childFunctionNode.Styles.EdgeColor = templates.OutlineColorLevel4
	childFunctionNode.ForceScale(2.3)

	return functionGraph.String()
}

func (g *Grapher) functionReachability(functionGraph *graph.Graph, rel *relations.ReachabilityFunctionResult) {
	parentFunctionNode := templates.TemplateFunctionNode(rel.ParentFunction)
	functionGraph.AddNode(parentFunctionNode)

	childFunctionNode := templates.TemplateFunctionNode(rel.ChildFunction)
	functionGraph.AddNode(childFunctionNode)

	var prevNode *graph.Node
	for _, path := range rel.Paths {
		var pathContainsExcludeFunction bool
		for _, function := range path {
			if _, isExcluded := rel.ExcludedFunctions[function]; isExcluded {
				pathContainsExcludeFunction = true
				break
			}
		}
		if pathContainsExcludeFunction {
			continue
		}

		for _, function := range path {
			functionNode := templates.TemplateFunctionNode(function)
			functionNode, _ = functionGraph.AddNode(functionNode)

			if prevNode != nil {
				functionGraph.AddEdgeByNode(prevNode, functionNode, templates.TemplateFunctionConnectionEdgeStyle())
			}

			prevNode = functionNode
		}
		prevNode = nil
	}
}
