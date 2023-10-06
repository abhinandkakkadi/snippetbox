package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/abhinandkakkadi/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// returns nicely formatted time as string (should only return one value, but can return error as second value)
func humanDate(t time.Time) string {

	// return empty string if time has zero value
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04") // custom template function
}

// this is a string-keyed map which acts as a lookup between the names of our custom
// template functions and the functions themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	// output is a slice of string and will be of form ./ui/html/pages/create/tmpl ect.
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		// Extract the file name (like 'home.tmpl') from full filepath
		name := filepath.Base(page)

		// template.FuncMap must be registered with the templates set before we call ParseFiles()
		// call template.New() to create an empty template set
		// use Funcs() method to register the template.FuncMap, and then parse the file as normal
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add any partials.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set as value with key being name of page (like 'home.tmpl')
		cache[name] = ts
	}

	return cache, nil

}
