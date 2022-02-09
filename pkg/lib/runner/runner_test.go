package runner

// import (
// 	"testing"

// 	"github.com/bluecolor/tractor/pkg/lib/meta"
// )

// func TestNew(t *testing.T) {
// 	c := meta.Connection{
// 		ConnectionType: "dummy",
// 	}
// 	if _, err := New(c, c); err != nil {
// 		t.Error(err)
// 	}
// }

// func TestRun(t *testing.T) {
// 	c := meta.Connection{
// 		ConnectionType: "dummy",
// 	}
// 	r, err := New(c, c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	p := meta.ExtParams{}
// 	if err := r.Run(p); err != nil {
// 		t.Error("expected no error")
// 	}
// }
