//go:build wasm

package html_test

import (
	"testing"

	"webtyp.com/dom"
	. "webtyp.com/html"
)

type SearchChild struct {
	dom.Element
	InputEvents *dom.SignalString // misuse signal to count events for simplicity in test
	count int
}

func (c *SearchChild) Render() *dom.Element {
	return Input("text").
		ID(c.GetID()+"-input").
		OnInput(func(e dom.Event) {
			c.count++
			// c.InputEvents.Set(...)
		})
}

type ParentWithChild struct {
	dom.Element
	child *SearchChild
	toggle *dom.SignalBool
}

func (c *ParentWithChild) Init(ctx dom.Ctx) {
	c.child = &SearchChild{}
	c.toggle = dom.NewBool(true)
}

func (c *ParentWithChild) Render() *dom.Element {
	return Div().Child(
		dom.Show(c.toggle, c.child),
		Button().ID("toggle-btn").OnClick(func(e dom.Event) {
			c.toggle.Toggle()
		}),
	)
}

func TestChildListenersWithSignals(t *testing.T) {
	SetupDOM(t)
	parent := &ParentWithChild{}
	dom.Render("root", parent)

	id := parent.child.GetID() + "-input"
	TriggerEvent(id, "input", "a")
	if parent.child.count != 1 {
		t.Errorf("expected 1 input event, got %d", parent.child.count)
	}

	// Toggle off and on
	TriggerEvent("toggle-btn", "click", "")
	TriggerEvent("toggle-btn", "click", "")

	TriggerEvent(id, "input", "b")
	if parent.child.count != 2 {
		t.Errorf("expected 2 input events, got %d", parent.child.count)
	}
}
