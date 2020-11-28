# Getting started

Welcome to the article on getting started with **phpstats**.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Start of work](#start-of-work)
- [Commands in the interactive shell](#commands-in-the-interactive-shell)
  - [info](#info)
  - [list](#list)
  - [graph](#graph)
  - [relation](#relation)
  - [brief](#brief)
  - [metrics](#metrics)
  - [about](#about)
- [Additional commands](#additional-commands)

## Installation

The first step is to install the Go toolkit if you don't already have it. To do this, go to the official [site](https://golang.org/dl/) and download the required package for your system.

After installing the downloaded package, run the following command:

```
go get -u -v github.com/i582/phpstats
```

After that, to run **phpstats**, just write the following in the terminal:

```
$ ~/go/bin/phpstats
```

If you want to build dependency graphs, then you must also install [Graphviz](https://graphviz.org/download/). After installation, make sure the path to `graphviz` is in the `Path` environment variable.

## Configuration

After installing **phpstats**, you need to create a configuration file for the analyzed project. This part of the setup needs to be done for each project that you want to analyze.

To do this, first, go to the root folder of the project you want to analyze.

Then, in the terminal, write the following command:

```
$ ~/go/bin/phpstats init
```

Answer the wizard's questions. Please note that when you enter the path to the source code, the correctness of the path is checked, so adding a path that does not exist will not allowed.

After that, a configuration file will appear in the project folder.

This completes the **phpstats** configuration.

## Start of work

The `collect` subcommand is used to start the analysis. Write the following in the terminal to run:

```
$ ~/go/bin/phpstats collect
```

If everything is configured correctly, then the analysis of the project will begin and, after it is over, you will receive the following:

```
Started
Indexing [./tests]
Linting
11 / 11 [-------------------------------------------------------------------------------------------------------] 100.00% ? p/s
Entering interactive mode (type "help" for commands)
>>>
```

After completing the analysis, you will be taken to an interactive shell. This interactive shell is used for all data interactions. For all available commands, write `help`.

In the next section, we'll take a look at each command separately.

### Commands in the interactive shell

#### `info`

The `info` command is used to get information about classes, files, functions and namespaces.

The command accepts the following subcommands:

1. `class` or `interface` — information about the class or interface;
2. `func` or `method` — information about a function or method;
3. `namespace` — information about the namespace;
4. `file` — information about the file.

For each subcommand, you need to pass the required name, be it a class or a function or something else.
Please note that the search is **not strict**, so it is not necessary to enter the full name, however, if several options are suitable for the entered name, then you will receive the first of them, to get the desired one you need to specify the name.

##### Example

```
>>> info class Foo
```

Will show information about the `Foo` class:

```
>>> info class Foo
Show info about Foo class

Class \Foo
  File:                          /path/to/file/with/foo/foo.php
  Afferent coupling:             0.00
  Efferent coupling:             3.00
  Instability:                   1.00
  Lack of Cohesion in Methods:   -1.00
  Lack of Cohesion in Methods 4: 1
  Count class dependencies:      3
  Count dependent classes:       0
```

#### `list`

The `list` command is used to list classes, files, functions, and namespaces.

The command accepts the following subcommands:

1. `class` — list of classes;
2. `interface` or `ifaces` — list of interfaces;
3. `func` — list of functions;
4. `method` — list of methods;
5. `namespace` — list of namespaces;
6. `file` — list of files.

For each subcommand, you need to pass the required name, be it a class or a function or something else.
Please note that the search is **not strict**, so it is not necessary to enter the full name, however, if several options are suitable for the entered name, then you will receive the first of them, to get the desired one you need to specify the name.

Also, each of the subcommands accepts the following flags:

1. `-c` — number of classes in the list;
2. `-o` — shift of the list from the beginning;
3. `--json` — path to the file for outputting information in the `json` format.
4. `--sort` — column number by which sorting will be performed;
5. `-r` — sort in reverse order.

The `func` subcommand also accepts the following flag:

1. `-e` — flag, when built-in functions will be displayed in the list.

The `namespace` subcommand also accepts the following flag:

1. `-l` — the level of namespaces to be displayed (default: 0 (top-level namespaces)).

##### Example

```
>>> list classes -c 5 -o 2
```

Shows a list of classes of 5 elements starting with the 3rd:

```
>>> list classes -c 5 -o 2
 #                      Name                      Aff     Eff    Instab   LCOM    LCOM 4   Class   Classes
                                                  coup   coup                              deps    depends
--- -------------------------------------------- ------ ------- -------- ------- -------- ------- ---------
 3   \Symfony\Component\DependencyInjection\      0.00   54.00     1.00   undef       55      54         0
     Tests\Compiler\AutowirePassTest
 4   \Symfony\Bundle\FrameworkBundle\             5.00   53.00     0.91   undef        2      53         5
     FrameworkBundle
 5   \Symfony\Component\Console\Tests\            0.00   48.00     1.00    0.95       79      48         0
     ApplicationTest
 6   \Symfony\Bundle\FrameworkBundle\Tests\       3.00   46.00     0.94    0.99        1      46         3
     DependencyInjection\FrameworkExtensionTest
 7   \Symfony\Component\Messenger\Tests\          0.00   44.00     1.00   undef        1      44         0
     DependencyInjection\MessengerPassTest
```

And the command:

```
>>> list classes -c 5 -o 2 --json classes.json
```

Will output information to the file `classes.json`:

```
[
	{
		"name": "\\Symfony\\Component\\DependencyInjection\\Tests\\Compiler\\AutowirePassTest",
		"file": "/path/to/class",
		"type": "Class",
		"aff": 0,
		"eff": 54,
		"instab": 1,
		"lcom": -1,
		"lcom4": 55,
		"countDeps": 54,
		"countDepsBy": 0
	},
	...
]
```

#### `graph`

The `graph` command is used to create dependency graphs for classes, files, functions, and namespaces.

See [building graphs](./graphs.md) for details.

#### `relation`

The `relation` command is used to get information about the relationship between functions (currently only the reachability of one function from another).

The command accepts the following subcommands:

1. `funcs` — the reachability of one function from another.

The `funcs` subcommand accepts the following flags:

1. `--parent` — function from which the reachability of another function will be checked.;
2. `--child` — function to find reachability..

##### Example

```
>>> relation funcs --parent \AD::ADMethod --child \AC::ACMethod
```

Display information about the reachability of the `\AC::ACMethod` method from the `\AD::ADMethod` method:

```
>>> relation funcs --parent \AD::ADMethod --child \AC::ACMethod
Reachability: true

Callstacks:
[\AD::ADMethod -> \AA::AAMethod -> \AC::ACMethod]
```

#### `brief`

The `brief` command is used to view brief information about the project.

##### Example

```
>>> brief
```

Will display brief information about the project:

```
>>> brief
General project statistics

Size
    Lines of Code (LOC):                                236
    Comment Lines of Code (CLOC):                         2 (0.85%)
    Non-Comment Lines of Code (NCLOC):                  234 (99.15%)

Metrics
    Cyclomatic Complexity
        Average Complexity per Class:                  0.00
            Maximum Class Complexity:                  0.00
            Minimum Class Complexity:                  0.00
        Average Complexity per Method:                 0.00
            Maximum Method Complexity:                 0.00
            Minimum Method Complexity:                 0.00
        Average Complexity per Functions:              0.00
            Maximum Functions Complexity:              0.00
            Minimum Functions Complexity:              0.00

    Count of Magic Numbers
        Average Class Count:                              0
            Maximum Class Count:                          4
            Minimum Class Count:                          0
        Average Method Count:                             0
            Maximum Method Count:                         2
            Minimum Method Count:                         0
        Average Functions Count:                          0
            Maximum Method Count:                         0
            Minimum Method Count:                         0

Structure
    Files:                                               11
    Namespaces:                                           0
    Interfaces:                                           2
    Classes                                              29
        Abstract Classes:                                 1 (3.23%)
        Concrete Classes:                                28 (96.77%)
    Methods:                                             17
    Constants:                                            8
    Functions:
        Named Functions:                                  3 (100.00%)
        Anonymous Functions:                              0 (0.00%)
```

#### `metrics`

The `metrics` command is used to view general information about the metrics being collected.

##### Example

```
>>> metrics
```

Will display information about the collected metrics:

```
>>> metrics
A brief description of the metrics.

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
        (Use the 'graph lcom4' command to build the graph.) to each other.
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
```

#### `about`

The `about` command is used to show general information about phpstats.

### Additional commands

Also, in addition to the commands above, there are commands for interacting with the interactive shell.

1. `help` — shows a page with all valid commands;
2. `clear` — clears the console;
3. `exit` — exits the interactive environment.



### Afterword

What's next? If you haven't looked at how to [build graphs](./graphs.md) yet, then it's time to see.