package stats

import (
	"encoding/gob"
	"io"
	"log"

	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
)

type GlobalContext struct {
	Funcs     *Functions
	Classes   *Classes
	Files     *Files
	Constants *Constants

	Decoded bool
	Encoded bool
}

func (ctx *GlobalContext) Version() string {
	return "1.0.0"
}

func (ctx *GlobalContext) Encode(writer io.Writer, checker linter.RootChecker) error {
	if meta.IsLoadingStubs() {
		return nil
	}

	ind := checker.(*rootIndexer)

	enc := gob.NewEncoder(writer)
	if err := enc.Encode(&ind.meta); err != nil {
		log.Printf("cache error: encode %s: %v", ind.ctx.Filename(), err)
		return err
	}

	return nil
}

func (ctx *GlobalContext) Decode(r io.Reader, filename string) error {
	if meta.IsLoadingStubs() {
		return nil
	}

	var m FileMeta

	dec := gob.NewDecoder(r)
	if err := dec.Decode(&m); err != nil {
		log.Printf("cache error: decode %s: %v", filename, err)
		return err
	}

	ctx.updateMeta(&m)

	return nil
}

func (ctx *GlobalContext) updateMeta(f *FileMeta) {
	for _, file := range f.Files.Files {
		f := NewFile(file.Path)

		ctx.Files.Add(f)
	}

	for _, class := range f.Classes.Classes {
		var cl *Class

		file, ok := ctx.Files.Get(class.File.Path)
		if !ok {
			log.Fatal("file not found")
		}

		if class.IsInterface {
			cl = NewInterface(class.Name, file)
		} else if class.IsAbstract {
			cl = NewAbstractClass(class.Name, file)
		} else {
			cl = NewClass(class.Name, file)
		}

		ctx.Classes.Add(cl)
	}

	for _, fn := range f.Funcs.Funcs {
		fun := NewFunctionInfo(fn.Name, fn.Pos)

		if fun.IsMethod() {
			class, ok := ctx.Classes.Get(fun.Name.ClassName)
			if !ok {
				return
			}
			class.AddMethod(fun)
			fun.Class = class
		}

		ctx.Funcs.Add(fun)
	}

	for _, constant := range f.Constants.Constants {
		ctx.Constants.Add(constant)
	}
}

func NewGlobalContext() *GlobalContext {
	return &GlobalContext{
		Funcs:     NewFunctionsInfo(),
		Classes:   NewClasses(),
		Files:     NewFiles(),
		Constants: NewConstants(),
	}
}

var GlobalCtx = NewGlobalContext()

type FileMeta struct {
	Classes   *Classes
	Funcs     *Functions
	Files     *Files
	Constants *Constants
}

func NewFileMeta() FileMeta {
	return FileMeta{
		Classes:   NewClasses(),
		Funcs:     NewFunctionsInfo(),
		Files:     NewFiles(),
		Constants: NewConstants(),
	}
}
