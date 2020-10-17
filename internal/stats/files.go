package stats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"phpstats/internal/utils"
)

type Files struct {
	sync.Mutex

	Files map[string]*File
}

func NewFiles() *Files {
	return &Files{
		Files: map[string]*File{},
	}
}

func (f *Files) Len() int {
	return len(f.Files)
}

func (f *Files) GetFullFileName(name string) ([]string, error) {
	var res []string

	for _, file := range f.Files {
		if strings.Contains(file.Path, name) {
			res = append(res, file.Path)
		}
	}

	if len(res) == 0 {
		return res, fmt.Errorf("class %s not found", name)
	}

	return res, nil
}

func (f *Files) GetAll(count int64, offset int64, sorted bool) []*File {
	var res []*File
	var index int64

	if offset < 0 {
		offset = 0
	}

	for _, fn := range f.Files {
		if !sorted {
			if index+offset > count && count != -1 {
				break
			}
		}

		res = append(res, fn)
		index++
	}

	if sorted {
		sort.Slice(res, func(i, j int) bool {
			if res[i].RequiredBy.Len() == res[j].RequiredBy.Len() {
				return res[i].Name < res[j].Name
			}

			return res[i].RequiredBy.Len() > res[j].RequiredBy.Len()
		})

		if count != -1 {
			if count+offset < int64(len(res)) {
				res = res[:count+offset]
			}

			if offset < int64(len(res)) {
				res = res[offset:]
			}
		}
	}

	return res
}

func (f *Files) Add(file *File) {
	f.Lock()
	defer f.Unlock()
	f.Files[file.Path] = file
}

func (f *Files) Get(path string) (*File, bool) {
	f.Lock()
	defer f.Unlock()
	file, ok := f.Files[path]
	return file, ok
}

func (f *Files) GetOrCreate(path string) *File {
	file, ok := f.Get(path)
	if !ok {
		file := NewFile(path)
		f.Add(file)
		return file
	}

	return file
}

type File struct {
	Name string
	Path string

	RequiredRoot  *Files
	RequiredBlock *Files
	RequiredBy    *Files

	Classes *Classes
	Funcs   *Functions
}

func NewFile(path string) *File {
	return &File{
		Name:          filepath.Base(path),
		Path:          path,
		RequiredRoot:  NewFiles(),
		RequiredBlock: NewFiles(),
		RequiredBy:    NewFiles(),
		Classes:       NewClasses(),
		Funcs:         NewFunctionsInfo(),
	}
}

func GraphvizRecursive(file *File, level int64, visited map[string]struct{}, maxRecursive int64, root, block bool) string {
	if level > maxRecursive {
		return ""
	}

	if _, ok := visited[file.Path]; ok {
		return ""
	}
	visited[file.Path] = struct{}{}

	var res string

	if root {
		for _, f := range file.RequiredRoot.Files {
			str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", file.Path, f.Path)
			if _, ok := visited[str]; !ok {
				res += str
				visited[str] = struct{}{}
			}

			res += GraphvizRecursive(f, level+1, visited, maxRecursive, root, block)
		}
	}

	if block {
		for _, f := range file.RequiredBlock.Files {
			str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", file.Path, f.Path)
			if _, ok := visited[str]; !ok {
				res += str
				visited[str] = struct{}{}
			}

			res += GraphvizRecursive(f, level+1, visited, maxRecursive, root, block)
		}
	}

	return res
}

func (f *File) GraphvizRecursive(maxRecursive int64, root, block bool) string {
	var res string

	res += "digraph test{\n"

	res += GraphvizRecursive(f, 0, make(map[string]struct{}), maxRecursive, root, block)

	res += "}\n"

	return res
}

