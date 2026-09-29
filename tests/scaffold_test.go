//go:build !wasm

package html_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// webtyp.com/app copies scaffold/client.go into every new project as
// web/client.go. It must compile against this module's API, so an API change
// that breaks it fails here, in the same commit, not in a user's project.
func TestScaffoldClientCompiles(t *testing.T) {
	if _, err := os.Stat(filepath.Join("..", "scaffold", "client.go")); err != nil {
		t.Fatalf("scaffold/client.go is consumed by webtyp.com/app and must exist: %v", err)
	}

	cmd := exec.Command("go", "build", "-o", filepath.Join(t.TempDir(), "client.wasm"), "./scaffold")
	cmd.Dir = ".."
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("scaffold/client.go does not compile for js/wasm: %v\n%s", err, out)
	}
}
