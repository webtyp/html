---
PLAN: "refactor: rename scaffold/client.go to scaffold/browser.go with //go:build browser"
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 1466854644911309856
PR: https://github.com/webtyp/html/pull/8
REVIEWED_COMMIT: de44e549824b02d2c9abba3d2d7b46c37398a59c
---

> Read `AGENTS.md` at the repo root before any change.

# `webtyp/html` — `scaffold/browser.go` (`//go:build browser`)

Part of [BUILD_DESTINATIONS_MASTER_PLAN.md](https://github.com/webtyp/docs/blob/main/BUILD_DESTINATIONS_MASTER_PLAN.md) (D1, D3, D4, D15): the browser entry point in every WebTyp project is `web/browser.go` carrying `//go:build browser` and compiled with `-tags browser` into `browser.wasm`. `webtyp.com/app` copies `scaffold/browser.go` from `webtyp.com/html/scaffold` when initializing a project.

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
```

## Design gate

1. **Prior art:** `package.json` `browser` condition and WebTyp's `BUILD_DESTINATIONS` D1/D15 (`web/browser.go` → `//go:build browser` → `browser.wasm`).
2. **Novice-name test:** `scaffold/browser.go` matches the destination file `web/browser.go` and build tag `//go:build browser` 1:1.
3. **Complexity ledger:** +0 files / −1 old name (`scaffold/client.go`).
4. **Where it belongs:** `webtyp/html/scaffold` owns the canonical browser entry template consumed by `webtyp/app`.
5. **What it deletes:** `scaffold/client.go` (renamed via `git mv` to `scaffold/browser.go`, no backwards-compatibility copy).

---

## Stage 1 — Rename `scaffold/client.go` to `scaffold/browser.go` and update build tag

File: `scaffold/browser.go` (renamed from `scaffold/client.go`)

1. Run `git mv scaffold/client.go scaffold/browser.go`.
2. Change line 1 from `//go:build wasm` to `//go:build browser`.
3. Change line 3 from `// web/client.go — WebAssembly entry point for this module.` to `// web/browser.go — Browser WebAssembly entry point for this module.`.

---

## Stage 2 — Update `tests/scaffold_test.go`

File: `tests/scaffold_test.go`

Replace the contents of `tests/scaffold_test.go` with:

```go
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
```

Run `gotest ./...` and verify all tests pass.

---

## Deletions and acceptance criteria

- `scaffold/client.go` does not exist (`test ! -e scaffold/client.go`).
- `grep -rn --exclude="LAST_PLAN_EXECUTED.md" "client\.go\|client\.wasm" .` returns empty.
- `gotest ./...` passes with zero failures.

## Stages

| # | Files | Goal |
|---|---|---|
| 1 | `scaffold/browser.go` (deleted: `scaffold/client.go`) | Rename scaffold to `scaffold/browser.go` with `//go:build browser` |
| 2 | `tests/scaffold_test.go` | Compile `./scaffold` with `GOOS=js GOARCH=wasm -tags browser` to `browser.wasm` |


## Executor notes

- All steps completed successfully.
- Renamed `scaffold/client.go` to `scaffold/browser.go`.
- Updated build tags and comments.
- Updated `tests/scaffold_test.go`.
- Verified all tests pass with `gotest`.
