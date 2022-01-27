package localfs

import (
	"io/ioutil"
	"os"
	"testing"
)

func containsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestFindDatasets(t *testing.T) {

	dir, err := ioutil.TempDir("/tmp", "tractor")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	provider, err := NewLocalFSProvider(map[string]interface{}{
		"path": dir,
	})
	if err != nil {
		t.Fatal(err)
	}
	files := []string{"abc.csv", "abcD.csv", "abcDe.csv"}
	for _, file := range files {
		err = ioutil.WriteFile(dir+"/"+file, []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		pattern string
		result  []string
	}{
		{"", files},
		{"^abc", files},
		{"^abcD", []string{"abcD.csv", "abcDe.csv"}},
	}

	for _, test := range tests {
		datasets, err := provider.FindDatasets(test.pattern)
		if err != nil {
			t.Fatal(err)
		}
		if len(datasets) != len(test.result) {
			t.Fatalf("Expected %d datasets, got %d", len(test.result), len(datasets))
		}
		names := []string{}
		for _, dataset := range datasets {
			names = append(names, dataset.Name)
		}
		for _, file := range test.result {
			if !containsStr(names, file) {
				t.Fatalf("Expected dataset %s to be in %v", file, names)
			}
		}
	}
}
