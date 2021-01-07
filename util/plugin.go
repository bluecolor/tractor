package util

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"regexp"
)

// List files with the given pattern
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

func findPluginFile(dir, name string) (os.FileInfo, error) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if name+".so" == file.Name() {
			return file, nil
		}
	}
	return nil, errors.New("Filed to find file " + name + ".so  under " + dir)
}

func getPlugin(pluginsPath, name, ptype string) (*plugin.Plugin, error) {
	ppath := path.Join(pluginsPath, ptype)
	file, err := findPluginFile(ppath, name)
	if err != nil {
		return nil, err
	}
	return plugin.Open(path.Join(ppath, file.Name()))

}

// GetPlugins Get input and output plugins
func GetPlugins(pluginsPath, inputPluginName, outputPluginName string) (*plugin.Plugin, *plugin.Plugin, error) {

	inputPlugin, err := getPlugin(pluginsPath, inputPluginName, "input")
	if err != nil {
		return nil, nil, err
	}

	outputPlugin, err := getPlugin(pluginsPath, outputPluginName, "output")
	if err != nil {
		return nil, nil, err
	}
	return inputPlugin, outputPlugin, nil
}
