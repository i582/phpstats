package uml

import (
	"fmt"
	"strings"

	"github.com/i582/phpstats/internal/stats/symbols"
	"github.com/i582/phpstats/internal/utils"
)

func GetUmlForFunction(f *symbols.Function) string {
	id := utils.NameToIdentifier(f.Name.String())
	idParts := strings.Split(id, `_`)
	shortName := idParts[len(idParts)-1]

	label := fmt.Sprintf("{%s}", shortName)

	return fmt.Sprintf("%s[label = \"%s\"]\n", id, label)
}
