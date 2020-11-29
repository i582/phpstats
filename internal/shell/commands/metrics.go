package commands

import (
	"fmt"

	"github.com/i582/phpstats/internal/shell"
	"github.com/i582/phpstats/internal/shell/flags"
)

func Metrics() *shell.Executor {
	metricsExecutor := &shell.Executor{
		Name:  "metrics",
		Help:  "shows referential information about the metrics being collected",
		Flags: flags.NewFlags(),
		Func: func(c *shell.Context) {
			fmt.Print(`A brief description of the metrics.

Afferent couplings (Ca): 
	The number of classes in other packages that depend upon classes within 
	the package is an indicator of the package's responsibility. 

Efferent couplings (Ce): 
	The number of classes in other packages that the classes in a package 
	depend upon is an indicator of the package's dependence on externalities.

Instability (I): 
	The ratio of efferent coupling (Ce) to total coupling (Ce + Ca) such that 
	I = Ce / (Ce + Ca). 
	This metric is an indicator of the package's resilience to change. 
	The range for this metric is 0 to 1, with I=0 indicating a completely stable 
	package and I=1 indicating a completely unstable package.

Abstractness (A):
	The ratio of the number of abstract classes in a group to the total number of classes.
	A = nA / nAll.
	nA   - the number of abstract classes in a group.
	nAll - the total number of classes.

	0 = the category is completely concrete.
	1 = the category is completely abstract.

Lack of Cohesion in Methods (LCOM):
	The result of subtracting from one the sum of the number of methods (CM_i) 
	that refer to a certain class field (i) for all fields, divided by the number 
	of methods (CM) multiplied by the number of fields (CF).

	LCOM = 1 - (\Sum{i eq [0, CF]}{CM_i}) / (CM * CF))

Lack of Cohesion in Methods 4 (LCOM4):
	The number of "connected components" in a class.
	A connected component is a set of related methods (and class-level variables). 
	There should be only one such a component in each class. 
	If there are 2 or more components, the class should be split into so many smaller classes.

	Which methods are related? Methods a and b are related if:
	
	  - they both access the same class-level variable, or
	  - a calls b, or b calls a.

	After determining the related methods, we draw a graph linking the related methods 
	(use the 'graph lcom4' command to build the graph) to each other. 
	LCOM4 equals the number of connected groups of methods.

	  - LCOM4=1  indicates a cohesive class, which is the "good" class.
	  - LCOM4>=2 indicates a problem. The class should be split into so many smaller classes.
	  - LCOM4=0  happens when there are no methods in a class. This is also a "bad" class.

	Information from https://www.aivosto.com/project/help/pm-oo-cohesion.html#LCOM4

Cyclomatic complexity (CC):
	The number of decision points.
	Cyclomatic complexity is basically a metric to figure out areas of code that needs 
	more attention for the maintainability. It would be basically an input to the refactoring. 
	It definitely gives an indication of code improvement area in terms of avoiding deep 
	nested loop, conditions etc.
	
	The decision points is conditional statements like if, for, while, foreach, case, default, 
	continue, break, goto, catch, ternary. coalesce, or, and.

Count of magic numbers (CMN):
	Magic numbers are any number in code that isn't immediately obvious to someone 
	with very little knowledge.

	Code with magic numbers is more difficult to understand and refactor, as it is not 
	always obvious what the author meant by it. The more magic numbers, the more difficult 
	it is to refactor the given code.

PHPStats (c) 2020
`)
		},
	}

	return metricsExecutor
}
