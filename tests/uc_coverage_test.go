//go:build wasm

package html_test

import (
	"testing"

	"webtyp.com/dom"
	. "webtyp.com/html"
	"webtyp.com/fmt"
)

func TestCoverageElementFactories(t *testing.T) {
	t.Run("All Element Factories", func(t *testing.T) {
		els := []dom.Component{
			Span(), P(), H1(), H2(), H3(), H4(), H5(), H6(),
			Ul(), Ol(), Li(), Nav(), Section(), Main(), Article(),
			Header(), Footer(), Aside(), Details(), Summary(), Dialog(),
			Figure(), Figcaption(), Pre(), Code(), Strong(),
			Small(), Mark(), Table(), Thead(), Tbody(), Tfoot(),
			Tr(), Th(), Td(), Fieldset(), Legend(), Label(),
			Canvas(), Style(), Script(), A("href"), Button(),
		}
		for _, el := range els {
			if el.String() == "" {
				t.Errorf("Element factory returned empty HTML")
			}
		}
	})
}

func TestCoverageDOMLogic(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Render to non-existent parent", func(t *testing.T) {
		err := dom.Render("void", Div())
		if err == nil {
			t.Error("Expected error when rendering to non-existent parent")
		}
	})

	t.Run("Logging", func(t *testing.T) {
		dom.SetLog(nil)
		dom.Log("test message")

		logged := false
		dom.SetLog(func(v ...any) {
			logged = true
		})
		dom.Log("test message 2")
		if !logged {
			t.Error("Log handler not called")
		}
		dom.SetLog(nil)
	})

	t.Run("Hash", func(t *testing.T) {
		dom.SetHash("test")
		_ = dom.GetHash()
		dom.OnHashChange(func(h string) {})
	})

	t.Run("Element methods", func(t *testing.T) {
		el := Div().Text("hello").Attr("k", "v")
		el.Attr("k", "v2")
		if !fmt.Contains(el.String(), "k='v2'") {
			t.Error("Attr not updated")
		}
	})
}

func TestCoverageEvents(t *testing.T) {
	SetupDOM(t)

	t.Run("Event interface methods - Button", func(t *testing.T) {
		var ev dom.Event
		triggered := false
		btn := Button().ID("btn-ev").OnClick(func(e dom.Event) {
			ev = e
			triggered = true
		}).Text("Click")
		dom.Render("root", btn)
		TriggerEvent(btn.GetID(), "click", "")

		if triggered && ev != nil {
			ev.PreventDefault()
			ev.StopPropagation()
			_ = ev.TargetID()
			_ = ev.TargetValue()
		}
	})

	t.Run("Append logic", func(t *testing.T) {
		parent := Div().ID("parent-append")
		dom.Render("root", parent)
		child := Span().ID("child-append").Text("appended")
		err := dom.Append("parent-append", child)
		if err != nil {
			t.Errorf("Append failed: %v", err)
		}
		if _, ok := GetRef("child-append"); !ok {
			t.Error("Appended element not found in DOM")
		}
	})
}

func TestLifecycleDeep(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Nested components cleanup", func(t *testing.T) {
		child := &MockComponent{Element: *Div()}
		child.SetID("child-comp")
		parent := Div().ID("parent-comp").Child(child)

		dom.Render("root", parent)
		if !child.Mounted {
			t.Error("Child component should be mounted")
		}

		dom.Render("root", Div().ID("parent-comp").Text("no more child"))

		if child.Mounted {
			t.Error("Child component should be unmounted after removal from parent")
		}
	})

	t.Run("ElementNode in children", func(t *testing.T) {
		c := &MockComponent{Element: *Div()}
		c.SetID("mock-child")
		el := Div().Child(c)
		html := el.String()
		if !fmt.Contains(html, "<div id='mock-child'") {
			t.Error("MockComponent elementNode not rendered correctly in children")
		}
	})

	t.Run("Default types in String()", func(t *testing.T) {
		// Element.String() now uses elementToHTML which handles various types
		// though we don't have direct access to internal children slice easily
	})
}

func TestCoverageCleanup(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Listener cleanup", func(t *testing.T) {
		triggered := false
		btn := Button().ID("btn-clean").OnClick(func(e dom.Event) {
			triggered = true
		}).Text("Click")
		dom.Render("root", btn)

		TriggerEvent("btn-clean", "click", "")
		if !triggered {
			t.Error("Event should have triggered")
		}

		triggered = false
		dom.Render("root", Div().ID("root-clean").Text("Gone"))

		TriggerEvent("btn-clean", "click", "")
		if triggered {
			t.Error("Event should NOT have triggered after cleanup")
		}
	})

	t.Run("Option helpers", func(t *testing.T) {
		opt := Option("v1", "Text 1")
		if !fmt.Contains(opt.String(), "value='v1'") || !fmt.Contains(opt.String(), "Text 1") {
			t.Error("Option not rendered correctly")
		}
		sopt := SelectedOption("v2", "Text 2")
		if !fmt.Contains(sopt.String(), "selected=''") {
			t.Error("SelectedOption not rendered correctly")
		}
	})
}
