package main

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// ErrNoConfigFile -
var ErrNoConfigFile = errors.New("config File does not exists")

// Config -
type Config struct {
	Version      string   `yaml:"version"`
	Kind         string   `yaml:"kind"`
	Host         string   `yaml:"host"`
	FaviconPath  string   `yaml:"faviconPath"`
	EntryPoint   string   `yaml:"entryPoint"`
	TemplatePath string   `yaml:"templatePath"`
	AssetsPath   string   `yaml:"assetsPath"`
	MarkdownPath string   `yaml:"markdownPath"`
	BindJS       bool     `yaml:"bindJS"`
	BindCSS      bool     `yaml:"bindCSS"`
	BindFavicon  bool     `yaml:"bindFavicon"`
	CSS          []string `yaml:"CSS"`
	JS           []string `yaml:"JS"`
}

// ParseConfig -
func (c *Config) ParseConfig(configFile string) error {
	congNotExist := func(filename string) bool {
		_, err := os.Stat(filename)
		return os.IsNotExist(err)
	}
	if congNotExist(configFile) {
		return ErrNoConfigFile
	}
	b, err := ioutil.ReadFile(configFile)
	err = yaml.Unmarshal(b, c)
	return err
}
