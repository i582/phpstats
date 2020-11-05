package config

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port         int64    `yaml:"port"`
	CacheDir     string   `yaml:"cacheDir"`
	DisableCache bool     `yaml:"disableCache"`
	ProjectPath  string   `yaml:"projectPath"`
	Exclude      []string `yaml:"exclude"`
	Groups       []Group  `yaml:"groups"`
	Extensions   []string `yaml:"extensions"`
}

type Group struct {
	Name      string
	Namespace string
}

func OpenConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config *Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) ToCliArgs() []string {
	var args []string

	if c.DisableCache {
		args = append(args, "-disable-cache")
	}

	if len(c.Exclude) != 0 {
		excludeRegexp := strings.Join(c.Exclude, "|")
		args = append(args, "-exclude", excludeRegexp)
	}

	if len(c.Extensions) != 0 {
		extensionsList := strings.Join(c.Extensions, ",")
		args = append(args, "-php-extensions", extensionsList)
	}

	return append(args, []string{
		"-cache-dir", c.CacheDir,
	}...)
}
