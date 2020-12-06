package graph

import (
	"fmt"
)

type Nodes map[string]*Node

func (n Nodes) String() string {
	var res string

	for _, node := range n {
		res += node.String() + "\n"
	}

	return res
}

func (n Nodes) DeleteNode(node *Node) {
	delete(n, node.Name)
}

func (n Nodes) AddNode(node *Node) {
	n[node.Name] = node
}

type Node struct {
	Name string

	Edges  Edges
	Styles NodeStyles
}

func (n *Node) String() string {
	var res string

	res += fmt.Sprintf("\t%s %s", n.Name, n.Styles)

	return res
}

func (n *Node) Scale(val float64) {
	fontSize := n.Styles.FontSize
	if fontSize == 0 {
		fontSize = 12
	}
	if fontSize != 12 {
		return
	}

	fontSize = int64(float64(fontSize) * val)
	n.Styles.FontSize = fontSize
}

func (n *Node) ForceScale(val float64) {
	fontSize := n.Styles.FontSize
	if fontSize == 0 {
		fontSize = 12
	}

	fontSize = int64(float64(12) * val)
	n.Styles.FontSize = fontSize
}

func (n *Node) addEdge(edge *Edge) {
	if n.Edges == nil {
		n.Edges = make(Edges)
	}

	n.Edges[edge.From.Name+edge.To.Name] = edge
}
