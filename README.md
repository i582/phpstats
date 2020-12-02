![](doc/logo_1.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/i582/phpstats)](https://goreportcard.com/report/github.com/i582/phpstats) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/i582/phpstats/master/LICENSE) ![](https://img.shields.io/badge/-%3C%3E%20with%20%E2%9D%A4-red)

# phpstats

`phpstats` is a tool for **collecting project statistics** and **building dependency graphs** for PHP, that allows you to find places in the code that can be **improved**.

It tries to **be fast**, ~150k LOC/s (*lines of code per second*) on Core i5 with SSD with ~3500Mb/s for reading.

This tool is written in [Go](https://golang.org/) and uses [NoVerify](https://github.com/VKCOM/noverify).

## Table of Contents

* [What is currently available?](#what-is-currently-available)
  * [Metrics](#metrics)
  * [Graphs](#graphs-graphviz-format-and-svg)
  * [Relation](#relation)
  * [Brief project information](#brief-project-information)
* [Install](#install)
* [Usage](#usage)
* [Config](#config)
* [Server](#server)
* [Contact](#contact)
* [Contributing](#contributing)
* [License](#license)

![](doc/screen.svg)

## What is currently available?

### Metrics

1. `Afferent couplings`:
   - for classes;
   - for namespaces;
2. `Efferent couplings`:
   - for classes;
   - for namespaces;
3. `Instability`:
   - for the classes;
   - for namespaces;
4. `Abstractness`;
5. `Lack of Cohesion in Methods`;
6. `Lack of Cohesion in Methods 4`;
7. `Cyclomatic Complexity`;
8. `Count of magic numbers in functions and methods`.

### Graphs (Graphviz format and svg)

1. Class (or interface) dependencies;
2. Class (interface) extend and implementation dependencies;
3. Function or method dependencies;
4. Links within a class (or graph for the LCOM 4 metric);
5. Links between files (included in global and in function);
6. Namespace dependencies graph;
7. Namespace structure graph.

See [building graphs](doc/graphs.md) for details.

### Relation

1. Checking the reachability of a function from another function and outputs the call stacks.

### Brief project information

See [example of brief command](./doc/brief-command-example.md) for details.

## Install

If you don't have the Go toolkit installed, then go to the official [site](https://golang.org/) and install it to continue.

After installation, run the following command in terminal.

```
go get -u -v github.com/i582/phpstats
```

After that you can use by writing `~/go/bin/phpstats` in the terminal.

If you want to work with **dependency graphs**, then you need to install the [Graphviz](https://graphviz.org/download/) utility to visualize graphs.

## Usage

```
$ phpstats collect         \
    [--config-path <dir>]  \
    [--cache-dir <dir>]    \
    [--disable-cache]      \
    [--port <value>]       \
    [--project-path <dir>] \
    [<analyze-dir>]
```

>  All flags and analysis directory are optional.

The `--config-path` flag sets the **path to the configuration file**. See [config](doc/config.md).

The `--cache-dir` flag sets a **custom cache directory**.

The `--disable-cache` flag **disables caching**.

The `--project-path` flag sets the directory relative to which **paths to files will be resolved when importing**. If the flag is not set, the directory is set to the value of the current analyzed directory.

The `--port` flag sets the **port for the server**. See the [server](#Server) part.

The analyzed directory can be omitted if the include field is **specified in the config** (*by default it is* `"./"`).  See [config](doc/config.md).

After collecting information you will be taken to an **interactive shell**, for help, enter `help`.

See [Getting started](doc/getting-start.md) for details.

### Metrics

To **view the metrics**, use the `info` command, which **shows information** about classes, functions, namespaces or files by their names. The **search is not strict**, so it is not necessary to enter the full name.

```
>>> info class ClassName
# show information about ClassName class.
```

For command information, write `info help`.

### Building Graphs

To **build graphs**, use the `graph` command. The `-o` flag is required and sets the file in which the graph will be placed.

```
>>> graph class -o graph.svg ClassName
# outputs the graph for the ClassName class dependencies to the graph.svg file.
```

When creating a graph, two files are created, one with the source code of the graph in the `graphviz` format and a file with the graph in `svg` format.

For command information, write `graph help`.

See [Building graphs](doc/graphs.md) for details.

### Relation

See [Relationships between symbols](doc/relation.md) for details.

### Brief project information

Use the `brief` command to **show brief information about the project**.

```
>>> brief
# shows brief information.
```

See [example of brief command](./doc/brief-command-example.md) for details.

### Brief metrics information

Use the `metrics` command to see a summary of the metrics being collected.

```
>>> metrics		
# shows brief information about of the colected metrics.		
```

## Config

The config allows **more flexible** and **convenient** control over the launch of the analyzer.

More details can be found on the [config page](doc/config.md).

## Server

> Server and API are under development.

A local server (port 8080 by default) is used to **interact with the analyzer from other programs**. The server, **by default**, is started every time an analysis is started.

### API

> All API responses are in `json` format.

`/info/class?name=value` — getting information about the class by its name (the name does not have to be completely the same, the search is not strict).

`/info/func?name=value` — getting information about a function by its name.

`/info/namespace?name=value` — getting information about a namespace by its name.

`/exit` — shutdown of the server.

`/analyzeStats` — getting the current analysis state.

## Contact

 For any questions — tg: `@petr_makhnev`.

## Contributing

Feel free to contribute to this project. I am always glad to new people.

## License

This project is under the **MIT License**. See the [LICENSE](https://github.com/i582/phpstats/blob/master/LICENSE) file for the full license text.
