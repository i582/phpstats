package walkers

import (
	"encoding/gob"
	"io"
	"log"

	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/cheggaaa/pb/v3"

	"github.com/i582/phpstats/internal/stats/filemeta"
	"github.com/i582/phpstats/internal/stats/symbols"
)

var GlobalCtx = NewGlobalContext()

type globalContext struct {
	Functions  *symbols.Functions
	Classes    *symbols.Classes
	Files      *symbols.Files
	Constants  *symbols.Constants
	Namespaces *symbols.Namespaces

	ProjectRoot string
	BarLinting  *pb.ProgressBar
}

func NewGlobalContext() *globalContext {
	return &globalContext{
		Functions:  symbols.NewFunctions(),
		Classes:    symbols.NewClasses(),
		Files:      symbols.NewFiles(),
		Constants:  symbols.NewConstants(),
		Namespaces: symbols.NewNamespaces(),
	}
}

func (ctx *globalContext) Version() string {
	return "1.0.1"
}

func (ctx *globalContext) Encode(writer io.Writer, checker linter.RootChecker) error {
	if meta.IsLoadingStubs() {
		return nil
	}

	ind := checker.(*rootIndexer)

	enc := gob.NewEncoder(writer)
	if err := enc.Encode(&ind.Meta); err != nil {
		log.Printf("cache error: encode %s: %v", ind.Ctx.Filename(), err)
		return err
	}

	return nil
}

func (ctx *globalContext) Decode(r io.Reader, filename string) error {
	if meta.IsLoadingStubs() {
		return nil
	}

	var m filemeta.FileMeta

	dec := gob.NewDecoder(r)
	if err := dec.Decode(&m); err != nil {
		log.Printf("cache error: decode %s: %v", filename, err)
		return err
	}

	ctx.UpdateMeta(&m)

	return nil
}

func (ctx *globalContext) UpdateMeta(f *filemeta.FileMeta) {
	for _, file := range f.Files.Files {
		f := symbols.NewFile(file.Path)

		ctx.Files.Add(f)
	}

	for _, class := range f.Classes.Classes {
		var cl *symbols.Class

		file, ok := ctx.Files.Get(class.File.Path)
		if !ok {
			log.Fatal("file not found")
		}

		if class.IsInterface {
			cl = symbols.NewInterface(class.Name, file)
		} else if class.IsAbstract {
			cl = symbols.NewAbstractClass(class.Name, file)
		} else {
			cl = symbols.NewClass(class.Name, file)
		}

		cl.Vendor = class.Vendor

		ctx.Classes.Add(cl)
	}

	for _, fn := range f.Funcs.Funcs {
		fun := symbols.NewFunction(fn.Name, fn.Pos)

		if fun.IsMethod() {
			class, ok := ctx.Classes.Get(fun.Name.ClassName)
			if !ok {
				return
			}
			class.AddMethod(fun)
			fun.Class = class
		}

		fun.CyclomaticComplexity = fn.CyclomaticComplexity
		fun.CountMagicNumbers = fn.CountMagicNumbers

		ctx.Functions.Add(fun)
	}

	if f.Constants != nil {
		for _, constant := range f.Constants.Constants {
			ctx.Constants.Add(constant)
		}
	}
}
