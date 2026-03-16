package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"stash.io/internal/models"
	"stash.io/ui"
)

// TemplateData holds dynamic data passed to HTML templates.
type TemplateData struct {
	CurrentYear     int
	Form            any
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// NewTemplateCache builds a map of parsed template sets keyed by page filename.
// Each entry includes the base layout, nav partial, and the page template.
func NewTemplateCache() (map[string]*template.Template, error) {
	// Map to act as cache
	cache := map[string]*template.Template{}

	// Get a slice of all filepaths matching the page template pattern
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*tmpl.html",
			page,
		}

		// Parse the base template file into a template set.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
