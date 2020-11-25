package graph

import (
	"fmt"
	"strings"
)

type SubGraphs map[string]*Graph

func (g SubGraphs) String() string {
	var res string

	for _, graph := range g {
		res += indentText(graph.String(), 1)
	}

	return res
}

type Graph struct {
	Name string

	IsSubgraph bool

	GraphStyle Styles
	NodeStyle  NodeStyles
	EdgeStyle  EdgeStyles

	SubGraphs SubGraphs
	Nodes     Nodes
	Edges     Edges
}

func (g *Graph) String() string {
	typ := "digraph"
	namePrefix := ""
	if g.IsSubgraph {
		typ = "subgraph"
		namePrefix = "cluster_"
	}

	return fmt.Sprintf(`%s %s%s {
	// graph styles
	%s
	// subgraphs
%s
	// node styles
	node %s
	// edge styles
	edge %s

	// nodes
%s
	// edges
%s
}
`, typ, namePrefix, g.Name, g.GraphStyle, g.SubGraphs, g.NodeStyle, g.EdgeStyle, g.Nodes, g.Edges)
}

func (g *Graph) AddNode(node *Node) (*Node, error) {
	if g.Nodes == nil {
		g.Nodes = make(Nodes)
	}

	foundNode, found := g.Nodes[node.Name]
	if found {
		return foundNode, fmt.Errorf("node with name '%s' already added", node.Name)
	}

	g.Nodes[node.Name] = node
	return node, nil
}

func (g *Graph) GetNode(name string) (*Node, bool) {
	if g.Nodes == nil {
		g.Nodes = make(Nodes)
		return nil, false
	}

	node, found := g.Nodes[name]
	return node, found
}

func (g *Graph) GetNodeInSubgraphs(name string) (*Node, bool) {
	if g.Nodes == nil {
		g.Nodes = make(Nodes)
	}

	node, found := g.Nodes[name]
	if !found {
		for _, graph := range g.SubGraphs {
			nodeInSubgraph, found := graph.GetNodeInSubgraphs(name)
			if found {
				return nodeInSubgraph, true
			}
		}
	}

	return node, found
}

func (g *Graph) AddEdge(from, to string, styles EdgeStyles) error {
	if g.Edges == nil {
		g.Edges = make(Edges)
	}

	edge, err := NewEdge(g, from, to, styles)
	if err != nil {
		return err
	}

	if _, found := g.Edges[from+to]; found {
		return fmt.Errorf("edge for '%s' and '%s' already added", from, to)
	}

	g.Edges[from+to] = edge
	return nil
}

func (g *Graph) AddEdgeByNode(from, to *Node, styles EdgeStyles) error {
	if g.Edges == nil {
		g.Edges = make(Edges)
	}

	edge := &Edge{
		From:   from,
		To:     to,
		Styles: styles,
	}

	if _, found := g.Edges[from.Name+to.Name]; found {
		return fmt.Errorf("edge for '%s' and '%s' already added", from.Name, to.Name)
	}

	from.addEdge(edge)
	to.addEdge(edge)

	g.Edges[from.Name+to.Name] = edge
	return nil
}

func (g *Graph) AddSubGraph(subgraph *Graph) *Graph {
	if g.SubGraphs == nil {
		g.SubGraphs = make(SubGraphs)
	}

	subgraph.IsSubgraph = true

	g.SubGraphs[subgraph.Name] = subgraph
	return subgraph
}

func (g *Graph) GetSubGraph(name string) (*Graph, bool) {
	if g.SubGraphs == nil {
		g.SubGraphs = make(SubGraphs)
		return nil, false
	}

	graph, found := g.SubGraphs[name]
	return graph, found
}

func (g *Graph) GetOrCreateSubGraph(nsName string) *Graph {
	nsName = strings.TrimPrefix(nsName, `\`)
	parts := strings.Split(nsName, `\`)
	return g.getOrCreateSubGraph(parts)
}

func (g *Graph) getOrCreateSubGraph(parts []string) *Graph {
	if len(parts) == 0 {
		return g
	}

	mainPart := parts[0]
	subGraph, found := g.GetSubGraph(mainPart)
	if !found {
		subGraph = &Graph{}
		*subGraph = *g
		subGraph.Edges = Edges{}
		subGraph.Nodes = Nodes{}
		subGraph.SubGraphs = SubGraphs{}

		subGraph.Name = mainPart

		g.AddSubGraph(subGraph)

		subGraph = subGraph.getOrCreateSubGraph(parts[1:])
	} else {
		subGraph = subGraph.getOrCreateSubGraph(parts[1:])
	}

	return subGraph
}
