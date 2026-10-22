---
PLAN: "feat!: scaffold browser main at scaffold/browser/main.go, copied to web/browser/main.go"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 4338760399953196305
---

> Read AGENTS.md at the repo root before any change.

# Scaffold: `scaffold/browser/main.go`

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
```

## Why

The ecosystem moves every binary main into its own folder under `web/` (the folder is the
binary): `web/browser/main.go` → `web/public/browser.wasm`. `webtyp.com/app` copies this
module's scaffold into new projects; after this change it copies `scaffold/browser/main.go` to
`web/browser/main.go`. A main carries no destination tag (its folder separates it), so the
scaffold loses `//go:build browser`. The browser compiler still passes `-tags browser`.

## Stage 1 — move and retag the scaffold

1. Move `scaffold/browser.go` → `scaffold/browser/main.go` (same content otherwise).
2. Delete its first line `//go:build browser` and the blank line after it.
3. Header comment: `// web/browser.go — Browser WebAssembly entry point for this module.` →
   `// web/browser/main.go — Browser WebAssembly entry point for this module.`
4. `gotest` compiles `./...` natively. If `scaffold/browser` does not compile natively (it
   imports packages that only build for wasm), put `//go:build wasm` as its first line — the
   platform tag, never `browser` — and say so in the PR description.

`test ! -e scaffold/browser.go` must hold.

## Stage 2 — `tests/scaffold_test.go`

- The existence check: `filepath.Join("..", "scaffold", "browser", "main.go")`, message
  `scaffold/browser/main.go is consumed by webtyp.com/app and must exist: %v`.
- The build: `exec.Command("go", "build", "-tags", "browser", "-o", filepath.Join(t.TempDir(), "browser.wasm"), "./scaffold/browser")`
  (env `GOOS=js GOARCH=wasm` as today); failure message
  `scaffold/browser/main.go does not compile for js/wasm -tags browser: %v\n%s`.
- The doc comment of the test: "webtyp.com/app copies scaffold/browser/main.go into every new
  project as web/browser/main.go."

## Delivery check

```bash
gotest
test ! -e scaffold/browser.go
grep -rn 'scaffold/browser\.go\|web/browser\.go' --include=*.go --include=*.md . | grep -v docs/PLAN.md   # empty
```

A PR whose tests could not run is a failed delivery. Environment:
`GOPRIVATE=webtyp.com/*,github.com/veltylabs/*`, `GOFLAGS=-mod=mod`.

Never pause to ask permission or report progress: when the delivery check passes, commit, push
and open the PR. Stuck for 30 minutes → open the PR anyway listing what fails.

## Stages

| # | Files | Result |
|---|---|---|
| 1 | `scaffold/browser.go` → `scaffold/browser/main.go` | scaffold in its folder, no destination tag |
| 2 | `tests/scaffold_test.go` | test follows the new path |

## Executor notes

All steps were successfully completed. I moved the scaffold file to its own directory and removed the `//go:build browser` build tag as requested. It successfully compiled with `go build` inside its directory natively, so I did not need to append `//go:build wasm`. Tests and delivery checks passed successfully.
