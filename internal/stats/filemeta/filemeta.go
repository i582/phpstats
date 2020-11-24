package filemeta

import (
	"github.com/i582/phpstats/internal/stats/symbols"
)

// FileMeta describes the data to be cached.
type FileMeta struct {
	Classes   *symbols.Classes
	Funcs     *symbols.Functions
	Files     *symbols.Files
	Constants *symbols.Constants

	CountCommentLine        int64
	CountAnonymousFunctions int64
}

// NewFileMeta returns a new FileMeta instance with pre-allocated fields.
func NewFileMeta() FileMeta {
	return FileMeta{
		Classes:   symbols.NewClasses(),
		Funcs:     symbols.NewFunctions(),
		Files:     symbols.NewFiles(),
		Constants: symbols.NewConstants(),
	}
}
