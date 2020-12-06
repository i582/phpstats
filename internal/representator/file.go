package representator

import (
	"encoding/json"
	"fmt"

	"github.com/gookit/color"
	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/stats/symbols"
)

type FileData struct {
	Name string `json:"name"`
	Path string `json:"path"`

	CountRequiredRoot  int64 `json:"countRequiredInRoot"`
	CountRequiredBlock int64 `json:"countRequiredInBlock"`
	CountRequiredBy    int64 `json:"countRequiredBy"`

	requiredRoot  *symbols.Files
	requiredBlock *symbols.Files
	requiredBy    *symbols.Files
}

func fileToData(f *symbols.File) *FileData {
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

func GetShortStringFileRepr(f *symbols.File) string {
	if f == nil {
		return ""
	}

	data := fileToData(f)

	return fmt.Sprintf("%-40s (%s)", data.Name, data.Path)
}

func GetStringFileRepr(f *symbols.File) string {
	if f == nil {
		return ""
	}

	data := fileToData(f)

	var res string

	res += fmt.Sprintf("File %s\n", data.Name)
	res += color.Gray.Sprintf("  Path %s\n", data.Path)

	res += cfmt.Sprintf("  {{Include files at the root}}::green:      %s\n", ColorOutputIntZeroableValue(data.CountRequiredRoot))
	res += cfmt.Sprintf("  {{Include files in the functions}}::green: %s\n", ColorOutputIntZeroableValue(data.CountRequiredBlock))
	res += cfmt.Sprintf("  {{Count of required}}::green:              %s\n", ColorOutputIntZeroableValue(data.CountRequiredBy))

	return res
}

func GetJsonFileRepr(f *symbols.File) (string, error) {
	data := fileToData(f)

	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetPrettifyJsonFilesRepr(f []*symbols.File) (string, error) {
	data := make([]*FileData, 0, len(f))

	for _, file := range f {
		data = append(data, fileToData(file))
	}

	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetJsonFileReprWithFlag(f *symbols.File) (string, error) {
	type Response struct {
		Data  *FileData `json:"data"`
		Found bool      `json:"found"`
	}
	var resp Response

	resp.Data = fileToData(f)
	resp.Found = f != nil

	res, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
