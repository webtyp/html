//go:build wasm

package html_test

import (
	"testing"

	"webtyp.com/dom"
	. "webtyp.com/html"
)

type trackableComp struct {
	dom.Element
	clickCount int
}

func (c *trackableComp) Render() *dom.Element {
	return Div().ID(c.GetID()).OnClick(func(e dom.Event) {
		c.clickCount++
	})
}

// TestUpdateLifecycle_NoDuplicateListeners verifies that re-rendering the same
// component via dom.Render does not accumulate duplicate event listeners.
// Without this, a click would fire the handler once per render
// (e.g. a form submitting multiple POST requests on a single user action).
func TestUpdateLifecycle_NoDuplicateListeners(t *testing.T) {
	doc := SetupDOM(t)

	comp := &trackableComp{}
	comp.SetID("trackable")

	if err := dom.Render("root", comp); err != nil {
		t.Fatalf("first Render failed: %v", err)
	}

	// Re-render (update) the same component.
	if err := dom.Render("root", comp); err != nil {
		t.Fatalf("second Render failed: %v", err)
	}

	rawEl := doc.Call("getElementById", "trackable")
	TriggerEvent("trackable", "click", "")
	_ = rawEl

	if comp.clickCount != 1 {
		t.Errorf("expected clickCount=1 after a single click post-update, got %d (bug: duplicate listeners)", comp.clickCount)
	}
}
