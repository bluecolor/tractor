package util

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"plugin"
)

func findAllPluginFiles(dir string) ([]os.FileInfo, error) {
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := []os.FileInfo{}
	for _, file := range items {
		if file.IsDir() {
			continue
		}
		files = append(files, file)
	}
	return files, nil
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

// GetPluginNamesByType ...
func GetPluginNamesByType(pluginsPath, ptype string) ([]string, error) {
	ppath := path.Join(pluginsPath, ptype)
	files, err := findAllPluginFiles(ppath)
	if err != nil {
		return nil, err
	}
	plugins := []string{}
	for _, file := range files {
		plugins = append(plugins, file.Name())
	}
	return plugins, nil
}
