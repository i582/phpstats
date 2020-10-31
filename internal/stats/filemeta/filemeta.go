package filemeta

import (
	"github.com/i582/phpstats/internal/stats/symbols"
)

type FileMeta struct {
	Classes   *symbols.Classes
	Funcs     *symbols.Functions
	Files     *symbols.Files
	Constants *symbols.Constants
}

func NewFileMeta() FileMeta {
	return FileMeta{
		Classes:   symbols.NewClasses(),
		Funcs:     symbols.NewFunctions(),
		Files:     symbols.NewFiles(),
		Constants: symbols.NewConstants(),
	}
}
