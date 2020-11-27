package grapher

import (
	"fmt"

	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/grapher/templates"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/stats/walkers"
	"github.com/i582/phpstats/internal/utils"
)

func (g *Grapher) ClassSuperGlobalsDeps(c *symbols.Class) string {
	graphName := "GraphFor_" + utils.NameToIdentifier(c.Name)
	classGraph := &graph.Graph{
		Name:       graphName,
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label:   "Class " + utils.NormalizeSlashes(c.Name) + " super globals constant.",
			Padding: 2.0,
		},
		NodeStyle: graph.NodeStyles{},
		EdgeStyle: graph.EdgeStyles{},
	}

	g.classSuperGlobalsDepsRecursive(classGraph, c)

	mainClassNode, found := classGraph.GetNode(utils.NameToIdentifier(c.Name))
	if found {
		mainClassNode.Scale(1.4)
		templates.ColorizeByScale(mainClassNode, 1.4)
	}

	return classGraph.String()
}

func (g *Grapher) classSuperGlobalsDepsRecursive(classGraph *graph.Graph, c *symbols.Class) {
	classNode := templates.TemplateClassNode(c)
	classGraph.AddNode(classNode)

	for _, constant := range c.UsedConstants.Constants {
		if !constant.IsEmbedded() {
			constantNode := templates.TemplateConstantNode(constant)
			classGraph.AddNode(constantNode)
			classGraph.AddEdgeByNode(classNode, constantNode, templates.TemplateClassConnectionEdgeStyle())
		}
	}
}

func (g *Grapher) ClassImplementsExtendsDeps(c *symbols.Class, maxRecursion int64) string {
	graphName := "GraphFor_" + utils.NameToIdentifier(c.Name)
	classGraph := &graph.Graph{
		Name:       graphName,
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label:   "Class " + utils.NormalizeSlashes(c.Name) + " inheritance graph.",
			Padding: 2.0,
		},
		NodeStyle: graph.NodeStyles{},
		EdgeStyle: graph.EdgeStyles{},
	}

	g.classImplementsExtendsDepsRecursive(classGraph, c, 0, maxRecursion, map[*symbols.Class]struct{}{})

	mainClassNode, found := classGraph.GetNode(utils.NameToIdentifier(c.Name))
	if found {
		mainClassNode.Scale(1.4)
		templates.ColorizeByScale(mainClassNode, 1.4)
	}

	return classGraph.String()
}

func (g *Grapher) classImplementsExtendsDepsRecursive(classGraph *graph.Graph, c *symbols.Class, levelRecursion, maxRecursion int64, visitedClasses map[*symbols.Class]struct{}) {
	classNode := templates.TemplateClassNode(c)
	classGraph.AddNode(classNode)

	if levelRecursion > maxRecursion {
		return
	}

	for _, implementClass := range c.Implements.Classes {
		implementClassNode := templates.TemplateClassNode(implementClass)
		implementClassNode, _ = classGraph.AddNode(implementClassNode)

		classGraph.AddEdgeByNode(classNode, implementClassNode, templates.TemplateImplementEdgeStyle())

		if _, found := visitedClasses[implementClass]; !found {
			visitedClasses[implementClass] = struct{}{}
			g.classImplementsExtendsDepsRecursive(classGraph, implementClass, levelRecursion+1, maxRecursion, visitedClasses)
		}
	}

	for _, extendedClass := range c.Extends.Classes {
		extendedClassNode := templates.TemplateClassNode(extendedClass)
		extendedClassNode, _ = classGraph.AddNode(extendedClassNode)

		classGraph.AddEdgeByNode(classNode, extendedClassNode, templates.TemplateExtendEdgeStyle())

		if _, found := visitedClasses[extendedClass]; !found {
			visitedClasses[extendedClass] = struct{}{}
			g.classImplementsExtendsDepsRecursive(classGraph, extendedClass, levelRecursion+1, maxRecursion, visitedClasses)
		}
	}

	for _, implementClass := range c.ImplementsBy.Classes {
		implementClassNode := templates.TemplateClassNode(implementClass)
		implementClassNode, _ = classGraph.AddNode(implementClassNode)

		classGraph.AddEdgeByNode(implementClassNode, classNode, templates.TemplateImplementEdgeStyle())

		if _, found := visitedClasses[implementClass]; !found {
			visitedClasses[implementClass] = struct{}{}
			g.classImplementsExtendsDepsRecursive(classGraph, implementClass, levelRecursion+1, maxRecursion, visitedClasses)
		}
	}

	for _, extendedClass := range c.ExtendsBy.Classes {
		extendedClassNode := templates.TemplateClassNode(extendedClass)
		extendedClassNode, _ = classGraph.AddNode(extendedClassNode)

		classGraph.AddEdgeByNode(extendedClassNode, classNode, templates.TemplateExtendEdgeStyle())

		if _, found := visitedClasses[extendedClass]; !found {
			visitedClasses[extendedClass] = struct{}{}
			g.classImplementsExtendsDepsRecursive(classGraph, extendedClass, levelRecursion+1, maxRecursion, visitedClasses)
		}
	}
}

