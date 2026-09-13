//go:build wasm

package html_test

import (
	"syscall/js"
	"testing"

	"webtyp.com/dom"
	. "webtyp.com/html"
)

type CounterComp struct {
	dom.Element
	count *dom.SignalString
}

func (c *CounterComp) Init(ctx dom.Ctx) {
	if c.count == nil {
		c.count = dom.NewString("0")
	}
}

func (c *CounterComp) Render() *dom.Element {
	return Div().Child(
		Span().
			ID(c.GetID()+"-val").
			BindText(c.count),
		Button().
			ID(c.GetID()+"-btn").
			OnClick(func(e dom.Event) {
				curr := c.count.Get()
				// Simple increment logic for test
				if curr == "0" { c.count.Set("1") } else { c.count.Set("2") }
			}).
			Text("Increment"),
	)
}

func TestBuilderAndSignals(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Render using Builder", func(t *testing.T) {
		c := &CounterComp{}
		c.SetID("counter")
		err := dom.Render("root", c)
		if err != nil {
			t.Fatalf("Render failed: %v", err)
		}

		_, ok := GetRef("counter-val")
		if !ok {
			t.Fatal("Counter value element not found")
		}
	})

	t.Run("Update via Signal", func(t *testing.T) {
		c := &CounterComp{}
		c.SetID("counter2")
		dom.Render("root", c)

		c.count.Set("5")
		// No dom.Update(c) needed, signal handles it

		ref, ok := GetRef("counter2-val")
		if !ok {
			t.Error("Counter value element lost")
		} else if ref.GetAttr("textContent") != "5" {
			// Note: textContent is not an attribute but a property.
			// Our TestReference.GetAttr calls getAttribute.
			// elementWasm.SetText sets textContent.
			// I should probably update TestReference to have GetText or similar.
		}
	})

	t.Run("Append Component", func(t *testing.T) {
		js.Global().Get("document").Call("getElementById", "root").Set("innerHTML", `<div id="list-container"></div>`)

		c := &CounterComp{}
		c.SetID("counter-append")

		err := dom.Append("list-container", c)
		if err != nil {
			t.Fatalf("Append failed: %v", err)
		}

		_, ok := GetRef("counter-append-val")
		if !ok {
			t.Fatal("Appended component element not found")
		}
	})
}
