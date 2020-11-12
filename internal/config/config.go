package config

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Include      []string `yaml:"include"`
	Exclude      []string `yaml:"exclude"`
	Port         int64    `yaml:"port"`
	CacheDir     string   `yaml:"cacheDir"`
	DisableCache bool     `yaml:"disableCache"`
	ProjectPath  string   `yaml:"projectPath"`
	Groups       []Group  `yaml:"groups"`
	Extensions   []string `yaml:"extensions"`
}

type Group struct {
	Name      string
	Namespace string
}

func OpenConfig(path string) (*Config, error, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err, nil
	}
	var config *Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, nil, err
	}

	return config, nil, nil
}

func (c *Config) ToCliArgs() []string {
	var args []string

	if c.DisableCache {
		args = append(args, "-disable-cache")
	}

	if len(c.Exclude) != 0 {
		excludeRegexp := strings.Join(c.Exclude, ",")
		args = append(args, "-index-only-files", excludeRegexp)
	}

	if len(c.Extensions) != 0 {
		extensionsList := strings.Join(c.Extensions, ",")
		args = append(args, "-php-extensions", extensionsList)
	}

	args = append(args, []string{
		"-cache-dir", c.CacheDir,
	}...)

	if len(c.Include) != 0 {
		args = append(args, c.Include...)
	}

	return args
}
