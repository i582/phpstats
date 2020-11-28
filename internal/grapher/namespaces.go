package grapher

import (
	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/grapher/templates"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func ColorizeNamespacesDepsGraph(classGraph *graph.Graph) {
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

func (g *Grapher) NamespacesDeps(n *symbols.Namespace, maxRecursion int64) string {
	graphName := "GraphFor_" + utils.NameToIdentifier(n.FullName)
	namespaceGraph := &graph.Graph{
		Name:       graphName,
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label:      "Namespace " + utils.NormalizeSlashes(n.FullName) + " dependencies",
			Padding:    2.0,
			NodeMargin: 1.5,
		},
		NodeStyle: graph.NodeStyles{},
		EdgeStyle: templates.TemplateFunctionConnectionEdgeStyle(),
	}

	g.namespacesDepsRecursive(namespaceGraph, n, 0, maxRecursion)

	ColorizeNamespacesDepsGraph(namespaceGraph)

	return namespaceGraph.String()
}

func (g *Grapher) namespacesDepsRecursive(namespaceGraph *graph.Graph, n *symbols.Namespace, levelRecursion, maxRecursion int64) {
	nsNode := templates.TemplateNamespaceNode(n)
	namespaceGraph.AddNode(nsNode)

	if levelRecursion > maxRecursion {
		return
	}

	for _, class := range n.Classes.Classes {
		for _, depClass := range class.Deps.Classes {
			if depClass.Namespace == nil {
				globalNamespaceName := "GlobalNamespace"
				globalNamespaceNode, found := namespaceGraph.GetNode(globalNamespaceName)
				if !found {
					globalNamespaceNode, _ = namespaceGraph.AddNode(templates.TemplateNamespaceNode(&symbols.Namespace{
						FullName: globalNamespaceName,
					}))
				}
				namespaceGraph.AddEdgeByNode(nsNode, globalNamespaceNode, templates.TemplateNamespaceConnectionEdgeStyle())
				continue
			}

			depNamespace := depClass.Namespace
			if n == depNamespace {
				continue
			}

			depNamespaceNode := templates.TemplateNamespaceNode(depNamespace)
			depNamespaceNode, _ = namespaceGraph.AddNode(depNamespaceNode)
			namespaceGraph.AddEdgeByNode(nsNode, depNamespaceNode, templates.TemplateNamespaceConnectionEdgeStyle())

			g.namespacesDepsRecursive(namespaceGraph, depNamespace, levelRecursion+1, maxRecursion)
		}
	}

	for _, fun := range n.Functions.Funcs {
		for _, function := range fun.Called.Funcs {
			if function.IsMethod() {
				continue
			}

			if function.Namespace == nil {
				globalNamespaceName := "GlobalNamespace"
				globalNamespaceNode, found := namespaceGraph.GetNode(globalNamespaceName)
				if !found {
					globalNamespaceNode, _ = namespaceGraph.AddNode(templates.TemplateNamespaceNode(&symbols.Namespace{
						FullName: globalNamespaceName,
					}))
				}
				namespaceGraph.AddEdgeByNode(nsNode, globalNamespaceNode, templates.TemplateNamespaceConnectionEdgeStyle())
				continue
			}

			depNamespace := function.Namespace
			if n == depNamespace {
				continue
			}

			depNamespaceNode := templates.TemplateNamespaceNode(depNamespace)
			depNamespaceNode, _ = namespaceGraph.AddNode(depNamespaceNode)
			namespaceGraph.AddEdgeByNode(nsNode, depNamespaceNode, templates.TemplateNamespaceConnectionEdgeStyle())

			g.namespacesDepsRecursive(namespaceGraph, depNamespace, levelRecursion+1, maxRecursion)
		}

		for _, depClass := range fun.Deps().Classes {
			if depClass.Namespace == nil {
				globalNamespaceName := "GlobalNamespace"
				globalNamespaceNode, found := namespaceGraph.GetNode(globalNamespaceName)
				if !found {
					globalNamespaceNode, _ = namespaceGraph.AddNode(templates.TemplateNamespaceNode(&symbols.Namespace{
						FullName: globalNamespaceName,
					}))
				}
				namespaceGraph.AddEdgeByNode(nsNode, globalNamespaceNode, templates.TemplateNamespaceConnectionEdgeStyle())
				continue
			}

			depNamespace := depClass.Namespace
			if n == depNamespace {
				continue
			}

			depNamespaceNode := templates.TemplateNamespaceNode(depNamespace)
			depNamespaceNode, _ = namespaceGraph.AddNode(depNamespaceNode)
			namespaceGraph.AddEdgeByNode(nsNode, depNamespaceNode, templates.TemplateNamespaceConnectionEdgeStyle())

			g.namespacesDepsRecursive(namespaceGraph, depNamespace, levelRecursion+1, maxRecursion)
		}
	}
}

func (g *Grapher) Namespaces(n *symbols.Namespace, maxRecursion int64) string {
	graphName := "GraphFor_" + utils.NameToIdentifier(n.FullName)
	namespaceGraph := &graph.Graph{
		Name:       graphName,
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label:      "Namespaces " + utils.NormalizeSlashes(n.FullName),
			Padding:    2.0,
			NodeMargin: 1.5,
		},
		NodeStyle: graph.NodeStyles{},
		EdgeStyle: templates.TemplateFunctionConnectionEdgeStyle(),
	}

	g.namespacesRecursive(namespaceGraph, n, 0, maxRecursion)

	funcNode, found := namespaceGraph.GetNodeInSubgraphs(utils.NameToIdentifier(n.FullName))
	if found {
		funcNode.Styles.FillColor = templates.FillColorLevel3
		funcNode.Styles.EdgeColor = templates.OutlineColorLevel3
		funcNode.Scale(1.7)
	}

	return namespaceGraph.String()
}

func (g *Grapher) namespacesRecursive(namespaceGraph *graph.Graph, n *symbols.Namespace, levelRecursion, maxRecursion int64) {
	nsNode := templates.TemplateNamespaceNode(n)
	namespaceGraph.AddNode(nsNode)

	if levelRecursion > maxRecursion {
		return
	}

	for _, namespace := range n.Childs.Namespaces {
		childNamespaceNode := templates.TemplateNamespaceNode(namespace)
		childNamespaceNode, _ = namespaceGraph.AddNode(childNamespaceNode)

		namespaceGraph.AddEdgeByNode(nsNode, childNamespaceNode, templates.TemplateNamespaceConnectionEdgeStyle())

		g.namespacesRecursive(namespaceGraph, namespace, levelRecursion+1, maxRecursion)
	}
}
