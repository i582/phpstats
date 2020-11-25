package config

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Include      []string  `yaml:"include"`
	Exclude      []string  `yaml:"exclude"`
	Port         int64     `yaml:"port"`
	CacheDir     string    `yaml:"cacheDir"`
	DisableCache bool      `yaml:"disableCache"`
	ProjectPath  string    `yaml:"projectPath"`
	UsePackages  bool      `yaml:"use-packages"`
	Packages     *Packages `yaml:"packages"`
	Extensions   []string  `yaml:"extensions"`
}

type Packages []*Package

type Package struct {
	Name       string   `yaml:"name"`
	Namespaces []string `yaml:"namespaces"`
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

func (c *Config) AddPackagesToContext(packages *Packages) {
	if !c.UsePackages {
		return
	}

	*packages = *c.Packages
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

func (p Packages) GetPackage(className string) (*Package, bool) {
	for _, pack := range p {
		for _, namespace := range pack.Namespaces {
			if strings.HasPrefix(className, namespace) {
				return pack, true
			}
		}
	}
	return nil, false
}
