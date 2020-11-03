![](doc/logo.png)

# phpstats

`phpstats` is a utility for collecting project statistics and building dependency graphs for PHP, that allows you to find places in the code that can be improved.

It tries to be fast, ~150k LOC/s (lines of code per second) on Core i5 with SSD with ~3500Mb/s for reading.

This tool is written in [Go](https://golang.org/) and uses [NoVerify](https://github.com/VKCOM/noverify).

![](doc/screen.png)

## What's currently available?

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
   
4. `Lack of Cohesion in Methods`;
5. `Lack of Cohesion in Methods 4`;
6. `Cyclomatic Complexity`.

### Graph output (Graphviz format and svg)

1. File dependencies, both for root and inside functions;

2. Class dependencies;


3. Function/method dependencies;


4. All project namespaces;


5. Specific namespace;


7. LCOM4.

### Tops

Tops displays information about the top functions, classes and files. The `top` command is used to display the top.

```
>>> top classes
# shows the top 10 classes.
```

To show the list in reverse, add the file `-r`. To control the count and offset in the list, use the `-c` and `-o` flags, respectively.

```
>>> top classes -c 100 -o 10 -r
# shows the top 100 classes from the end, starting from the 10th.
```

Supported output to a file in `json` format, for this add the `--output` flag and the path to the file to which you want to write the list.

```
>>> top classes --output top-classes.json
# outputs the top 10 classes to top-classes.json file.
```

#### Classes

- by Lack of cohesion in methods;
- by Lack of cohesion in methods 4;
- by Afferent coupling;
- by Efferent coupling;
- by Instability;
- by the number of classes on which it depends;
- by the number of classes dependent on it.

#### Functions

- by  the number of classes on which it depends;
- by the number of classes dependent on it;
- by uses count;
- by cyclomatic complexity.

### Brief project information

- Count of classes;
- Count of methods;
- Count of constants;
- Count of functions;
- Count of files;
- Count of lines of code.

## Install

If you don't have the Go toolkit installed, then go to the official [site](https://golang.org/) and install it to continue.

After installation, run the following command in terminal.

```
go get -u -v github.com/i582/phpstats
```

After that you can use it simply by writing `phpstats` in the terminal.

If you want to work with dependency graphs, then you need to install the [Graphviz](https://graphviz.org/download/) utility to visualize graphs.

## Usage

```
$ phpstats collect [--port <value>] [--project-path <dir>] [--cache-dir <dir>] <analyze-dir>
```

The `--project-path` flag sets the directory relative to which paths to files will be resolved when importing. If the flag is not set, the directory is set to the value of the current analyzed directory.

The `--cache-dir` flag sets a custom cache directory.

The `--port` flag sets the port for the server. See the [server](#Server) part.

After collecting information you will be taken to an interactive shell, for help, enter "help".

#### Building graphs

The graph command is used to build graphs. The required flag is the `-o` flag which sets the output file with the graph.

When creating a graph, two files are created, one with the source code of the graph in the graphviz language and a file with the graph in svg format.

For command information, write `graph help`.

Example:

```
>>> graph class -o graph.svg ClassName
```

## Server

> Server and API are under development.

A local server (port 8080 by default) is used to interact with the analyzer from other programs. The server, by default, is started every time an analysis is started.

### API

`/info/class?name=value` — getting information about the class by its name (the name does not have to be completely the same, the search is not strict).

`/info/func?name=value` — getting information about a function by its name.

`/info/namespace?name=value` — getting information about a namespace by its name.

`/exit` — shutdown of the server.

`/analyzeStats` — get the current analysis state.

## License

MIT

---

**<>** with ❤