//go:build wasm

package html_test

import (
	"testing"

	"webtyp.com/dom"
	. "webtyp.com/html"
)

// R1: BindText patterns capture state updated after Signal.Set()
type StateCapturer struct {
	dom.Element
	Value *dom.SignalString
	LastCaptured string
}

func (c *StateCapturer) Init(ctx dom.Ctx) {
	c.Value = dom.NewString("v1")
}

func (c *StateCapturer) Render() *dom.Element {
	return Div().Child(
		Button().ID("btn-r1").OnClick(func(e dom.Event) {
			c.LastCaptured = c.Value.Get()
		}),
	)
}

func TestR1StateCapture(t *testing.T) {
	SetupDOM(t)
	c := &StateCapturer{}
	dom.Render("root", c)

	c.Value.Set("v2")

	TriggerEvent("btn-r1", "click", "")
	if c.LastCaptured != "v2" {
		t.Errorf("R1: expected captured value 'v2', got %q", c.LastCaptured)
	}
}

// R3: Race conditions / Re-wiring during Signal updates
type RewireComp struct {
	dom.Element
	Fired int
	Trigger *dom.SignalBool
}

func (c *RewireComp) Init(ctx dom.Ctx) {
	c.Trigger = dom.NewBool(false)
}

func (c *RewireComp) Render() *dom.Element {
	return Div().Child(
		Button().ID("btn-r3").OnClick(func(e dom.Event) {
			c.Fired++
			c.Trigger.Toggle()
		}),
	)
}

func TestR3Rewiring(t *testing.T) {
	SetupDOM(t)
	c := &RewireComp{}
	dom.Render("root", c)

	for i := 1; i <= 5; i++ {
		TriggerEvent("btn-r3", "click", "")
		if c.Fired != i {
			t.Errorf("R3 iteration %d: expected Fired=%d, got %d", i, i, c.Fired)
		}
	}
}
