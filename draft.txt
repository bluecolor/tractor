package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"regexp"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/util"
)

func listFiles(dir, pattern string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filteredFiles := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			return nil, err
		}
		if matched {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

func main() {
	configFile := "mappings.yml"
	// util.PrintMappingList(configFile)
	// cmd.Execute()
	pluginsPath := "bin/plugins"

	if _, err := os.Stat(pluginsPath); err != nil {
		fmt.Printf("failed to find plugins path: %v\n", err)
		os.Exit(1)
	}

	plugins, err := listFiles(pluginsPath, `.*.so`)
	if err != nil {
		fmt.Printf("failed to load plugins: %v\n", err)
		os.Exit(1)
	}
	for _, plug := range plugins {
		p, err := plugin.Open(path.Join(pluginsPath, plug.Name()))
		if err != nil {
			fmt.Printf("failed to open plugin %s: %v\n", plug.Name(), err)
			os.Exit(1)
		}
		symbol, err := p.Lookup("Run")
		if err != nil {
			fmt.Printf("failed find Run method in plugin %s: %v\n", plug.Name(), err)
			os.Exit(1)
		}

		mappings := util.GetMappings(configFile)
		symbol.(func(api.Config))(mappings[0]["Demo"].Input)
	}
}
