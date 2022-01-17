package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Field struct {
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	Length    int64  `yaml:"length"`
	Precision int64  `yaml:"precision"`
	Scale     int64  `yaml:"scale"`
}

func (f *Field) String() string {
	return fmt.Sprintf("name:%s, type:%s", f.Name, f.Type)
}

type Catalog struct {
	Name       string  `yaml:"name"`
	Properties []Field `yaml:"fields"`
}

type Input struct {
	Plugin  string                 `yaml:"plugin"`
	Config  map[string]interface{} `yaml:"config"`
	Catalog *Catalog               `yaml:"catalog"`
}

type Output struct {
	Plugin  string                 `yaml:"plugin"`
	Config  map[string]interface{} `yaml:"config"`
	Catalog *Catalog               `yaml:"catalog"`
}

type Mapping struct {
	Name   string `yaml:"name"`
	Input  Input  `yaml:"input"`
	Output Output `yaml:"output"`
}

type Options struct {
	LogLevel string `yaml:"log_level"`
}

type Config struct {
	Options  Options   `yaml:"options"`
	Mappings []Mapping `yaml:"mappings"`
}

func NewConfig() *Config {
	return &Config{
		Options: Options{LogLevel: "info"},
	}
}

func (c *Config) LoadConfig(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}
	return nil
}
