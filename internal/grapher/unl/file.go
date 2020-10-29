package uml

import (
	"fmt"

	"github.com/i582/phpstats/internal/stats"
)

func GetUmlForFile(f *stats.File) string {
	id := f.UniqueId()
	shortName := f.Name

	label := fmt.Sprintf("{%s}", shortName)

	return fmt.Sprintf("%s[label = \"%s\"]\n", id, label)
}
