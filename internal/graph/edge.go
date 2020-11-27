package graph

import (
	"fmt"
	"strings"
)

type Edges map[string]*Edge

func (n Edges) String() string {
	var res string

	for _, edge := range n {
		res += edge.String() + "\n"
	}

	res = strings.TrimSuffix(res, "\n")

	return res
}

type Edge struct {
	From *Node
	To   *Node

	Styles EdgeStyles
}

func NewEdge(g *Graph, from, to string, styles EdgeStyles) (*Edge, error) {
	fromNode, found := g.GetNode(from)
	if !found {
		return nil, fmt.Errorf("node '%s' not found in graph '%s'", from, g.Name)
	}
	toNode, found := g.GetNode(to)
	if !found {
		return nil, fmt.Errorf("node '%s' not found in graph '%s'", to, g.Name)
	}

	edge := &Edge{
		From:   fromNode,
		To:     toNode,
		Styles: styles,
	}

	fromNode.addEdge(edge)
	toNode.addEdge(edge)

	return edge, nil
}

func (n Edge) String() string {
	var res string

	res += fmt.Sprintf("\t%s -> %s %s", n.From.Name, n.To.Name, n.Styles)

	return res
}
