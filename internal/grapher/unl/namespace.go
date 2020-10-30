package uml

import (
	"fmt"

	"github.com/i582/phpstats/internal/stats"
	"github.com/i582/phpstats/internal/utils"
)

func GetUmlForNamespace(n *stats.Namespace) string {
	id := utils.ClassNameNormalize(n.FullName)

	label := fmt.Sprintf("{%s}", n.Name)

	return fmt.Sprintf("%s[label = \"%s\"]\n", id, label)
}
