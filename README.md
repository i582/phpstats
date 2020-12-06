![](doc/logo.png)

![Build Status](https://github.com/i582/phpstats/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/i582/phpstats)](https://goreportcard.com/report/github.com/i582/phpstats) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/i582/phpstats/master/LICENSE) ![](https://img.shields.io/badge/-%3C%3E%20with%20%E2%9D%A4-red)

# phpstats

**PhpStats** is a tool that *collects statistics* for the code of your project and, based on these statistics, *calculates various qualitative metrics* of the code, *builds the necessary graphs*, and also *finds the relationships between symbols* in the system. 

It tries to be fast, at the moment—about **150k lines of code per second** on a MacBook Pro 2019 with Core i5.

The tool is built on top of [NoVerify](https://github.com/VKCOM/noverify) and written in [Go](https://golang.org/).

You can find the **documentation** for the project [here](https://i582.github.io/phpstats-docs/).

## Table of Contents

* [What is supported?](#what-is-supported)
  * [Metrics](#code-metrics)
  * [Graphs](#dependency-graphs)
  * [Relations](#relations-between-symbols)
  * [Brief project information](#brief-information-about-the-project)
* [About the project](#about-the-project)
* [Contacts](#contacts)
* [Contributing](#contributing)
* [License](#license)

![](doc/screen.svg)

## What is supported?

PhpStats currently represents five areas:

1. Collecting code **metrics**;
2. Building **dependency graphs**;
3. Analysis of **relationships between symbols**;
4. Gathering **brief information** about the project;
5. Analysis of the **reachability** of a function.

It also allows you to **view lists** of *classes, interfaces, functions, methods, files and namespaces* in a **tabular form** with the **ability to sort by metrics**.

*Let's look at each point separately.*

### Code metrics

PhpStats currently calculates the following metrics:

1. Afferent couplings:
2. Efferent couplings:
3. Instability:
4. Abstractness;
5. Lack of Cohesion in Methods;
6. Lack of Cohesion in Methods 4 (*or the number of connected components of the class*);
7. Cyclomatic Complexity;
8. Count of magic numbers in functions and methods;
9. Count fully typed methods.

### Dependency graphs

PhpStats is currently building the following dependency graphs:

1. Class (or interface) dependencies;
2. Class (interface) extend and implementation dependencies;
3. Function or method dependencies;
4. Links within a class (*or graph for the LCOM 4 metric*);
5. Links between files (*included in global and in functions*);
6. Namespace dependencies graph;
7. Namespace structure graph;
8. Function reachability graph.

[Graphviz](https://graphviz.org/download/) is used to build graphs.

### Relations between symbols

PhpStats is currently analyzing the following relations:

1. **For class-class relations:**
   1. Whether one class is **extends** another and vice versa;
   2. Whether the class **implements** the interface or vice versa;
   3. What methods, fields and constants are **used by one class used by another** and in which methods this happens.

2. **For class-function relations:**
   1. Function **belong** to class;
   2. The class is **used inside** the function;
   3. **Used class** members in functions;
   4. The function is **used in the class** (*+ all methods where this function is used*).

3. **For function-function relations:**
   1. Functions **belong to the same class**;
   2. Does the **first function use the second** and vice versa;
   3. Whether the **first function is reachable from the second through calls** and vice versa (*+ call stacks to reach the function*).

### Brief information about the project

PhpStats is currently collecting various **brief** information about the project presented below:

```
General Test project statistics

Size
    Lines of Code (LOC):                             611240
    Comment Lines of Code (CLOC):                    109340 (17.89%)
    Non-Comment Lines of Code (NCLOC):               501900 (82.11%)

Metrics
    Cyclomatic Complexity
        Average Complexity per Class:                  5.55
            Maximum Class Complexity:              29954.00
            Minimum Class Complexity:                  0.00
        Average Complexity per Method:                 1.01
            Maximum Method Complexity:               142.00
            Minimum Method Complexity:                 0.00
        Average Complexity per Functions:              0.00
            Maximum Functions Complexity:              6.00
            Minimum Functions Complexity:              0.00

    Count of Magic Numbers
        Average Class Count:                              0
            Maximum Class Count:                       5055
            Minimum Class Count:                          2
        Average Method Count:                             0
            Maximum Method Count:                        50
            Minimum Method Count:                         0
        Average Functions Count:                          0
            Maximum Method Count:                         2
            Minimum Method Count:                         0

Structure
    Files:                                             5323
    Namespaces:                                        1680
    Interfaces:                                         423
    Traits                                               10
    Classes                                            4974
        Abstract Classes:                               218 (4.04%)
        Concrete Classes:                              4756 (95.96%)
    Methods:                                          29738
    Constants:                                         1152
    Functions:
        Named Functions:                                 66 (3.32%)
        Anonymous Functions:                           1921 (96.68%)
```

## About the project

PhpStats is © 2020-2020 by Petr Makhnev.

### Contacts

Have any questions—welcome in telegram: [@petr_makhnev](https://t.me/petr_makhnev).

### Contributing

Feel free to contribute to this project. I am always glad to new people.

### License

PhpStats is distributed by an [MIT license](https://github.com/i582/phpstats/tree/master/LICENSE).