func ColorizeClassGraph(classGraph *graph.Graph) {
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

		node.Scale(scaleForNode)
		templates.ColorizeByScale(node, scaleForNode)
	}
}

func (g *Grapher) ClassDeps(c *symbols.Class, maxRecursion int64, withGroups bool) string {
	graphName := "GraphFor_" + utils.NameToIdentifier(c.Name)
	classGraph := &graph.Graph{
		Name:       graphName,
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label:   "Class " + utils.NormalizeSlashes(c.Name) + " dependencies",
			Padding: 2.0,
		},
		NodeStyle: graph.NodeStyles{},
		EdgeStyle: graph.EdgeStyles{},
	}

	g.classDepsRecursive(classGraph, c, 0, maxRecursion, withGroups)
	ColorizeClassGraph(classGraph)

	return classGraph.String()
}

func (g *Grapher) classDepsRecursive(classGraph *graph.Graph, c *symbols.Class, levelRecursion, maxRecursion int64, withGroups bool) {
	classNode := templates.TemplateClassNode(c)

	if withGroups {
		pack, found := walkers.GlobalCtx.Packages.GetPackage(c.Name)
		if found {
			subgraphName := utils.NameToIdentifier(pack.Name)
			subgraph, found := classGraph.GetSubGraph(subgraphName)
			if !found {
				subgraph = classGraph.AddSubGraph(&graph.Graph{
					Name:       subgraphName,
					GraphStyle: classGraph.GraphStyle,
					NodeStyle:  classGraph.NodeStyle,
					EdgeStyle:  classGraph.EdgeStyle,
				})
				subgraph.GraphStyle.Label = "Package " + pack.Name
				subgraph.GraphStyle.GraphMargin = 75
				subgraph.GraphStyle.Style = "filled"
				subgraph.GraphStyle.BorderColor = templates.DefaultSubgraphFillColor
			}

			subgraph.AddNode(classNode)
		} else {
			subgraphName := "GlobalPackage"
			subgraph, found := classGraph.GetSubGraph(subgraphName)
			if !found {
				subgraph = classGraph.AddSubGraph(&graph.Graph{
					Name:       subgraphName,
					GraphStyle: classGraph.GraphStyle,
					NodeStyle:  classGraph.NodeStyle,
					EdgeStyle:  classGraph.EdgeStyle,
				})
				subgraph.GraphStyle.Label = "Global Package"
				subgraph.GraphStyle.GraphMargin = 75
				subgraph.GraphStyle.Style = "filled"
				subgraph.GraphStyle.BorderColor = templates.DefaultSubgraphFillColor
			}

			subgraph.AddNode(classNode)
		}
	} else {
		classGraph.AddNode(classNode)
	}

	if levelRecursion > maxRecursion {
		return
	}

	for _, implementClass := range c.Implements.Classes {
		implementClassNode := templates.TemplateClassNode(implementClass)
		implementClassNode, _ = classGraph.AddNode(implementClassNode)

		classGraph.AddEdgeByNode(classNode, implementClassNode, templates.TemplateImplementEdgeStyle())

		g.classDepsRecursive(classGraph, implementClass, levelRecursion+1, maxRecursion, withGroups)
	}

	for _, extendedClass := range c.Extends.Classes {
		extendedClassNode := templates.TemplateClassNode(extendedClass)
		extendedClassNode, _ = classGraph.AddNode(extendedClassNode)

		classGraph.AddEdgeByNode(classNode, extendedClassNode, templates.TemplateExtendEdgeStyle())

		g.classDepsRecursive(classGraph, extendedClass, levelRecursion+1, maxRecursion, withGroups)
	}

	for _, depsClass := range c.Deps.Classes {
		depsClassNode := templates.TemplateClassNode(depsClass)
		depsClassNode, _ = classGraph.AddNode(depsClassNode)

		classGraph.AddEdgeByNode(classNode, depsClassNode, templates.TemplateClassConnectionEdgeStyle())

		g.classDepsRecursive(classGraph, depsClass, levelRecursion+1, maxRecursion, withGroups)
	}

}

