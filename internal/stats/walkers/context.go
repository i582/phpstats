package walkers

import (
	"encoding/gob"
	"io"
	"log"
	"path/filepath"
	"regexp"

	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/cheggaaa/pb/v3"

	"github.com/i582/phpstats/internal/config"
	"github.com/i582/phpstats/internal/stats/filemeta"
	"github.com/i582/phpstats/internal/stats/symbols"
)

// GlobalCtx stores all functions, classes, files, constants and namespaces.
var GlobalCtx = newGlobalContext()

type globalContext struct {
	Functions  *symbols.Functions
	Classes    *symbols.Classes
	Files      *symbols.Files
	Constants  *symbols.Constants
	Namespaces *symbols.Namespaces

	Packages *config.Packages

	ProjectRoot   string
	ExcludeRegexp *regexp.Regexp
	BarLinting    *pb.ProgressBar

	CountFiles int64

	CountCommentLine        int64
	CountAnonymousFunctions int64
}

func newGlobalContext() *globalContext {
	return &globalContext{
		Functions:  symbols.NewFunctions(),
		Classes:    symbols.NewClasses(),
		Files:      symbols.NewFiles(),
		Constants:  symbols.NewConstants(),
		Namespaces: symbols.NewNamespaces(),
		Packages:   &config.Packages{},
	}
}

// Version returns the current version of the cache.
func (ctx *globalContext) Version() string {
	return "1.0.1"
}

// Encode caches the data of one rootWalker of one file.
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

// Decode recovers data from cache.
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

	ctx.UpdateMeta(&m, filename)

	return nil
}

// UpdateMeta recovers data by collecting it from each file.
func (ctx *globalContext) UpdateMeta(f *filemeta.FileMeta, filename string) {
	for range f.Files.Files {
		ctx.CountFiles++
	}

	if ctx.ExcludeRegexp != nil && ctx.ExcludeRegexp.MatchString(filepath.ToSlash(filename)) {
		return
	}

	ctx.CountCommentLine += f.CountCommentLine
	ctx.CountAnonymousFunctions += f.CountAnonymousFunctions

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
		} else if class.IsTrait {
			cl = symbols.NewTrait(class.Name, file)
		} else {
			cl = symbols.NewClass(class.Name, file)
		}

		cl.IsVendor = class.IsVendor

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

			var cl *symbols.Class

			file, ok := ctx.Files.Get(constant.Class.File.Path)
			if !ok {
				log.Print("file not found")
			}

			if constant.Class.IsInterface {
				cl = symbols.NewInterface(constant.Class.Name, file)
			} else if constant.Class.IsAbstract {
				cl = symbols.NewAbstractClass(constant.Class.Name, file)
			} else if constant.Class.IsTrait {
				cl = symbols.NewTrait(constant.Class.Name, file)
			} else {
				cl = symbols.NewClass(constant.Class.Name, file)
			}

			cl.IsVendor = constant.Class.IsVendor

			constant.Class = cl
		}
	}
}
