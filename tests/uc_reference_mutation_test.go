//go:build wasm

package html_test

import (
	"strings"
	"syscall/js"
	"testing"

	"webtyp.com/dom"
	. "webtyp.com/html"
)

func TestReference_SetValue_PreservesListeners(t *testing.T) {
	SetupDOM(t)
	dom.Render("root", Input("text").ID("f1-input"))
	ref, _ := dom.Get("f1-input")
	fired := false
	ref.OnInput(func(e dom.Event) { fired = true })

	ref.SetValue("fixed")

	TriggerEvent("f1-input", "input", "fixed")

	if !fired {
		t.Error("FIX 1: SetValue destroyed the 'input' listener")
	}
	if ref.Value() != "fixed" {
		t.Errorf("Expected value 'fixed', got %q", ref.Value())
	}
}

func TestReference_SetAttr_PreservesListeners(t *testing.T) {
	SetupDOM(t)
	dom.Render("root", Button().ID("f2-btn").Text("Click"))
	ref, _ := dom.Get("f2-btn")
	clicked := false
	ref.OnClick(func(e dom.Event) { clicked = true })

	ref.SetAttr("disabled", "")
	ref.RemoveAttr("disabled")

	TriggerEvent("f2-btn", "click", "")

	if !clicked {
		t.Error("FIX 2: SetAttr/RemoveAttr destroyed the 'click' listener")
	}
}

func TestReference_SetText_PreservesAttributes(t *testing.T) {
	SetupDOM(t)
	dom.Render("root", Span().ID("f3-span").Attr("aria-live", "polite"))
	ref, _ := dom.Get("f3-span")

	ref.SetText("fixed text")

	rawEl := js.Global().Get("document").Call("getElementById", "f3-span")
	innerHTML := rawEl.Get("innerHTML").String()
	if strings.Contains(innerHTML, "<span") {
		t.Errorf("FIX 3: SetText nested a span. innerHTML=%q", innerHTML)
	}
	if innerHTML != "fixed text" {
		t.Errorf("Expected 'fixed text', got %q", innerHTML)
	}

	ariaLive := rawEl.Call("getAttribute", "aria-live").String()
	if ariaLive != "polite" {
		t.Errorf("FIX 3: aria-live lost after SetText, got %q", ariaLive)
	}
}

func TestReference_SetText_NoHTMLInjection(t *testing.T) {
	SetupDOM(t)
	dom.Render("root", Span().ID("f4-span"))
	ref, _ := dom.Get("f4-span")

	ref.SetText("<img src=x onerror=alert(1)>")

	rawEl := js.Global().Get("document").Call("getElementById", "f4-span")
	innerHTML := rawEl.Get("innerHTML").String()
	if strings.Contains(innerHTML, "<img") {
		t.Errorf("SECURITY: SetText interpreted HTML. innerHTML=%q", innerHTML)
	}
}
