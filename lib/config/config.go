package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Field struct {
	Name      string                 `yaml:"name"`
	Type      string                 `yaml:"type"`
	Length    int64                  `yaml:"length"`
	Precision int64                  `yaml:"precision"`
	Scale     int64                  `yaml:"scale"`
	Source    string                 `yaml:"source"`
	Options   map[string]interface{} `yaml:"options"`
}

func (f *Field) String() string {
	return fmt.Sprintf("name:%s, type:%s", f.Name, f.Type)
}

type Catalog struct {
	Name          string  `yaml:"name"`
	AutoMapFields bool    `yaml:"auto_map_fields"`
	Fields        []Field `yaml:"fields"`
}

func (c *Catalog) GetFieldMap() map[string]*Field {
	fieldMap := make(map[string]*Field)
	for _, f := range c.Fields {
		fieldMap[f.Name] = &f
	}
	return fieldMap
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
func (c *Config) Load(path string) error {
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
	return nil, fmt.Errorf("unable to find mapping [%s]", name)
}
