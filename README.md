![](/doc/logo.png)

# phpstats

phpstats — это небольшая утилита для сбора статистики проектов на PHP.

### Установка

```
go get github.com/i582/phpstats
```

### Использование

```
collect [--project-path <value>] <dir>
```

Флаг `--project-path` устанавливает директорию относительно которой будут разрешаться пути к файлам при импортировании. Если флаг не проставлен, директория устанавливается в значение текущей анализируемой директории.

После сбора информации вы попадете в интерактивную оболочку, для помощи введите `help`.

```
>>> help
Commands:
  clear                clear screen
  exit                 exit program
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

  help                 help page
```

### Roadmap

#### Команда `info`

1. Вывод информации о классе
   * [x] Афферентность (количество классов от которых зависит класс)
     * [x] Учитывать константы
     * [x] Учитывать методы
     * [x] Учитывать статические методы
     * [x] Учитывать использование new
   * [x] Эфферентность (количество классов которые зависит от класса)
     * [x] Учитывать константы
     * [x] Учитывать методы
     * [x] Учитывать статические методы
     * [x] Учитывать использование new
   * [x] Стабильность Эфферентность / (Эфферентность + Афферентность)
   * [x] Расчет LCOM
   * [x] Какие реализует интерфейсы
     * [ ] Выводить рекурсивно?
   * [x] От какого класса наследуется
     * [ ] Выводить рекурсивно?
* [x] Список методов
  
2. Вывод информации о функции/методе

   * [x] Вывод места определения (или информацию о том, что функция встроенная)
   * [x] Вывод количества использований
   * [x] Вывод вызываемых внутри функций
3. Вывод информации о файле

   * [x] Вывод подключаемых файлов в корне
   * [x] Вывод подключаемых файлов в функциях
   * [ ] Вывод классов определенных внутри
   * [ ] Вывод функций определенных внутри

#### Команда `list`

1. Выводить список классов
   * [x] Возможность указывать количество
   * [x] Возможность указывать сдвиг
2. Выводить список интерфейсов
   * [x] Возможность указывать количество
   * [x] Возможность указывать сдвиг
3. Выводить список файлов
   * [x] Возможность указывать количество
   * [x] Возможность указывать сдвиг
4. Выводить список функций
   * [x] Возможность указывать количество
   * [x] Возможность указывать сдвиг
5. Выводить список методов
   * [x] Возможность указывать количество
   * [x] Возможность указывать сдвиг

#### Команда `graph`

1. Вывод информации о файле
    * [x] Вывод зависимостей файла в Graphviz формате
      * [x] Возможность указать максимальный уровень вложенности
      * [x] Разделять подключения в корне и функциях



### Лицензия

MIT