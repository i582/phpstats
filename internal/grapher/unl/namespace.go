package uml

import (
	"fmt"

	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func GetUmlForNamespace(n *symbols.Namespace) string {
	id := utils.NameToIdentifier(n.FullName)

	label := fmt.Sprintf("{%s}", n.Name)

	return fmt.Sprintf("%s[label = \"%s\"]\n", id, label)
}
