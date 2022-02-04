package runner

import (
	"testing"

	"github.com/bluecolor/tractor/pkg/lib/meta"
)

func TestNew(t *testing.T) {
	r := New(meta.ExtParams{})
	if r == nil {
		t.Error("expected a new runner")
	}
}
