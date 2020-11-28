package getter

import (
	"sort"
	"strings"

	"github.com/i582/phpstats/internal/stats/symbols"
)

type FilesGetOptions struct {
	Count       int64
	Offset      int64
	SortColumn  int64
	ReverseSort bool
}

func GetFilesByOptions(f *symbols.Files, opt FilesGetOptions) []*symbols.File {
	files := make([]*symbols.File, 0, f.Len())

	if opt.Offset < 0 {
		opt.Offset = 0
	}

	for _, fn := range f.Files {
		files = append(files, fn)
	}

	sort.Slice(files, func(i, j int) bool {
		var file1 int
		var file2 int
		switch opt.SortColumn {
		case 0, 1: // Name
			fun1 := strings.ToLower(files[i].Name)
			fun2 := strings.ToLower(files[j].Name)
			if opt.ReverseSort {
				fun1, fun2 = fun2, fun1
			}
			return fun1 < fun2

		case 2: // RequiredRoot
			file1 = files[i].RequiredRoot.Len()
			file2 = files[j].RequiredRoot.Len()
		case 3: // RequiredBlock
			file1 = files[i].RequiredBlock.Len()
			file2 = files[j].RequiredBlock.Len()
		case 4: // RequiredBy
			file1 = files[i].RequiredBy.Len()
			file2 = files[j].RequiredBy.Len()
		default:
			return i < j
		}

		if opt.ReverseSort {
			file1, file2 = file2, file1
		}

		return file1 > file2
	})

	if opt.Count+opt.Offset < int64(len(files)) {
		files = files[:opt.Count+opt.Offset]
	}

	if opt.Offset < int64(len(files)) {
		files = files[opt.Offset:]
	}

	return files
}
