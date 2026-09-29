# Agent Guide — `webtyp/html`

Constraints for agents working on the HTML tag builder. Read this before any change.

---

## Construction Harness — typed & explicit (the WebTyp approach)

This library is the **tag builder** layer of WebTyp's construction harness: the typed, explicit API
is what keeps an agent that doesn't know the library from building wrong code. Model it on
`webtyp/json` (typed methods per primitive, zero `any` in the data path).

- **No `any` in the builder.** Tag constructors are **no-arg** (`Span()`, `Div()`, `Button()`); compose
  via typed `*dom.Element` methods (`Text`, `Child`, `Attr`, `Class`, `Set`). Never a `Tag(...any)`
  slot. Constructors may take **semantic** args only (`Input(type)`, `Option(value,text)`).
- **Reuse `fmt` types — no duplicates.** Attributes are `fmt.KeyValue` (declared in `fmt/parse.go`);
  children are `*dom.Element`. Do **not** declare a `Node`/`Text`/attr wrapper type.
- **Reactive content is typed** — anything that changes goes only through a signal binding
  (`BindText`/`Bind*`), which requires a `*dom.Signal*`. The builder composes only **static** content.
- **Explicit names** — `.Text` (static) vs `.BindText` (reactive); reading a call states intent.
- **Docs are minimal "how" instructions, not long skills** — if a rule must be *remembered*, close it
  with types, not prose.

(Ecosystem rationale: skill **api-design** — the construction harness, in
`webtyp/devskills/skills/api-design/SKILL.md`.)

---

## WASM / TinyGo

- No Go stdlib: use `webtyp.com/fmt`. DOM types via `webtyp.com/dom` (dot-imported),
  never `syscall/js`. `switch` not `map`; no `defer/recover`; embed `dom.Element` by value.

## `scaffold/client.go` — consumed by `webtyp.com/app`

`webtyp.com/app` copies this file, from the `html` version it installs, as the default `web/client.go`
of every new project. Do not move or rename it. When an API change breaks it,
`tests/scaffold_test.go` fails; update the scaffold in the same commit.

## Testing & Docs

```bash
go install webtyp.com/devflow/cmd/gotest@latest
gotest
```

`gotest`, never `go test`; stdlib assertions only; dual WASM/stdlib. Publish with `gopush 'message'`.
Documentation first: update `docs/` and re-index `README.md` before code. Diagrams: `flowchart TD`,
no `subgraph`, `<br/>` for breaks.
