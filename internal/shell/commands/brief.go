package commands

import (
	"fmt"

	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
)

func percent(x, y int64) float64 {
	if y == 0 {
		return 0
	}

	return (float64(x) / float64(y)) * 100
}

func Brief() *shell.Executor {
	briefExecutor := &shell.Executor{
		Name:  "brief",
		Help:  "shows general information",
		Flags: flags.NewFlags(),
		Func: func(c *shell.Context) {
			countLines := walkers.GlobalCtx.Files.CountLines()

			maxMethodCMN, minMethodCMN, avgMethodCMN := walkers.GlobalCtx.Functions.MaxMinAvgMethodCountMagicNumbers()
			maxFunctionCMN, minFunctionCMN, avgFunctionCMN := walkers.GlobalCtx.Functions.MaxMinAvgFunctionsCountMagicNumbers()
			maxClassCMN, minClassCMN, avgClassCMN := walkers.GlobalCtx.Classes.MaxMinAvgCountMagicNumbers()

			maxMethodCC, minMethodCC, avgMethodCC := walkers.GlobalCtx.Functions.MaxMinAvgMethodCyclomaticComplexity()
			maxFunctionCC, minFunctionCC, avgFunctionCC := walkers.GlobalCtx.Functions.MaxMinAvgFunctionsCyclomaticComplexity()
			maxClassCC, minClassCC, avgClassCC := walkers.GlobalCtx.Classes.MaxMinAvgCyclomaticComplexity()

			fmt.Printf("General project statistics\n\n")

			fmt.Println("Size")

			fmt.Printf("    Lines of Code (LOC):                           %8d\n", countLines)
			fmt.Printf("    Comment Lines of Code (CLOC):                  %8d (%.2f%%)\n", walkers.GlobalCtx.CountCommentLine, percent(walkers.GlobalCtx.CountCommentLine, countLines))
			fmt.Printf("    Non-Comment Lines of Code (NCLOC):             %8d (%.2f%%)\n", countLines-walkers.GlobalCtx.CountCommentLine, 100-percent(walkers.GlobalCtx.CountCommentLine, countLines))
			fmt.Println()

			fmt.Println("Metrics")

			fmt.Printf("    Cyclomatic Complexity\n")

			fmt.Printf("        Average Complexity per Class:              %8.2f\n", avgClassCC)
			fmt.Printf("            Maximum Class Complexity:              %8.2f\n", maxClassCC)
			fmt.Printf("            Minimum Class Complexity:              %8.2f\n", minClassCC)

			fmt.Printf("        Average Complexity per Method:             %8.2f\n", avgMethodCC)
			fmt.Printf("            Maximum Method Complexity:             %8.2f\n", maxMethodCC)
			fmt.Printf("            Minimum Method Complexity:             %8.2f\n", minMethodCC)

			fmt.Printf("        Average Complexity per Functions:          %8.2f\n", avgFunctionCC)
			fmt.Printf("            Maximum Functions Complexity:          %8.2f\n", maxFunctionCC)
			fmt.Printf("            Minimum Functions Complexity:          %8.2f\n", minFunctionCC)
			fmt.Println()

			fmt.Printf("    Count of Magic Numbers\n")

			fmt.Printf("        Average Class Count:                       %8d\n", avgClassCMN)
			fmt.Printf("            Maximum Class Count:                   %8d\n", maxClassCMN)
			fmt.Printf("            Minimum Class Count:                   %8d\n", minClassCMN)

			fmt.Printf("        Average Method Count:                      %8d\n", avgMethodCMN)
			fmt.Printf("            Maximum Method Count:                  %8d\n", maxMethodCMN)
			fmt.Printf("            Minimum Method Count:                  %8d\n", minMethodCMN)

			fmt.Printf("        Average Functions Count:                   %8d\n", avgFunctionCMN)
			fmt.Printf("            Maximum Method Count:                  %8d\n", maxFunctionCMN)
			fmt.Printf("            Minimum Method Count:                  %8d\n", minFunctionCMN)
			fmt.Println()

			fmt.Println("Structure")

			fmt.Printf("    Files:                                         %8d\n", walkers.GlobalCtx.Files.Len())
			fmt.Printf("    Namespaces:                                    %8d\n", walkers.GlobalCtx.Namespaces.Count())
			fmt.Printf("    Interfaces:                                    %8d\n", walkers.GlobalCtx.Classes.CountIfaces())
			fmt.Printf("    Classes                                        %8d\n", int64(walkers.GlobalCtx.Classes.Len())-walkers.GlobalCtx.Classes.CountIfaces())
			fmt.Printf("        Abstract Classes:                          %8d (%.2f%%)\n", walkers.GlobalCtx.Classes.CountAbstractClasses(), percent(walkers.GlobalCtx.Classes.CountAbstractClasses(), int64(walkers.GlobalCtx.Classes.Len())))
			fmt.Printf("        Concrete Classes:                          %8d (%.2f%%)\n", walkers.GlobalCtx.Classes.CountClasses(), 100-percent(walkers.GlobalCtx.Classes.CountAbstractClasses(), int64(walkers.GlobalCtx.Classes.Len())))
			fmt.Printf("    Methods:                                       %8d\n", walkers.GlobalCtx.Functions.CountMethods())
			fmt.Printf("    Constants:                                     %8d\n", walkers.GlobalCtx.Constants.Len())
			fmt.Printf("    Functions:\n")
			fmt.Printf("        Named Functions:                           %8d (%.2f%%)\n", walkers.GlobalCtx.Functions.CountFunctions(), percent(walkers.GlobalCtx.Functions.CountFunctions(), walkers.GlobalCtx.Functions.CountFunctions()+walkers.GlobalCtx.CountAnonymousFunctions))
			fmt.Printf("        Anonymous Functions:                       %8d (%.2f%%)\n", walkers.GlobalCtx.CountAnonymousFunctions, percent(walkers.GlobalCtx.CountAnonymousFunctions, walkers.GlobalCtx.Functions.CountFunctions()+walkers.GlobalCtx.CountAnonymousFunctions))

			fmt.Println()
		},
	}

	return briefExecutor
}
