package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"sync"

	"github.com/bluecolor/tractor/api"
	"github.com/bluecolor/tractor/api/message"
	c "github.com/bluecolor/tractor/util/constants"
	"github.com/spf13/viper"
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

// GetMappingPlugins ...
func GetMappingPlugins(mapping *Mapping, args ...interface{}) (*api.TractorPlugin, *api.TractorPlugin, error) {
	var path string
	if len(args) > 0 {
		path = args[0].(string)
	} else {
		path = viper.GetString(c.TractorPluginsPath)
	}
	iname := mapping.Input["plugin"].(string)
	oname := mapping.Output["plugin"].(string)
	return GetPlugins(path, iname, oname)
}

// GetPlugins Get input and output plugins
func GetPlugins(pluginsPath, inputPluginName, outputPluginName string) (*api.TractorPlugin, *api.TractorPlugin, error) {

	iplugin, err := getPlugin(pluginsPath, inputPluginName, c.MappingInputKey)
	if err != nil {
		return nil, nil, err
	}

	oplugin, err := getPlugin(pluginsPath, outputPluginName, c.MappingOutputKey)
	if err != nil {
		return nil, nil, err
	}

	ip := api.TractorPlugin{Plugin: iplugin}
	op := api.TractorPlugin{Plugin: oplugin}

	if err := setRun(&ip, iplugin); err != nil {
		return nil, nil, err
	}
	if err := setRun(&op, oplugin); err != nil {
		return nil, nil, err
	}

	if err != nil {
		fmt.Printf("Failed find Run method in input plugin %s: %v\n", iplugin, err)
		os.Exit(1)
	}

	return &ip, &op, nil
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

func setRun(p *api.TractorPlugin, plug *plugin.Plugin) error {
	symbol, err := plug.Lookup(c.PluginRunMethodName)
	if err != nil {
		return err
	}
	p.Run = symbol.(func(*sync.WaitGroup, []byte, chan *message.Message) error)
	return nil
}
