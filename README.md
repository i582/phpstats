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
3. Stability:
   1. for the classes;
   2. for namespaces;
4. Lack of Cohesion in Methods for classes;
5. LCOM4 for classes.

#### Graph output (Graphviz format)

1. File dependencies, both for root and inside functions.
2. Class dependencies.
3. Function/method dependencies.
4. Lcom4

## Install

```
go get github.com/i582/phpstats
```

## Using

```
collect [--project-path <value>] [--cache-dir <value>] <dir>
```

The `--project-path` flag sets the directory relative to which paths to files will be resolved when importing. If the flag is not set, the directory is set to the value of the current analyzed directory.
The `--cache-dir` flag sets a custom cache directory.

After collecting information, you will be taken to an interactive shell, type `help` for help.

```
>>> help
  info                                info about
     class (or interface) <value>     info about class or interface
       [-f]                           output full information
       [-metrics]                     output only metrics

     file  <value>                    info about file
       [-r <value>]                   output recursive (default: 5)
       [-f]                           output full information

     func (or method) <value>         info about function or method
       [-f]                           output full information

     namespace  <value>               info about namespace

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
       [-o <value>]                   offset in list (default: 0)
       [-f]                           show full information
       [-c <value>]                   count in list (default: 10)

     interfaces                       show list of interfaces
       [-o <value>]                   offset in list (default: 0)
       [-f]                           show full information
       [-c <value>]                   count in list (default: 10)

  graph                               dependencies graph view
     class  <value>                   dependency graph for class
        -o <value>                    output file
       [-r <value>]                   recursive level (default: 5)
       [-show]                        show graph file in console

     func  <value>                    dependency graph for function
       [-show]                        show graph file in console
        -o <value>                    output file
       [-r <value>]                   recursive level (default: 5)

     file  <value>                    dependency graph for file
        -o <value>                    output file
       [-r <value>]                   recursive level (default: 5)
       [-root]                        only root require
       [-block]                       only block require
       [-show]                        show graph file in console

     lcom4  <value>                   show lcom4 connected class components
        -o <value>                    output file
       [-show]                        show graph file in console

  top                                 shows top of
     funcs                            show top of functions
       [-by-as-dep <value>]           top functions by as dependency (default: 10)
       [-by-uses <value>]             top functions by uses count (default: 10)
       [-r]                           sort reverse
       [-c <value>]                   count in list (default: 10)
       [-o <value>]                   offset in list (default: 0)
       [-by-deps <value>]             top functions by dependencies (default: 10)

     classes                          show top of classes
       [-by-eff <value>]              top classes by efferent coupling (default: 10)
       [-by-stab <value>]             top classes by stability (default: 10)
       [-by-lcom <value>]             top classes by Lack of cohesion in methods (default: 10)
       [-r]                           sort reverse
       [-by-aff <value>]              top classes by afferent coupling (default: 10)
       [-by-lcom4 <value>]            top classes by Lack of cohesion in methods 4 (default: 10)
       [-by-deps <value>]             top classes by dependencies (default: 10)
       [-by-as-dep <value>]           top classes by as dependency (default: 10)
       [-c <value>]                   count in list (default: 10)
       [-o <value>]                   offset in list (default: 0)

  brief                               shows general information
  help                                help page
  clear                               clear screen
  exit                                exit the program

```

## License

MIT