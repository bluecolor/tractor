package csv

import (
	"io/ioutil"
	"path"
)

func (c *Csv) getFiles() (files []string, err error) {
	fileInfos, err := ioutil.ReadDir(c.Path)
	if err != nil {
		return nil, err
	}
	for _, f := range fileInfos {
		if !f.IsDir() {
			files = append(files, path.Join(c.Path, f.Name()))
		}
	}
	// fmt.Printf("%v\n", files)
	return files, nil
}
