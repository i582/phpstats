package templates

import (
	"github.com/i582/phpstats/internal/graph"
)

const (
	DefaultSubgraphFillColor = "#F0EEED"
	DefaultFillColor         = "#EDECEA"
	DefaultOutlineColor      = "#B2AC9F"
	DefaultEdgeColor         = "#B2AC9F"

	FillColorLevel1    = "#EDEBE8"
	OutlineColorLevel1 = "#CDC4B6"

	FillColorLevel2    = "#EDE7E2"
	OutlineColorLevel2 = "#B28A60"

	FillColorLevel3    = "#EDE3DC"
	OutlineColorLevel3 = "#B27540"

	FillColorLevel4    = "#EDDBD5"
	OutlineColorLevel4 = "#B22F00"
)

func ColorByScale(scale float64) (string, string) {
	fillColor := DefaultFillColor
	outlineColor := DefaultOutlineColor

	switch {
	case scale >= 1.6:
		fillColor = FillColorLevel4
		outlineColor = OutlineColorLevel4

	case scale >= 1.3:
		fillColor = FillColorLevel3
		outlineColor = OutlineColorLevel3

	case scale >= 1:
		fillColor = FillColorLevel2
		outlineColor = OutlineColorLevel2

	case scale > 0.6:
		fillColor = FillColorLevel1
		outlineColor = OutlineColorLevel1
	}

	return fillColor, outlineColor
}

func ColorizeByScale(node *graph.Node, scale float64) {
	fillColor, outlineColor := ColorByScale(scale)

	node.Styles.FillColor = fillColor
	node.Styles.EdgeColor = outlineColor
}