func RequireRecursive(file *File, level int, visited map[string]struct{}, maxRecursive int, isRootRequire bool) string {
	if level > maxRecursive {
		return ""
	}

	if _, ok := visited[file.Path]; ok {
		return ""
		// return fmt.Sprintf("%s<цикл, файл был подключен выше по иерархии> %s", GenIndent(level), file.ExtraShortString(0))
	}
	visited[file.Path] = struct{}{}

	var res string

	var prefix string

	if isRootRequire {
		prefix = " r↳ "
	} else {
		prefix = " f↳ "
	}

	if len(visited) == 1 {
		res += file.ShortString(level)
	} else {
		res += file.ExtraShortStringWithPrefix(level, prefix)
	}

	var countNonLoopRequiredRoot int
	for _, f := range file.RequiredRoot.Files {
		if _, ok := visited[f.Path]; !ok {
			countNonLoopRequiredRoot++
		}
	}

	if countNonLoopRequiredRoot != 0 && level < maxRecursive {
		// res += fmt.Sprintf("%s   (R) Подключаемые файлы в корне:\n", GenIndent(level))
	}
	for _, f := range file.RequiredRoot.Files {
		res += RequireRecursive(f, level+1, visited, maxRecursive, true)
	}

	var countNonLoopRequiredBlock int
	for _, f := range file.RequiredBlock.Files {
		if _, ok := visited[f.Path]; !ok {
			countNonLoopRequiredBlock++
		}
	}

	if countNonLoopRequiredBlock != 0 && level < maxRecursive {
		// res += fmt.Sprintf("%s   (F) Подключаемые файлы в функциях:\n", GenIndent(level))
	}
	for _, f := range file.RequiredBlock.Files {
		res += RequireRecursive(f, level+1, visited, maxRecursive, false)
	}

	return res
}

func (f *File) FullStringRecursive(maxRecursive int) string {
	var res string

	res += RequireRecursive(f, 0, make(map[string]struct{}), maxRecursive, false)

	return res
}

func (f *File) FullString(level int) string {
	var res string

	res += f.ShortString(level)

	if f.RequiredRoot.Len() != 0 {
		res += fmt.Sprintf("%sПодключаемые файлы в корне:\n", utils.GenIndent(level))
	} else {
		res += fmt.Sprintf("%sПодключаемых файлов в корне нет\n", utils.GenIndent(level))
	}
	for _, f := range f.RequiredRoot.Files {
		res += f.ExtraShortString(level + 1)
	}

	if f.RequiredBlock.Len() != 0 {
		res += fmt.Sprintf("%sПодключаемые файлы в функциях:\n", utils.GenIndent(level))
	} else {
		res += fmt.Sprintf("%sПодключаемых файлов в функциях нет\n", utils.GenIndent(level))
	}
	for _, f := range f.RequiredBlock.Files {
		res += f.ExtraShortString(level + 1)
	}

	res += "\n"

	return res
}

func (f *File) ShortString(level int) string {
	var res string

	res += fmt.Sprintf("%sИмя:  %s\n", utils.GenIndent(level), f.Name)
	res += fmt.Sprintf("%sПуть: %s\n", utils.GenIndent(level), f.Path)

	return res
}

func (f *File) ExtraShortString(level int) string {
	return f.ExtraShortStringWithPrefix(level, "")
}

func (f *File) ExtraShortStringWithPrefix(level int, prefix string) string {
	var res string

	res += fmt.Sprintf("%s%s%-30s (%s)\n", utils.GenIndent(level), prefix, f.Name, f.Path)

	return res
}

func (f *File) AddRequiredFile(file *File) {
	_, ok := f.RequiredBlock.Get(file.Path)
	if !ok {
		f.RequiredBlock.Add(file)
	}
}

func (f *File) AddRequiredRootFile(file *File) {
	_, ok := f.RequiredRoot.Get(file.Path)
	if !ok {
		f.RequiredRoot.Add(file)
	}
}

func (f *File) AddRequiredByFile(file *File) {
	_, ok := f.RequiredBy.Get(file.Path)
	if !ok {
		f.RequiredBy.Add(file)
	}
}

func (f *File) AddClass(class *Class) {
	_, ok := f.Classes.Get(class.Name)
	if !ok {
		f.Classes.Add(class)
	}
}

func (f *File) AddFunc(fun *Function) {
	_, ok := f.Funcs.Get(fun.Name)
	if !ok {
		f.Funcs.Add(fun)
	}
}

// GobEncode is a custom gob marshaller
func (f *File) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(f.Name)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(f.Path)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// GobDecode is custom gob unmarshaller
func (f *File) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&f.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&f.Path)
	if err != nil {
		return err
	}
	return nil
}

// GobDecode is custom gob unmarshaller
func (f *Files) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&f.Files)
	if err != nil {
		return err
	}
	return nil
}

// GobEncode is a custom gob marshaller
func (f *Files) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(f.Files)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