func AddedCountLinksInGraphNodes(classGraph *graph.Graph) {
	for _, node := range classGraph.Nodes {
		node.Styles.Label += "\\n(links: " + fmt.Sprint(len(node.Edges)) + ")"
	}
}

func (g *Grapher) Lcom4(c *symbols.Class) string {
	graphName := "GraphFor_" + utils.NameToIdentifier(c.Name)
	classGraph := &graph.Graph{
		Name:       graphName,
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label:      "Lack of Cohesion in Methods 4 (LCOM4) graph for " + utils.NormalizeSlashes(c.Name),
			Padding:    2.5,
			NodeMargin: 1.5,
		},
		NodeStyle: graph.NodeStyles{},
		EdgeStyle: graph.EdgeStyles{},
	}

	g.lcom4(classGraph, c)

	AddedCountLinksInGraphNodes(classGraph)

	return classGraph.String()
}

func (g *Grapher) lcom4(classGraph *graph.Graph, c *symbols.Class) {
	for _, method := range c.Methods.Funcs {
		methodNode := templates.TemplateFunctionNode(method)
		classGraph.AddNode(methodNode)

		methodNode.Styles.FillColor = templates.FillColorLevel3
		methodNode.Styles.EdgeColor = templates.OutlineColorLevel3
		methodNode.Scale(1.5)

		for _, calledFunction := range method.Called.Funcs {
			if calledFunction.Class == c {
				calledFunctionNode, _ := classGraph.AddNode(templates.TemplateFunctionNode(calledFunction))

				calledFunctionNode.Styles.FillColor = templates.FillColorLevel3
				calledFunctionNode.Styles.EdgeColor = templates.OutlineColorLevel3
				calledFunctionNode.Scale(1.5)

				classGraph.AddEdgeByNode(methodNode, calledFunctionNode, graph.EdgeStyles{Color: templates.OutlineColorLevel3})
			}
		}
	}

	for _, field := range c.Fields.Fields {
		fieldNode, _ := classGraph.AddNode(templates.TemplateFieldNode(field))

		for _, function := range field.Used.Funcs {
			if function.Class == c {
				functionNode, _ := classGraph.AddNode(templates.TemplateFunctionNode(function))

				functionNode.Styles.FillColor = templates.FillColorLevel3
				functionNode.Styles.EdgeColor = templates.OutlineColorLevel3
				functionNode.Scale(1.5)

				classGraph.AddEdgeByNode(functionNode, fieldNode, graph.EdgeStyles{Color: templates.OutlineColorLevel1, Style: "dashed"})
			}
		}
	}

	for _, constant := range c.Constants.Constants {
		constantNode, _ := classGraph.AddNode(templates.TemplateConstantNode(constant))

		for _, function := range constant.Used.Funcs {
			if function.Class == c {
				functionNode, _ := classGraph.AddNode(templates.TemplateFunctionNode(function))

				functionNode.Styles.FillColor = templates.FillColorLevel3
				functionNode.Styles.EdgeColor = templates.OutlineColorLevel3
				functionNode.Scale(1.5)

				classGraph.AddEdgeByNode(functionNode, constantNode, graph.EdgeStyles{Color: templates.OutlineColorLevel2, Style: "dotted"})
			}
		}
	}

	for _, node := range classGraph.Nodes {
		if len(node.Edges) == 0 {
			subgraph, found := classGraph.GetSubGraph("WithoutConnections")
			if !found {
				subgraph = classGraph.AddSubGraph(&graph.Graph{
					Name:       "WithoutConnections",
					GraphStyle: classGraph.GraphStyle,
					NodeStyle:  classGraph.NodeStyle,
					EdgeStyle:  classGraph.EdgeStyle,
				})
				subgraph.GraphStyle.Label = "Symbols without connections"
			}
			subgraph.AddNode(node)
			classGraph.Nodes.DeleteNode(node)
		}
	}
}
