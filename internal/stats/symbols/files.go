package symbols

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Files struct {
	m sync.Mutex

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

func (f *Files) CountLines() int64 {
	var count int64
	for _, file := range f.Files {
		count += file.CountLines
	}
	return count
}

func (f *Files) GetFileByPartOfName(name string) (*File, error) {
	classes, err := f.GetFullFileName(name)
	if err != nil {
		return nil, err
	}

	class, found := f.Get(classes[0])
	if !found {
		return nil, fmt.Errorf("file %s not found", name)
	}
	return class, nil
}

func (f *Files) GetFullFileName(name string) ([]string, error) {
	var res []string

	for _, file := range f.Files {
		if file.Path == name {
			return []string{file.Path}, nil
		}

		if strings.Contains(file.Path, name) {
			res = append(res, file.Path)
		}
	}

	if len(res) == 0 {
		return res, fmt.Errorf("file %s not found", name)
	}

	return res, nil
}

func (f *Files) Add(file *File) {
	f.m.Lock()
	f.Files[file.Path] = file
	f.m.Unlock()
}

func (f *Files) Get(path string) (*File, bool) {
	f.m.Lock()
	file, ok := f.Files[path]
	f.m.Unlock()
	return file, ok
}

type File struct {
	Name string
	Path string

	RequiredRoot  *Files
	RequiredBlock *Files
	RequiredBy    *Files

	Classes *Classes
	Funcs   *Functions

	CountLines int64
}

func NewFile(path string) *File {
	return &File{
		Name:          filepath.Base(path),
		Path:          path,
		RequiredRoot:  NewFiles(),
		RequiredBlock: NewFiles(),
		RequiredBy:    NewFiles(),
		Classes:       NewClasses(),
		Funcs:         NewFunctions(),
	}
}

var filenameNormalizer = regexp.MustCompile("[^0-9a-zA-Z]")

func (f *File) UniqueId() string {
	dir, _ := filepath.Split(f.Path)
	return filepath.Base(dir) + "__" + filenameNormalizer.ReplaceAllString(f.Name, "_")
}

func (f *File) AddRequiredFile(file *File) {
	f.RequiredBlock.Add(file)
}

func (f *File) AddRequiredRootFile(file *File) {
	f.RequiredRoot.Add(file)
}

func (f *File) AddRequiredByFile(file *File) {
	f.RequiredBy.Add(file)
}

func (f *File) AddClass(class *Class) {
	f.Classes.Add(class)
}

func (f *File) AddFunc(fun *Function) {
	f.Funcs.Add(fun)
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
