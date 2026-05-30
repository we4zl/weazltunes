//go:build !windows

package player

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewMPVControlUsesPrivateSocketDir(t *testing.T) {
	control, ok := newMPVControl().(*socketControl)
	if !ok {
		t.Fatal("expected socket control")
	}
	if control.tmpDir == "" {
		t.Fatal("expected private temp directory")
	}
	if filepath.Dir(control.path) != control.tmpDir {
		t.Fatalf("socket path %q is not inside temp dir %q", control.path, control.tmpDir)
	}
	info, err := os.Stat(control.tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Fatalf("socket dir mode = %o, want 700", got)
	}

	control.close()
	if _, err := os.Stat(control.tmpDir); !os.IsNotExist(err) {
		t.Fatalf("socket dir still exists after close: %v", err)
	}
}
