package conf

import (
	"os"
	"testing"

	"github.com/bluecolor/tractor/pkg/test"
)

func TestLoadConf(t *testing.T) {
	envfile, dirname, err := test.GenEnvFile()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dirname)
	})
	if _, err := LoadConfig(envfile); err != nil {
		t.Error(err)
	}
}
