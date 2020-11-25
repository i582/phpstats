package templates

import (
	"github.com/i582/phpstats/internal/graph"
)

func TemplateImplementEdgeStyle() graph.EdgeStyles {
	return graph.EdgeStyles{
		ArrowTail: "empty",
		Style:     "dashed",
		Width:     2,
		Color:     OutlineColorLevel4,
		FontColor: OutlineColorLevel4,
		Label:     "impl",
		ToolTip:   "Implement",
	}
}

func TemplateExtendEdgeStyle() graph.EdgeStyles {
	return graph.EdgeStyles{
		ArrowTail: "empty",
		Style:     "dotted",
		Width:     2,
		Color:     OutlineColorLevel3,
		FontColor: OutlineColorLevel3,
		Label:     "ext",
		ToolTip:   "Extends",
	}
}

func TemplateClassConnectionEdgeStyle() graph.EdgeStyles {
	return graph.EdgeStyles{
		ArrowTail: "empty",
		Color:     DefaultEdgeColor,
	}
}

func TemplateFunctionConnectionEdgeStyle() graph.EdgeStyles {
	return graph.EdgeStyles{
		ArrowTail: "empty",
		Color:     DefaultEdgeColor,
	}
}

func TemplateRootFileEdgeStyle() graph.EdgeStyles {
	return graph.EdgeStyles{
		ArrowTail: "empty",
		Style:     "dashed",
		Width:     2,
		Color:     OutlineColorLevel4,
		FontColor: OutlineColorLevel4,
		ToolTip:   "Included in root",
	}
}

func TemplateBlockFileEdgeStyle() graph.EdgeStyles {
	return graph.EdgeStyles{
		ArrowTail: "empty",
		Style:     "dotted",
		Width:     2,
		Color:     OutlineColorLevel3,
		FontColor: OutlineColorLevel3,
		ToolTip:   "Included in block",
	}
}
