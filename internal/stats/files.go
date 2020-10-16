package stats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Files struct {
	sync.Mutex

	Files map[string]*File
}

func (f *Files) String() string {
	var res string

	for _, file := range f.Files {
		res += fmt.Sprint(file)
	}

	return res
}

func NewFiles() *Files {
	return &Files{
		Files: map[string]*File{},
	}
}

func (f *Files) Len() int {
	return len(f.Files)
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

func (f *Files) GetFullFileName(name string) (string, error) {
	for _, file := range f.Files {
		if strings.Contains(file.Path, name) {
			return file.Path, nil
		}
	}

	return "", fmt.Errorf("file %s not found", name)
}

func (f *Files) Graphviz() string {
	var res string
	res += "digraph test{\n"

	var count int

	// Outer:
	for _, file := range f.Files {
		for _, requiredFile := range file.RequiredRoot.Files {
			// if !strings.Contains(file.Path, `www\VK\API`) {
			// 	continue
			// }
			res += fmt.Sprintf("   \"%s\" -> \"%s\"\n", file.Path, requiredFile.Path)
			count++
		}

		for _, requiredFile := range file.RequiredBlock.Files {
			// if !strings.Contains(file.Path, `www\VK\API`) {
			// 	continue
			// }
			res += fmt.Sprintf("   \"%s\" -> \"%s\" [style=dotted]\n", file.Path, requiredFile.Path)
			count++
		}
	}
	res += "}"
	return res
}

func (f *Files) GetFilesIncludedFile(name string) ([]*File, error) {
	file, ok := f.Get(name)

	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}

	return file.RequiredBy.GetAll(-1, false), nil
}

func (f *Files) GetAll(count int, sorted bool) []*File {
	var res []*File
	var index int

	for _, fn := range f.Files {
		if !sorted {
			if index > count && count != -1 {
				break
			}
		}

		res = append(res, fn)
		index++
	}

	if sorted {
		sort.Slice(res, func(i, j int) bool {
			return len(res[i].RequiredBy.Files) > len(res[j].RequiredBy.Files)
		})
		if count < len(res) {
			res = res[:count]
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
	if file, ok := f.Get(path); !ok {
		file := NewFile(path)
		f.Add(file)
		return file
	} else {
		return file
	}
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

func genIndent(level int) string {
	var res string
	for i := 0; i < level; i++ {
		res += "   "
	}
	return res
}

func GraphvizRecursive(file *File, level int, visited map[string]struct{}, maxRecursive int, isRootRequire bool) string {
	if level > maxRecursive {
		return ""
	}

	if _, ok := visited[file.Path]; ok {
		return ""
		// return fmt.Sprintf("%s<цикл, файл был подключен выше по иерархии> %s", genIndent(level), file.ExtraShortString(0))
	}
	visited[file.Path] = struct{}{}

	var res string

	for _, f := range file.RequiredRoot.Files {
		str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", file.Path, f.Path)
		if _, ok := visited[str]; !ok {
			res += str
			visited[str] = struct{}{}
		}

		res += GraphvizRecursive(f, level+1, visited, maxRecursive, true)
	}

	for _, f := range file.RequiredBlock.Files {
		str := fmt.Sprintf("   \"%s\" -> \"%s\"\n", file.Path, f.Path)
		if _, ok := visited[str]; !ok {
			res += str
			visited[str] = struct{}{}
		}

		res += GraphvizRecursive(f, level+1, visited, maxRecursive, false)
	}

	return res
}

func (f *File) GraphvizRecursive(maxRecursive int) string {
	var res string

	res += "digraph test{\n"

	res += GraphvizRecursive(f, 0, make(map[string]struct{}), maxRecursive, false)

	res += "}\n"

	return res
}

func RequireRecursive(file *File, level int, visited map[string]struct{}, maxRecursive int, isRootRequire bool) string {
	if level > maxRecursive {
		return ""
	}

	if _, ok := visited[file.Path]; ok {
		return ""
		// return fmt.Sprintf("%s<цикл, файл был подключен выше по иерархии> %s", genIndent(level), file.ExtraShortString(0))
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

	// used-in file TicketAuthorRights.php
	if countNonLoopRequiredRoot != 0 && level < maxRecursive {
		// res += fmt.Sprintf("%s   (R) Подключаемые файлы в корне:\n", genIndent(level))
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
		// res += fmt.Sprintf("%s   (F) Подключаемые файлы в функциях:\n", genIndent(level))
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
		res += fmt.Sprintf("%sПодключаемые файлы в корне:\n", genIndent(level))
	} else {
		res += fmt.Sprintf("%sПодключаемых файлов в корне нет\n", genIndent(level))
	}
	for _, f := range f.RequiredRoot.Files {
		res += f.ExtraShortString(level + 1)
	}

	if f.RequiredBlock.Len() != 0 {
		res += fmt.Sprintf("%sПодключаемые файлы в функциях:\n", genIndent(level))
	} else {
		res += fmt.Sprintf("%sПодключаемых файлов в функциях нет\n", genIndent(level))
	}
	for _, f := range f.RequiredBlock.Files {
		res += f.ExtraShortString(level + 1)
	}

	res += "\n"

	return res
}

func (f *File) ShortString(level int) string {
	var res string

	res += fmt.Sprintf("%sИмя:  %s\n", genIndent(level), f.Name)
	res += fmt.Sprintf("%sПуть: %s\n", genIndent(level), f.Path)

	return res
}

func (f *File) ExtraShortString(level int) string {
	return f.ExtraShortStringWithPrefix(level, "")
}

func (f *File) ExtraShortStringWithPrefix(level int, prefix string) string {
	var res string

	res += fmt.Sprintf("%s%s%-30s (%s)\n", genIndent(level), prefix, f.Name, f.Path)

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
