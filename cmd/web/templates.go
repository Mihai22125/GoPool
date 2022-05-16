package main

import (
	"time"
	"path/filepath"
	"html/template"
	"github.com/mihai22125/goPool/pkg/models"
	"github.com/mihai22125/goPool/pkg/forms"

)

type templateData struct {
	AuthenticatedUser int
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	Pool              *models.Pool
	Pools             []*models.Pool
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 15:04:05")
}

func intRange(start, end int) []int {
	n := end - start
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = start + i
	}
	return result
}

var functions = template.FuncMap {
	"humanDate": humanDate,
	"intRange": intRange,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}