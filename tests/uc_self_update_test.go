//go:build wasm

package html_test

import (
	"testing"

	"webtyp.com/dom"
	. "webtyp.com/html"
)

type SelfUpdater struct {
	dom.Element
	InputFired  int
	SelectFired int
	toggle *dom.SignalBool
}

func (c *SelfUpdater) Init(ctx dom.Ctx) {
	c.toggle = dom.NewBool(false)
}

func (c *SelfUpdater) Render() *dom.Element {
	return Div().Child(
		Input("text").ID(c.GetID()+"-search").OnInput(func(e dom.Event) {
			c.InputFired++
			c.toggle.Set(true)
		}),
		dom.Show(c.toggle, func() *dom.Element {
			return Div().ID(c.GetID() + "-options").OnClick(func(e dom.Event) {
				c.SelectFired++
			}).Text("Options")
		}),
	)
}

func TestSelfUpdateRewiresOnMountListeners(t *testing.T) {
	SetupDOM(t)
	c := &SelfUpdater{}
	dom.Render("root", c)

	TriggerEvent(c.GetID()+"-search", "input", "a")
	if c.InputFired != 1 {
		t.Errorf("InputFired=%d, want 1", c.InputFired)
	}

	TriggerEvent(c.GetID()+"-options", "click", "")
	if c.SelectFired != 1 {
		t.Errorf("SelectFired=%d, want 1", c.SelectFired)
	}
}
