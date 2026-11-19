package html

import (
	. "webtyp.com/dom"
)

// DocumentOptions configures the document shell.
type DocumentOptions struct {
	Lang         string // default "en"
	Title        string
	CSSURL       string // <link rel="stylesheet">
	CSSIntegrity string // optional SRI hash (e.g. "sha384-...") for CSSURL
	JSURL        string // <script src defer>
	JSIntegrity  string // optional SRI hash (e.g. "sha384-...") for JSURL
	FaviconURL   string // <link rel="icon">
	Head         string // extra markup for <head> (optional)

	// Description, Canonical and Image are the per-page metadata a search or
	// answer engine reads. They are typed fields rather than markup pushed
	// through Head so that a site emitting many pages cannot silently ship
	// two of them sharing one description — the single most common way a
	// multi-page static site loses its per-page ranking.
	Description string // <meta name="description"> and og:description
	Canonical   string // absolute URL; <link rel="canonical"> and og:url
	Image       string // absolute URL; og:image

	// JSONLD is a schema.org document embedded as
	// <script type="application/ld+json">. It stays a string because its
	// shape is caller-defined (LocalBusiness, MedicalClinic, Service, …) and
	// modelling every schema.org type is not this package's job. Emitted
	// verbatim via dom.Raw/dom.Trust — HTML-escaping it would corrupt the
	// JSON (quotes become &quot;) — so the caller is responsible for it being
	// valid JSON built from their own data, never HTML content from a
	// request/database/third party passed straight through.
	JSONLD string
}

func metaNamed(name, content string) *Element {
	return NewElement("meta").NoCloseTag().Attr("name", name).Attr("content", content)
}

// metaProperty emits the property= form Open Graph requires; a name= meta is
// ignored by crawlers reading OG tags.
func metaProperty(property, content string) *Element {
	return NewElement("meta").NoCloseTag().Attr("property", property).Attr("content", content)
}

// Document builds the complete document shell as *Element (the <html> element).
// Does not include <!DOCTYPE html> — use DocumentString for the full output.
func Document(opts DocumentOptions, body ...Component) *Element {
	if opts.Lang == "" {
		opts.Lang = "en"
	}

	head := NewElement("head")
	head.Child(NewElement("meta").NoCloseTag().Attr("charset", "utf-8"))
	head.Child(NewElement("meta").NoCloseTag().
		Attr("name", "viewport").
		Attr("content", "width=device-width, initial-scale=1, viewport-fit=cover"))

	if opts.Title != "" {
		head.Child(NewElement("title").Text(opts.Title))
		head.Child(metaProperty("og:title", opts.Title))
	}
	if opts.Description != "" {
		head.Child(metaNamed("description", opts.Description))
		head.Child(metaProperty("og:description", opts.Description))
	}
	if opts.Canonical != "" {
		head.Child(NewElement("link").NoCloseTag().Attr("rel", "canonical").Attr("href", opts.Canonical))
		head.Child(metaProperty("og:url", opts.Canonical))
	}
	if opts.Image != "" {
		head.Child(metaProperty("og:image", opts.Image))
	}
	if opts.JSONLD != "" {
		head.Child(NewElement("script").Attr("type", "application/ld+json").Raw(Trust(opts.JSONLD)))
	}
	if opts.FaviconURL != "" {
		head.Child(NewElement("link").NoCloseTag().Attr("rel", "icon").Attr("href", opts.FaviconURL))
	}
	if opts.CSSURL != "" {
		link := NewElement("link").NoCloseTag().Attr("rel", "stylesheet").Attr("href", opts.CSSURL)
		if opts.CSSIntegrity != "" {
			link.Attr("integrity", opts.CSSIntegrity).Attr("crossorigin", "anonymous")
		}
		head.Child(link)
	}
	if opts.JSURL != "" {
		script := NewElement("script").Attr("src", opts.JSURL).Attr("defer", "")
		if opts.JSIntegrity != "" {
			script.Attr("integrity", opts.JSIntegrity).Attr("crossorigin", "anonymous")
		}
		head.Child(script)
	}
	if opts.Head != "" {
		head.Text(opts.Head)
	}

	bodyEl := NewElement("body").Child(body...)

	return NewElement("html").Attr("lang", opts.Lang).Child(head, bodyEl)
}

// DocumentString renders the full HTML document including <!DOCTYPE html>.
// This is what assetmin writes to disk as index.html.
func DocumentString(opts DocumentOptions, body ...Component) string {
	return "<!DOCTYPE html>" + Document(opts, body...).String()
}
