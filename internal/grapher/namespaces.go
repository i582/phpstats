package grapher

import (
	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/grapher/templates"
	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

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
