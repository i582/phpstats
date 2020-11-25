package templates

import (
	"fmt"
	"log"
	"testing"

	"github.com/i582/phpstats/internal/graph"
	"github.com/i582/phpstats/internal/stats/symbols"
)

func TestTemplateClassNode(t *testing.T) {
	g := &graph.Graph{
		Name:       "SomeTestGraph",
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label: "test graph",
		},
		NodeStyle: graph.NodeStyles{
			Shape:      "rect",
			FillColor:  "#ff00ff",
			EdgeColor:  "#000000",
			Style:      "filled",
			FontSize:   12,
			FontFamily: "Times New Roman",
			FontColor:  "#000000",
		},
		EdgeStyle: graph.EdgeStyles{
			ArrowTail: "empty",
		},
	}

	_, err := g.AddNode(&graph.Node{
		Name: "TestNode",
		Styles: graph.NodeStyles{
			Label:     "some test node",
			Shape:     "rect",
			FillColor: "#00ff00",
			FontSize:  20,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = g.AddNode(&graph.Node{
		Name: "TestNode1",
		Styles: graph.NodeStyles{
			Label:     "some  node",
			Shape:     "rect",
			FillColor: "#0000ff",
			FontSize:  15,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = g.AddEdge("TestNode", "TestNode1", graph.EdgeStyles{Style: "dashed"})
	if err != nil {
		log.Fatal(err)
	}

	subGraph := &graph.Graph{
		Name:       "SomeTestSubGraph",
		IsSubgraph: false,
		GraphStyle: graph.Styles{
			Label: "test sub graph",
		},
		NodeStyle: graph.NodeStyles{
			Shape:      "ellipse",
			FillColor:  "#0000ff",
			Style:      "filled",
			FontSize:   15,
			FontFamily: "Times New Roman",
			FontColor:  "#000000",
		},
		EdgeStyle: graph.EdgeStyles{
			ArrowTail: "empty",
		},
	}
	_, err = subGraph.AddNode(&graph.Node{
		Name: "TestSubNode1",
		Styles: graph.NodeStyles{
			Label:     "some sub node",
			Shape:     "ellipse",
			FillColor: "#00ffff",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	arrayNode, err := subGraph.AddNode(TemplateClassNode(&symbols.Class{
		Name: "VK\\Utils\\Array",
	}))
	if err != nil {
		log.Fatal(err)
	}

	arrayNode.Scale(1.6)
	ColorizeByScale(arrayNode, 1.6)

	usersNode, err := subGraph.AddNode(TemplateClassNode(&symbols.Class{
		Name:        "VK\\Common\\Users",
		IsInterface: true,
	}))
	if err != nil {
		log.Fatal(err)
	}

	usersNode.Scale(0.6)
	ColorizeByScale(usersNode, 0.6)

	g.AddEdgeByNode(arrayNode, usersNode, graph.EdgeStyles{
		Style: "dashed",
		Color: DefaultEdgeColor,
	})

	g.AddSubGraph(subGraph)

	fmt.Print(g)
}
