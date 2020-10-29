![](/doc/logo.png)

# phpstats

phpstats is a small utility for collecting statistics of PHP projects based on [NoVerify](https://github.com/VKCOM/noverify).

#### Metrics

The following metrics are currently available:

1. Afferent couplings:
   1. for classes;
   2. for namespaces;
2. Efferent couplings:
   1. for classes;
   2. for namespaces;
3. Instability:
   1. for the classes;
   2. for namespaces;
4. Lack of Cohesion in Methods for classes;
5. LCOM4 for classes.

#### Graph output (Graphviz format)

1. File dependencies, both for root and inside functions.

2. Class dependencies.

  ![](/doc/class_graph.svg)

3. Function/method dependencies.

  

  ![](/doc/func_graph.svg)

4. LCOM4

## Install

```
go get github.com/i582/phpstats
```

## Using

```
collect [--server] [--project-path <value>] [--cache-dir <value>] <dir>
```

The `--project-path` flag sets the directory relative to which paths to files will be resolved when importing. If the flag is not set, the directory is set to the value of the current analyzed directory.

The `--cache-dir` flag sets a custom cache directory.

The `--server` flag sets whether the server will be available for queries while the analyzer is running.. See [server](#Server) part.

After collecting information, you will be taken to an interactive shell, type `help` for help.

```
 >>> help
  brief                               shows general information
 
  info                                info about
     class (or interface) <value>     info about class or interface

     func (or method) <value>         info about function or method

     file  <value>                    info about file
       [-f]                           output full information
       [-r <value>]                   output recursive (default: 5)

     namespace  <value>               info about namespace

  graph                               dependencies graph view
     lcom4  <value>                   show lcom4 connected class components
       [-show]                        show graph file in console
        -o <value>                    output file

     file  <value>                    dependency graph for file
       [-block]                       only block require
       [-show]                        show graph file in console
        -o <value>                    output file
       [-r <value>]                   recursive level (default: 5)
       [-root]                        only root require

     class (or interface) <value>     dependency graph for class or interface
        -o <value>                    output file
       [-r <value>]                   recursive level (default: 5)
       [-show]                        show graph file in console

     func (or method) <value>         dependency graph for function or method
        -o <value>                    output file
       [-show]                        show graph file in console

  list                                list of
     funcs                            show list of functions
       [-c <value>]                   count in list (default: 10)
       [-o <value>]                   offset in list (default: 0)
       [-e]                           show embedded functions

     methods                          show list of methods
       [-c <value>]                   count in list (default: 10)
       [-o <value>]                   offset in list (default: 0)

     files                            show list of files
       [-c <value>]                   count in list (default: 10)
       [-o <value>]                   offset in list (default: 0)
       [-f]                           show full information

     classes                          show list of classes
       [-c <value>]                   count in list (default: 10)
       [-o <value>]                   offset in list (default: 0)

     interfaces                       show list of interfaces
       [-c <value>]                   count in list (default: 10)
       [-o <value>]                   offset in list (default: 0)
       [-f]                           show full information

  top                                 shows top of
     funcs                            show top of functions
       [-by-deps]                     top functions by dependencies
       [-by-as-dep]                   top functions by as dependency
       [-by-uses]                     top functions by uses count
       [-r]                           sort reverse
       [-c <value>]                   count in list (default: 10)
       [-o <value>]                   offset in list (default: 0)

     classes                          show top of classes
       [-by-eff]                      top classes by efferent coupling
       [-by-lcom]                     top classes by Lack of cohesion in methods
       [-by-deps]                     top classes by dependencies
       [-r]                           sort reverse
       [-by-aff]                      top classes by afferent coupling
       [-by-instab]                     top classes by instability
       [-by-lcom4]                    top classes by Lack of cohesion in methods 4
       [-by-as-dep]                   top classes by as dependency
       [-c <value>]                   count in list (default: 10)
       [-o <value>]                   offset in list (default: 0)

  help                                help page
  clear                               clear screen
  exit                                exit the program
>>>
```

## Server

> Server and api are under development.

A local server is used to interact with the analyzer from other programs. To enable the server, pass the `--server` flag when starting the analyzer.

The following APIs are currently available:

`/info/class?name=value` — getting information about the class by its name (the name does not have to be completely the same, the search is not strict).

`/info/func?name=value` — getting information about a function by its name.

`/exit` — shutdown of the server.

`/analyzeStats` — get the current analysis state.

## License

MIT