package uml

import (
	"fmt"
	"strings"

	"github.com/i582/phpstats/internal/stats"
	"github.com/i582/phpstats/internal/utils"
)

func GetUmlForFunction(f *stats.Function) string {
	id := utils.ClassNameNormalize(f.Name.String())
	idParts := strings.Split(id, `_`)
	shortName := idParts[len(idParts)-1]

	var typ string

	label := fmt.Sprintf("{%s%s}", typ, shortName)

	return fmt.Sprintf("%s[label = \"%s\"]\n", id, label)
}
