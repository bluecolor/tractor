package localfs

import (
	"io/ioutil"
	"regexp"

	"github.com/bluecolor/tractor/pkg/lib/cat/meta"
)

type LocalFS struct {
	Path string `json:"path"`
}

func NewLocalFS(config map[string]string) (*LocalFS, error) {
	return &LocalFS{
		Path: config["path"],
	}, nil
}
func (f *LocalFS) FindDatasets(pattern string) ([]meta.Dataset, error) {
	files, err := ioutil.ReadDir(f.Path)
	if err != nil {
		return nil, err
	}
	datasets := []meta.Dataset{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if pattern != "" {
			match, _ := regexp.MatchString(pattern, fileName)
			if !match {
				continue
			}
		}
		datasets = append(datasets, meta.Dataset{
			Name: fileName,
		})
	}
	return datasets, nil
}
