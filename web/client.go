//go:build wasm

package main

import (
	"webtyp.com/dom"
	. "webtyp.com/html"
	"webtyp.com/fmt"
)

// --- App State & Components ---

type App struct {
	dom.Element
	currentRoute *dom.SignalString
	counter      *dom.SignalString
	countVal     int
}

func (a *App) Init(ctx dom.Ctx) {
	a.currentRoute = dom.NewString("#home")
	a.counter = dom.NewString("0")

	// 0. Restore theme from localStorage
	if theme, _ := dom.LocalStorageGet("theme"); theme != "" {
		dom.SetDocumentAttr("data-theme", theme)
	}

	// 1. Inject minimal CSS
	css := `
		body { font-family: sans-serif; margin: 0; background: #f4f4f9; color: #333; }
		nav { background: #333; padding: 1rem; display: flex; gap: 1rem; }
		nav a { color: white; text-decoration: none; cursor: pointer; padding: 0.2rem 0.5rem; border-radius: 4px; }
		nav a:hover { background: #555; }
		nav a.active { background: #0037ff; }
		.container { padding: 2rem; max-width: 800px; margin: 0 auto; }
		.card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
		.btn-group { display: flex; gap: 0.5rem; align-items: center; margin-top: 1rem; }
		button { padding: 0.5rem 1rem; cursor: pointer; border: none; border-radius: 4px; background: #007bff; color: white; }
		button:hover { background: #0056b3; }
		.count { font-size: 1.5rem; font-weight: bold; min-width: 3rem; text-align: center; }
	`
	renderStyle(css)

	// 2. Setup Routing
	dom.OnHashChange(func(hash string) {
		a.currentRoute.Set(hash)
	})

	// Initial route
	hash := dom.GetHash()
	if hash == "" {
		hash = "#home"
		dom.SetHash("#home")
	}
	a.currentRoute.Set(hash)
}

func (a *App) Render() *dom.Element {
	return Div().Child(
		// Navigation Bar
		Nav().Child(
			NavLink("Home", "#home", a.currentRoute),
			NavLink("About", "#about", a.currentRoute),
		),

		// Content Area
		Div().Child(
			dom.Show(dom.DeriveBool(func() bool { return a.currentRoute.Get() == "#about" }), a.renderAbout()),
			dom.Show(dom.DeriveBool(func() bool { return a.currentRoute.Get() == "#home" || a.currentRoute.Get() == "" }), a.renderHome()),
		).Class("container"),
	)
}

func (a *App) renderAbout() *dom.Element {
	return Div().Child(
		H1().Text("Sobre Esta Libreria."),
		P().Text("webtyp/dom is a minimalist, WASM-optimized DOM toolkit for Go."),
		P().Text("It features a JSX-like Builder API, Elm-inspired state management, and no Virtual DOM overhead."),
	).Class("card")
}

func (a *App) renderHome() *dom.Element {
	return Div().Child(
		H1().Text("Counter Example"),
		P().Text("This demonstrates local state updates and hash routing."),
		Div().Child(
			Button().Text("-").OnClick(func(e dom.Event) {
				a.countVal--
				a.counter.Set(fmt.Sprint(a.countVal))
			}),
			Span().BindText(a.counter).Class("count"),
			Button().Text("+").OnClick(func(e dom.Event) {
				a.countVal++
				a.counter.Set(fmt.Sprint(a.countVal))
			}),
		).Class("btn-group"),

		H2().Text("Persistence & Attributes"),
		P().Text("The theme is persisted in localStorage and applied to <html>."),
		Div().Child(
			Button().Text("Toggle Theme").OnClick(func(e dom.Event) {
				current := dom.GetDocumentAttr("data-theme")
				next := "dark"
				if current == "dark" {
					next = "light"
				}
				dom.SetDocumentAttr("data-theme", next)
				dom.LocalStorageSet("theme", next)
			}),
		).Class("btn-group"),
	).Class("card")
}

// --- Helpers ---

func NavLink(text, hash string, currentRoute *dom.SignalString) *dom.Element {
	return A(hash).
		Text(text).
		OnClick(func(e dom.Event) {
			e.PreventDefault()
			dom.SetHash(hash)
		}).
		BindClass("active", dom.DeriveBool(func() bool {
			return currentRoute.Get() == hash
		}))
}

func renderStyle(css string) {
	// Inject style into head
	dom.Append("head", Style().Text(css))
}

func main() {
	app := &App{}
	dom.Render("app", app)
	select {}
}
