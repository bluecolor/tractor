package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Property struct {
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	Length    int64  `yaml:"length"`
	Precision int64  `yaml:"precision"`
	Scale     int64  `yaml:"scale"`
}

type Catalog struct {
	Name       string     `yaml:"name"`
	Properties []Property `yaml:"properties"`
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

type Settings struct {
	LogLevel string `yaml:"log_level"`
}

type Config struct {
	Settings Settings  `yaml:"settings"`
	Mappings []Mapping `yaml:"mappings"`
}

func NewConfig() *Config {
	return &Config{
		Settings: Settings{LogLevel: "info"},
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

func (c *Config) GetMapping(name string) (*Mapping, error) {
	for _, m := range c.Mappings {
		if m.Name == name {
			return &m, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Unable to find mapping [%s]", name))
}
