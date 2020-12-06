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
  * [Reachability](#reachability-of-functions)
  * [Brief project information](#brief-information-about-the-project)
* [About the project](#about-the-project)
* [Contacts](#contacts)
* [Contributing](#contributing)
* [License](#license)

![](doc/screen.svg)

## What is supported?

**PhpStats** currently represents five areas:

1. Collecting code **metrics**;
2. Building **dependency graphs**;
3. Analysis of **relationships between symbols**;
4. Gathering **brief information** about the project;
5. Analysis of the **reachability** of a function.

It also allows you to **view lists** of *classes, interfaces, functions, methods, files and namespaces* in a **tabular form** with the **ability to sort by metrics**.

*Let's look at each point separately.*

### Code metrics

**PhpStats** currently calculates the following metrics:

1. Afferent couplings:
2. Efferent couplings:
3. Instability:
4. Abstractness;
5. Lack of Cohesion in Methods;
6. Lack of Cohesion in Methods 4 (*or the number of connected components of the class*);
7. Cyclomatic Complexity;
8. Count of magic numbers in functions and methods;
9. Count fully typed methods.

See the documentation [part](https://i582.github.io/phpstats-docs/docs/capabilities/metrics/) for details.

### Dependency graphs

**PhpStats** is currently building the following dependency graphs:

1. Class (or interface) dependencies;
2. Class (interface) extend or implementation dependencies;
3. Function or method dependencies;
4. Links within a class (*or graph for the LCOM 4 metric*);
5. Links between files (*included in global and in functions*);
6. Namespace dependencies graph;
7. Namespace structure graph;
8. Function reachability graph.

See the documentation [part](https://i582.github.io/phpstats-docs/docs/capabilities/graphs/) for details.

### Relations between symbols

**PhpStats** is currently analyzing the following relations:

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

See the documentation [part](https://i582.github.io/phpstats-docs/docs/capabilities/relations/) for details.

### Reachability of functions

See the documentation [part](https://i582.github.io/phpstats-docs/docs/capabilities/function_reachability/) for details.

### Brief information about the project

See the documentation [part](https://i582.github.io/phpstats-docs/docs/capabilities/brief-information/) for details.

## About the project

**PhpStats** is © 2020-2020 by Petr Makhnev.

### Contacts

Have any questions—welcome in telegram: [@petr_makhnev](https://t.me/petr_makhnev).

### Contributing

Feel free to contribute to this project. I am always glad to new people.

### License

PhpStats is distributed by an [MIT license](https://github.com/i582/phpstats/tree/master/LICENSE).