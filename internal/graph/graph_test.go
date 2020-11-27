package graph

import (
	"fmt"
	"log"
	"testing"
)

func TestGraph(t *testing.T) {
	g := &Graph{
		Name:       "SomeTestGraph",
		IsSubgraph: false,
		GraphStyle: Styles{
			Label: "test graph",
		},
		NodeStyle: NodeStyles{
			Shape:      "rect",
			FillColor:  "#ff00ff",
			EdgeColor:  "#000000",
			Style:      "filled",
			FontSize:   12,
			FontFamily: "Times New Roman",
			FontColor:  "#000000",
		},
		EdgeStyle: EdgeStyles{
			ArrowTail: "empty",
		},
	}

	_, err := g.AddNode(&Node{
		Name: "TestNode",
		Styles: NodeStyles{
			Label:     "some test node",
			Shape:     "rect",
			FillColor: "#00ff00",
			FontSize:  20,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = g.AddNode(&Node{
		Name: "TestNode1",
		Styles: NodeStyles{
			Label:     "some  node",
			Shape:     "rect",
			FillColor: "#0000ff",
			FontSize:  15,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = g.AddEdge("TestNode", "TestNode1", EdgeStyles{Style: "dashed"})
	if err != nil {
		log.Fatal(err)
	}

	subGraph := &Graph{
		Name:       "SomeTestSubGraph",
		IsSubgraph: false,
		GraphStyle: Styles{
			Label: "test sub graph",
		},
		NodeStyle: NodeStyles{
			Shape:      "ellipse",
			FillColor:  "#0000ff",
			Style:      "filled",
			FontSize:   15,
			FontFamily: "Times New Roman",
			FontColor:  "#000000",
		},
		EdgeStyle: EdgeStyles{
			ArrowTail: "empty",
		},
	}
	_, err = subGraph.AddNode(&Node{
		Name: "TestSubNode1",
		Styles: NodeStyles{
			Label:     "some sub node",
			Shape:     "ellipse",
			FillColor: "#00ffff",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	g.AddSubGraph(subGraph)

	fmt.Print(g)
}
