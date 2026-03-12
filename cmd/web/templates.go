package main

import (
	"html/template"
	"path/filepath"
	"time"

	"stash.io/internal/models"
)

// TemplateData holds dynamic data passed to HTML templates.
type TemplateData struct {
	CurrentYear int
	Form        any
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
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
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
