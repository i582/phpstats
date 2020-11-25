package graph

import (
	"fmt"
)

type Styles struct {
	Label       string
	Padding     float64
	NodeMargin  float64
	GraphMargin float64
	Style       string
	FontColor   string
	BorderColor string
}

func (s Styles) String() string {
	var res string

	if s.Label != "" {
		res += makePropertyWithoutIndent("label", s.Label)
	}
	if s.Padding != 0 {
		res += makePropertyWithoutIndent("pad", fmt.Sprint(s.Padding))
	}
	if s.NodeMargin != 0 {
		res += makePropertyWithoutIndent("ranksep", fmt.Sprint(s.NodeMargin))
	}
	if s.GraphMargin != 0 {
		res += makePropertyWithoutIndent("margin", fmt.Sprint(s.GraphMargin))
	}
	if s.Style != "" {
		res += makePropertyWithoutIndent("style", fmt.Sprint(s.Style))
	}
	if s.FontColor != "" {
		res += makePropertyWithoutIndent("fontcolor", fmt.Sprint(s.FontColor))
	}
	if s.BorderColor != "" {
		res += makePropertyWithoutIndent("color", fmt.Sprint(s.BorderColor))
	}

	return res
}

type NodeStyles struct {
	Label string

	Shape      string
	FillColor  string
	EdgeColor  string
	Style      string
	FontSize   int64
	FontFamily string
	FontColor  string
}

func makeProperty(name, value string) string {
	return fmt.Sprintf("\t\t%s=\"%s\"\n", name, value)
}

func makePropertyWithoutIndent(name, value string) string {
	return fmt.Sprintf("%s=\"%s\"\n", name, value)
}

func (s NodeStyles) String() string {
	var res string
	res += "[\n"

	if s.Label != "" {
		res += makeProperty("label", s.Label)
	}
	if s.Shape != "" {
		res += makeProperty("shape", s.Shape)
	}
	if s.FillColor != "" {
		res += makeProperty("fillcolor", s.FillColor)
	}
	if s.EdgeColor != "" {
		res += makeProperty("color", s.EdgeColor)
	}
	if s.Style != "" {
		res += makeProperty("style", s.Style)
	}
	if s.FontSize != 0 {
		res += makeProperty("fontsize", fmt.Sprint(s.FontSize))
	}
	if s.FontFamily != "" {
		res += makeProperty("fontname", s.FontFamily)
	}

	if res == "[\n" {
		return "[]"
	}

	res += fmt.Sprint("\t]")

	return res
}

type EdgeStyles struct {
	ArrowTail string
	Style     string
	Color     string
	FontColor string
	Width     float64
	Label     string
	ToolTip   string
}

func (s EdgeStyles) String() string {
	var res string
	res += fmt.Sprint("[\n")

	if s.ArrowTail != "" {
		res += makeProperty("arrowtail", s.ArrowTail)
	}
	if s.Style != "" {
		res += makeProperty("style", s.Style)
	}
	if s.Color != "" {
		res += makeProperty("color", s.Color)
	}
	if s.FontColor != "" {
		res += makeProperty("fontcolor", s.FontColor)
	}
	if s.Width != 0 {
		res += makeProperty("penwidth", fmt.Sprint(s.Width))
	}
	if s.Label != "" {
		res += makeProperty("label", s.Label)
	}
	if s.ToolTip != "" {
		res += makeProperty("labeltooltip", s.ToolTip)
	}

	if res == "[\n" {
		return "[]"
	}

	res += fmt.Sprint("\t]")

	return res
}
