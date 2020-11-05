# Config

The phpstats config uses the [yaml](https://cloudslang-docs.readthedocs.io/en/v1.0/yaml_overview.html) markup language.

By default, the config looks like this.

```yaml
# The port on which the server will be launched
# to interact with the analyzer from other programs.
# By default, it is 8080
# port: 8080

# The path where the cache will be stored.
# Caching can significantly speed up data collection.
# By default, it is set to the value of the temporary folder + /phpstats.
# cacheDir: ""

# Disables caching.
# By default, it is false
# disableCache: false

# Path to the project relative to which all imports are allowed.
# By default, it is equal to the analyzed directory.
# projectPath: ""

# File extensions to be included in the analysis.
# By default, it is php, inc, php5, phtml.
# extensions:
#  - "php"
#  - "inc"
#  - "php5"
#  - "phtml"
```

All you need to do is uncomment the required field and write the required value into it.

