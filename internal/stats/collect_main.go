package stats

import (
	"os"

	"github.com/VKCOM/noverify/src/cmd"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
)

func CollectMain() error {
	linter.RegisterBlockChecker(func(ctx *linter.BlockContext) linter.BlockChecker {
		if meta.IsIndexingComplete() {
			return &blockChecker{
				ctx:  ctx,
				root: ctx.RootState()["vklints-root"].(*rootChecker),
			}
		}

		return &blockIndexer{
			ctx:  ctx,
			root: ctx.RootState()["vklints-root"].(*rootIndexer),
		}
	})

	linter.RegisterRootCheckerWithCacher(GlobalCtx, func(ctx *linter.RootContext) linter.RootChecker {
		if meta.IsIndexingComplete() {
			checker := &rootChecker{
				ctx: ctx,
			}
			ctx.State()["vklints-root"] = checker // Save for block checkers
			return checker
		}

		indexer := &rootIndexer{
			ctx:  ctx,
			meta: NewFileMeta(),
		}
		ctx.State()["vklints-root"] = indexer // Save for block checkers
		return indexer
	})

	var arg string
	// var withFile bool
	// var filepath string
	if len(os.Args) > 3 {
		arg = os.Args[1]
		if arg == "-f" {
			// filepath = os.Args[2]
			// withFile = true
			os.Args = append(os.Args[:1], os.Args[3:]...)
		} else {
			os.Args = append(os.Args[:1], os.Args[2:]...)
		}
	}
	//
	// onlyMethods := arg == "methods"
	// onlyFuncs := arg == "funcs"
	// all := arg != "methods" && arg != "funcs"

	// if !withFile {
	_, _ = cmd.Run(&cmd.MainConfig{
		BeforeReport: func(*linter.Report) bool {
			return false
		},
	})
	//
	// 	data, err := GlobalCtx.GobEncode()
	// 	if err != nil {
	// 		fmt.Errorf("error encode global context: %v", err)
	// 	}
	//
	// 	file, err := os.OpenFile("data.pl", os.O_CREATE | os.O_RDWR, os.ModePerm)
	// 	if err != nil {
	// 		log.Fatalf("file not open %v", err)
	// 	}
	//
	// 	fmt.Fprint(file, data)
	// 	file.Close()
	//
	// 	return nil
	// }
	//
	// data, err := ioutil.ReadFile(filepath)
	// if err != nil {
	// 	fmt.Errorf("error read file: %v", err)
	// }
	//
	// err = GlobalCtx.GobDecode(data)
	// if err != nil {
	// 	fmt.Errorf("error decode file: %v", err)
	// }
	// GlobalCtx.Classes.CalculateClassDeps()
	//
	// count := 0
	// for _, class := range GlobalCtx.Classes.Classes {
	// 	efferent := float64(len(class.Deps.Classes))
	// 	afferent := float64(len(class.DepsBy.Classes))
	//
	// 	var stability float64
	// 	if efferent+afferent == 0 {
	// 		stability = 0
	// 	} else {
	// 		stability = efferent / (efferent + afferent)
	// 	}
	//
	// 	count++
	// 	fmt.Printf("%60s: %f\n", class.Name, stability)
	//
	// 	for _, dev := range class.Deps.Classes {
	// 		fmt.Printf("%60s  ← Зависит от %s (абстрактный: %t)\n", "", dev.Name, dev.IsAbstract)
	// 	}
	//
	// 	if len(class.Deps.Classes) == 0 {
	// 		fmt.Printf("%60s  Зависимостей нет\n", "")
	// 	}
	//
	// 	for _, dev := range class.DepsBy.Classes {
	// 		fmt.Printf("%60s  → От него зависит %s\n", "", dev.Name)
	// 	}
	//
	// 	if len(class.DepsBy.Classes) == 0 {
	// 		fmt.Printf("%60s  Никто не зависит от класса\n", "")
	// 	}
	// }
	//
	// fmt.Println(count)

	// funcs, err := GlobalCtx.Funcs.GetFuncsThatCalledFunc(`\VK\API\Library\DeprecatedWrappers::wrapText`)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// funcs := GlobalCtx.Funcs.GetAll(onlyMethods || true, onlyFuncs, all && false, 20, true)
	// for _, fn := range funcs {
	// 	fmt.Fprintf(os.Stdout, "%s: %d (class %s) at %s:%d\n", fn.Name, fn.UsesCount, fn.Class, fn.Pos.Filename, fn.Pos.Line)
	//
	// 	// for _, called := range fn.Called.GetAll(false, true, false, -1, false) {
	// 	// 	fmt.Printf("   called %s\n", called.Name)
	// 	// }
	// }

	// for _, file := range GlobalCtx.Files.Files {
	// 	if len(file.RequiredBy.Files) == 0 {
	// 		continue
	// 	}
	//
	// 	fmt.Println(file.Path, len(file.RequiredBy.Files))
	// }

	// fmt.Println(GlobalCtx.Files.Graphviz())
	//
	// file, err := os.OpenFile("test.gv", os.O_CREATE | os.O_RDWR, os.ModePerm)
	// if err != nil {
	// 	log.Fatalf("file not open %v", err)
	// }
	//
	// fmt.Fprint(file, GlobalCtx.Files.Graphviz())
	// file.Close()

	// for _, class := range GlobalCtx.Classes.Classes {
	// 	for _, method := range class.Methods.Funcs {
	// 		if len(method.Called.Funcs) == 0 {
	// 			continue
	// 		}
	//
	// 		fmt.Println(class.Name)
	//
	// 		fmt.Printf("   %s (%s)\n", method.Name, method.Class.Name)
	// 		for _, called := range method.Called.Funcs {
	// 			if called.Class != nil && method.Class != nil && called.Class.File != method.Class.File {
	// 				fmt.Println("      called class", called.Class)
	// 			}
	// 		}
	// 	}
	// }

	// res := GlobalCtx.Classes.Graphviz()
	// fmt.Println(res)
	//
	// file, err := os.OpenFile("test.gv", os.O_CREATE | os.O_RDWR, os.ModePerm)
	// if err != nil {
	// 	log.Fatalf("file not open %v", err)
	// }
	//
	// fmt.Fprint(file, res)
	// file.Close()

	// classes, ok := GlobalCtx.Classes.GetUsedClassesInClass(`Features`)
	// if !ok {
	// 	log.Fatalf("class not found")
	// }
	//
	// for _, class := range classes.Classes {
	// 	fmt.Println(class.Name)
	// }

	// cmd := exec.Name( "dot","-Tsvg",`test.gv` )
	// cmd.Run()
	//
	// cmd = exec.Name( "open",`test.svg` )
	// cmd.Run()

	// path, err := GlobalCtx.Files.GetFullFileName("member.lib.php")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// files, err := GlobalCtx.Files.GetFilesIncludedFile(path) // `C:\projects\vkcom\www\lib\members.lib.php`
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, file := range files {
	// 	fmt.Println(file.Path)
	// }

	return nil
}
