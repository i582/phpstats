package graph

import (
	"strings"
)

func genIndent(level int64) string {
	var res string
	for i := int64(0); i < level; i++ {
		res += "\t"
	}
	return res
}

func indentText(text string, level int64) string {
	lines := strings.Split(text, "\n")

	for index := range lines {
		lines[index] = genIndent(level) + lines[index]
	}

	return strings.Join(lines, "\n")
}
