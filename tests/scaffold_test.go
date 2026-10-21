//go:build !wasm

package html_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// webtyp.com/app copies scaffold/browser.go into every new project as
// web/browser.go. It must compile against this module's API with -tags browser,
// so an API change that breaks it fails here, in the same commit, not in a
// user's project.
func TestScaffoldBrowserCompiles(t *testing.T) {
	if _, err := os.Stat(filepath.Join("..", "scaffold", "browser.go")); err != nil {
		t.Fatalf("scaffold/browser.go is consumed by webtyp.com/app and must exist: %v", err)
	}

	cmd := exec.Command("go", "build", "-tags", "browser", "-o", filepath.Join(t.TempDir(), "browser.wasm"), "./scaffold")
	cmd.Dir = ".."
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("scaffold/browser.go does not compile for js/wasm -tags browser: %v\n%s", err, out)
	}
}
