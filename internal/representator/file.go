package representator

import (
	"encoding/json"
	"fmt"

	"github.com/i582/phpstats/internal/stats"
)

type FileData struct {
	Name string `json:"name"`
	Path string `json:"path"`

	CountRequiredRoot  int64
	CountRequiredBlock int64
	CountRequiredBy    int64

	requiredRoot  *stats.Files
	requiredBlock *stats.Files
	requiredBy    *stats.Files
}

func fileToData(f *stats.File) *FileData {
	return &FileData{
		Name: f.Name,
		Path: f.Path,

		CountRequiredBlock: int64(f.RequiredBlock.Len()),
		CountRequiredRoot:  int64(f.RequiredRoot.Len()),
		CountRequiredBy:    int64(f.RequiredBy.Len()),

		requiredBlock: f.RequiredBlock,
		requiredRoot:  f.RequiredRoot,
		requiredBy:    f.RequiredBy,
	}
}

func GetShortFileRepr(f *stats.File) string {
	data := fileToData(f)

	return fmt.Sprintf("%-40s (%s)", data.Name, data.Path)
}

func GetFileRepr(f *stats.File) string {
	data := fileToData(f)

	var res string

	res += fmt.Sprintf("File %s\n", data.Name)
	res += fmt.Sprintf("  Path %s\n", data.Path)

	res += fmt.Sprintf("  Include files at the root:      %d\n", data.CountRequiredRoot)
	for _, f := range data.requiredRoot.Files {
		res += fmt.Sprintf("\t%-40s (%s)\n", f.Name, f.Path)
	}

	res += fmt.Sprintf("  Include files in the functions: %d\n", data.CountRequiredBlock)
	for _, f := range data.requiredBlock.Files {
		res += fmt.Sprintf("\t%-40s (%s)\n", f.Name, f.Path)
	}

	res += fmt.Sprintf("  Count of required:              %d\n", data.CountRequiredBy)

	return res
}

func GetJsonFileRepr(f *stats.File) ([]byte, error) {
	data := fileToData(f)

	res, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return res, nil
}
