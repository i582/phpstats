package commands

import (
	"fmt"

	"github.com/i582/cfmt"

	"github.com/i582/phpstats/internal/representator"
	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
	"github.com/i582/phpstats/internal/stats/walkers"
	"github.com/i582/phpstats/internal/utils"
)

func colorInt(value int64) string {
	return representator.ColorWidthOutputIntZeroableValue(value, 8)
}

func colorFloat(value float64) string {
	return representator.ColorWidthOutputFloatZeroableValue(value, 8)
}

func colorPercent(value float64) string {
	return representator.ColorOutputFloatZeroablePercentValue(value)
}

func Brief() *shell.Executor {
	briefExecutor := &shell.Executor{
		Name:  "brief",
		Help:  "shows brief information about the project",
		Flags: flags.NewFlags(),
		Func: func(c *shell.Context) {
			countLines := walkers.GlobalCtx.Files.CountLines()

			maxMethodCMN, minMethodCMN, avgMethodCMN := walkers.GlobalCtx.Functions.MaxMinAvgMethodCountMagicNumbers()
			maxFunctionCMN, minFunctionCMN, avgFunctionCMN := walkers.GlobalCtx.Functions.MaxMinAvgFunctionsCountMagicNumbers()
			maxClassCMN, minClassCMN, avgClassCMN := walkers.GlobalCtx.Classes.MaxMinAvgCountMagicNumbers()

			maxMethodCC, minMethodCC, avgMethodCC := walkers.GlobalCtx.Functions.MaxMinAvgMethodCyclomaticComplexity()
			maxFunctionCC, minFunctionCC, avgFunctionCC := walkers.GlobalCtx.Functions.MaxMinAvgFunctionsCyclomaticComplexity()
			maxClassCC, minClassCC, avgClassCC := walkers.GlobalCtx.Classes.MaxMinAvgCyclomaticComplexity()

			cfmt.Printf("General '%s' project statistics\n\n", walkers.GlobalCtx.ProjectName)

			cfmt.Println("Size")

			cfmt.Printf("    {{Lines of Code (LOC)}}::green:                           %s\n", colorInt(countLines))
			cfmt.Printf("    {{Comment Lines of Code (CLOC)}}::green:                  %s %s\n", colorInt(walkers.GlobalCtx.CountCommentLine), colorPercent(utils.Percent(walkers.GlobalCtx.CountCommentLine, countLines)))
			cfmt.Printf("    {{Non-Comment Lines of Code (NCLOC)}}::green:             %s %s\n", colorInt(countLines-walkers.GlobalCtx.CountCommentLine), colorPercent(100-utils.Percent(walkers.GlobalCtx.CountCommentLine, countLines)))
			cfmt.Println()

			cfmt.Println("Metrics")

			cfmt.Printf("    {{Cyclomatic Complexity}}::green\n")

			cfmt.Printf("        {{Average Complexity per Class}}::green:              %s\n", colorFloat(avgClassCC))
			cfmt.Printf("            {{Maximum Class Complexity}}::green:              %s\n", colorFloat(maxClassCC))
			cfmt.Printf("            {{Minimum Class Complexity}}::green:              %s\n", colorFloat(minClassCC))

			cfmt.Printf("        {{Average Complexity per Method}}::green:             %s\n", colorFloat(avgMethodCC))
			cfmt.Printf("            {{Maximum Method Complexity}}::green:             %s\n", colorFloat(maxMethodCC))
			cfmt.Printf("            {{Minimum Method Complexity}}::green:             %s\n", colorFloat(minMethodCC))

			cfmt.Printf("        {{Average Complexity per Functions}}::green:          %s\n", colorFloat(avgFunctionCC))
			cfmt.Printf("            {{Maximum Functions Complexity}}::green:          %s\n", colorFloat(maxFunctionCC))
			cfmt.Printf("            {{Minimum Functions Complexity}}::green:          %s\n", colorFloat(minFunctionCC))
			cfmt.Println()

			cfmt.Printf("    {{Count of Magic Numbers}}::green\n")

			cfmt.Printf("        {{Average Class Count}}::green:                       %s\n", colorInt(avgClassCMN))
			cfmt.Printf("            {{Maximum Class Count}}::green:                   %s\n", colorInt(maxClassCMN))
			cfmt.Printf("            {{Minimum Class Count}}::green:                   %s\n", colorInt(minClassCMN))

			cfmt.Printf("        {{Average Method Count}}::green:                      %s\n", colorInt(avgMethodCMN))
			cfmt.Printf("            {{Maximum Method Count}}::green:                  %s\n", colorInt(maxMethodCMN))
			cfmt.Printf("            {{Minimum Method Count}}::green:                  %s\n", colorInt(minMethodCMN))

			cfmt.Printf("        {{Average Functions Count}}::green:                   %s\n", colorInt(avgFunctionCMN))
			cfmt.Printf("            {{Maximum Method Count}}::green:                  %s\n", colorInt(maxFunctionCMN))
			cfmt.Printf("            {{Minimum Method Count}}::green:                  %s\n", colorInt(minFunctionCMN))
			cfmt.Println()

			cfmt.Println("Structure")

			cfmt.Printf("    {{Files}}::green:                                         %s\n", colorInt(int64(walkers.GlobalCtx.Files.Len())))
			cfmt.Printf("    {{Namespaces}}::green:                                    %s\n", colorInt(walkers.GlobalCtx.Namespaces.Count()))
			cfmt.Printf("    {{Interfaces}}::green:                                    %s\n", colorInt(walkers.GlobalCtx.Classes.CountIfaces()))
			cfmt.Printf("    {{Traits}}::green:                                        %s\n", colorInt(walkers.GlobalCtx.Classes.CountTraits()))
			cfmt.Printf("    {{Classes}}::green                                        %s\n", colorInt(walkers.GlobalCtx.Classes.CountClasses()))
			cfmt.Printf("        {{Abstract Classes}}::green:                          %s %s\n", colorInt(walkers.GlobalCtx.Classes.CountAbstractClasses()), colorPercent(utils.Percent(walkers.GlobalCtx.Classes.CountAbstractClasses(), int64(walkers.GlobalCtx.Classes.Len()))))
			cfmt.Printf("        {{Concrete Classes}}::green:                          %s %s\n", colorInt(walkers.GlobalCtx.Classes.CountConcreteClasses()), colorPercent(100-utils.Percent(walkers.GlobalCtx.Classes.CountAbstractClasses(), int64(walkers.GlobalCtx.Classes.Len()))))
			cfmt.Printf("    {{Methods}}::green:                                       %s\n", colorInt(walkers.GlobalCtx.Functions.CountMethods()))
			cfmt.Printf("    {{Constants}}::green:                                     %s\n", colorInt(int64(walkers.GlobalCtx.Constants.Len())))
			cfmt.Printf("    {{Functions}}::green:\n")
			cfmt.Printf("        {{Named Functions}}::green:                           %s %s\n", colorInt(walkers.GlobalCtx.Functions.CountFunctions()), colorPercent(utils.Percent(walkers.GlobalCtx.Functions.CountFunctions(), walkers.GlobalCtx.Functions.CountFunctions()+walkers.GlobalCtx.CountAnonymousFunctions)))
			cfmt.Printf("        {{Anonymous Functions}}::green:                       %s %s\n", colorInt(walkers.GlobalCtx.CountAnonymousFunctions), colorPercent(utils.Percent(walkers.GlobalCtx.CountAnonymousFunctions, walkers.GlobalCtx.Functions.CountFunctions()+walkers.GlobalCtx.CountAnonymousFunctions)))

			fmt.Println()
		},
	}

	return briefExecutor
}
