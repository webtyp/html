package html_test

import (
	"strings"
	"testing"

	. "webtyp.com/html"
)

func TestDocument(t *testing.T) {
	opts := DocumentOptions{
		Lang:       "es",
		Title:      "Mi App",
		CSSURL:     "/styles.css",
		JSURL:      "/bundle.js",
		FaviconURL: "/favicon.ico",
		Head:       "<meta name=\"theme-color\" content=\"#000\">",
	}

	doc := Document(opts, Div().Text("Hola Mundo"))

	got := doc.String()

	checks := []string{
		"lang='es'",
		"<title>Mi App</title>",
		"rel='stylesheet'",
		"href='/styles.css'",
		"src='/bundle.js'",
		"defer",
		"rel='icon'",
		"href='/favicon.ico'",
		"<meta charset='utf-8'",
		"theme-color",
		"<div>Hola Mundo</div>",
	}

	for _, check := range checks {
		if !strings.Contains(got, check) {
			t.Errorf("expected to find %q in output, but didn't.\nGot: %s", check, got)
		}
	}
}

func TestDocumentString(t *testing.T) {
	got := DocumentString(DocumentOptions{Title: "App"})
	if !strings.HasPrefix(got, "<!DOCTYPE html>") {
		t.Errorf("expected <!DOCTYPE html> prefix, got: %s", got[:30])
	}
	if !strings.Contains(got, "<html") {
		t.Error("expected <html element")
	}
}

func TestDocumentViewportDefault(t *testing.T) {
	got := DocumentString(DocumentOptions{})
	if !strings.Contains(got, `name='viewport'`) && !strings.Contains(got, `name="viewport"`) {
		t.Errorf("empty DocumentOptions must emit viewport meta, got: %s", got)
	}
	if !strings.Contains(got, "viewport-fit=cover") {
		t.Errorf("empty DocumentOptions must emit viewport-fit=cover, got: %s", got)
	}
}

func TestDocumentViewportOrder(t *testing.T) {
	got := DocumentString(DocumentOptions{
		CSSURL: "/styles.css",
		JSURL:  "/bundle.js",
	})
	charsetIdx := indexOfMeta(got, "charset")
	viewportIdx := indexOfMeta(got, "viewport")
	linkIdx := strings.Index(got, "<link")
	scriptIdx := strings.Index(got, "<script")
	if charsetIdx < 0 {
		t.Fatal("charset meta missing")
	}
	if viewportIdx < 0 {
		t.Fatal("viewport meta missing")
	}
	if charsetIdx > viewportIdx {
		t.Errorf("charset must appear before viewport: charset=%d viewport=%d", charsetIdx, viewportIdx)
	}
	if linkIdx >= 0 && viewportIdx > linkIdx {
		t.Errorf("viewport must appear before link: viewport=%d link=%d", viewportIdx, linkIdx)
	}
	if scriptIdx >= 0 && viewportIdx > scriptIdx {
		t.Errorf("viewport must appear before script: viewport=%d script=%d", viewportIdx, scriptIdx)
	}
}

func TestDocumentViewportOnce(t *testing.T) {
	got := DocumentString(DocumentOptions{})
	count := strings.Count(got, "viewport")
	// name=viewport + content containing viewport-fit → at least 2 substrings of "viewport"
	// Guard: exactly one meta viewport tag (one name='viewport' or name="viewport")
	nameCount := strings.Count(got, `name='viewport'`) + strings.Count(got, `name="viewport"`)
	if nameCount != 1 {
		t.Errorf("expected exactly one viewport meta, got %d (viewport substrings=%d)\n%s", nameCount, count, got)
	}
}

func TestDocumentViewportNoZoomLock(t *testing.T) {
	got := DocumentString(DocumentOptions{})
	if strings.Contains(got, "user-scalable") {
		t.Errorf("must not contain user-scalable, got: %s", got)
	}
	if strings.Contains(got, "maximum-scale") {
		t.Errorf("must not contain maximum-scale, got: %s", got)
	}
}

func indexOfMeta(s, needle string) int {
	// Prefer name= / charset= attribute occurrence inside a meta tag region.
	idx := strings.Index(s, needle)
	return idx
}

// TestDocumentSEOMetadata locks in the per-page metadata a multi-page static
// site needs to rank each URL on its own. Open Graph tags must use
// property=, not name= — a crawler reading OG ignores the name= form, so
// getting this wrong yields tags that look right in the source and do
// nothing.
func TestDocumentSEOMetadata(t *testing.T) {
	opts := DocumentOptions{
		Title:       "Oftalmología en Chillán",
		Description: "Atención oftalmológica accesible en el centro de Chillán.",
		Canonical:   "https://example.com/especialidades/oftalmologia/",
		Image:       "https://example.com/img/consultorio.webp",
		JSONLD:      `{"@context":"https://schema.org","@type":"MedicalClinic"}`,
	}

	got := DocumentString(opts, Div().Text("contenido"))

	for _, want := range []string{
		`<meta name='description' content='Atención oftalmológica accesible en el centro de Chillán.'>`,
		`<meta property='og:description' content='Atención oftalmológica accesible en el centro de Chillán.'>`,
		`<meta property='og:title' content='Oftalmología en Chillán'>`,
		`<link rel='canonical' href='https://example.com/especialidades/oftalmologia/'>`,
		`<meta property='og:url' content='https://example.com/especialidades/oftalmologia/'>`,
		`<meta property='og:image' content='https://example.com/img/consultorio.webp'>`,
		`<script type='application/ld+json'>{"@context":"https://schema.org","@type":"MedicalClinic"}</script>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s\ngot:\n%s", want, got)
		}
	}
}

// TestDocumentOmitsUnsetSEOFields guards the zero value: a document that
// declares no SEO metadata must not emit empty meta tags, which are worse
// than absent ones (an empty description competes with the crawler's own
// snippet extraction).
func TestDocumentOmitsUnsetSEOFields(t *testing.T) {
	got := DocumentString(DocumentOptions{Title: "Solo título"}, Div())

	for _, unwanted := range []string{
		`name='description'`, `rel='canonical'`, `og:url`, `og:image`, `application/ld+json`,
	} {
		if strings.Contains(got, unwanted) {
			t.Errorf("emitted %s for an unset field\ngot:\n%s", unwanted, got)
		}
	}
	if !strings.Contains(got, `<meta property='og:title' content='Solo título'>`) {
		t.Errorf("og:title should follow a set Title, got:\n%s", got)
	}
}
