![](/doc/logo.png)

# phpstats

phpstats is a small utility for collecting statistics of PHP projects based on [NoVerify](https://github.com/VKCOM/noverify).

#### Metrics

The following metrics are currently available:

1. Afferent couplings
   1. for classes
   2. for namespaces
2. Efferent couplings
   1. for classes
   2. for namespaces
3. Stability
   1. for the classes
   2. for namespaces
4. Lack of Cohesion in Methods for classes

#### Graph output (Graphviz format)

1. File dependencies, both for root and inside functions.
2. Class dependencies.
3. Function/method dependencies.

## Install

```
go get github.com/i582/phpstats
```

## Using

```
collect [--project-path <value>] <dir>
```

The `--project-path` flag sets the directory relative to which paths to files will be resolved when importing. If the flag is not set, the directory is set to the value of the current analyzed directory.

After collecting information, you will be taken to an interactive shell, type `help` for help.

```
>>> help
Commands:
  info                 info about
     class <value>     info about class
       [-f]            output full information

     func <value>      info about function
       [-f]            output full information

     file <value>      info about file
       [-f]            output full information
       [-r <value>]    output recursive (default: 5)

     namespace <value> info about namespace

  list                 list of
     interfaces        show list of interfaces
       [-o <value>]    offset in list (default: 0)
       [-f]            show full information
       [-c <value>]    count in list (default: 10)

     funcs             show list of functions
       [-o <value>]    offset in list (default: 0)
       [-e]            show embedded functions
       [-c <value>]    count in list (default: 10)

     methods           show list of methods
       [-c <value>]    count in list (default: 10)
       [-o <value>]    offset in list (default: 0)

     files             show list of files
       [-c <value>]    count in list (default: 10)
       [-o <value>]    offset in list (default: 0)
       [-f]            show full information

     classes           show list of classes
       [-c <value>]    count in list (default: 10)
       [-o <value>]    offset in list (default: 0)
       [-f]            show full information

  graph                dependencies graph view
     file <value>      dependency graph for file
       [-show]         show graph file in console
        -o <value>     output file
       [-r <value>]    recursive level (default: 5)
       [-root]         only root require
       [-block]        only block require

     class <value>     dependency graph for class
        -o <value>     output file
       [-r <value>]    recursive level (default: 5)
       [-show]         show graph file in console

     func <value>      dependency graph for function
       [-r <value>]    recursive level (default: 5)
       [-show]         show graph file in console
        -o <value>     output file
        
  brief                shows general information
  clear                clear screen
  exit                 exit program
  help                 help page
```

## License

MIT