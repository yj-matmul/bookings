package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/yj-matmul/bookings/internal/config"
	"github.com/yj-matmul/bookings/internal/models"
)

// implement custom function when go lang have no built-in function I want to
var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}
var pathToTemplates = "./templates"

var app *config.AppConfig

// Add
func Add(a, b int) int {
	return a + b
}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate returns time in YYYY-MM-DD
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// Iterate returns a slice of ints, starting at 1, going to count
func Iterate(count int) []int {
	var items []int

	for i := 1; i <= count; i++ {
		items = append(items, i)
	}

	return items
}

// FormatDate returns time in layout
func FormatDate(t time.Time, layout string) string {
	return t.Format(layout)
}

// AddDefaultData sets default template data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// Template renders template using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	// get the template cache from the app config
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from the template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		log.Println("CreateTemplateCache, glob page.html", err)
		return templateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			log.Println("CreateTemplateCache, parsing page.html", err)
			return templateCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			log.Println("CreateTemplateCache, glob layout.html", err)
			return templateCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				log.Println("CreateTemplateCache, parsing layout.html", err)
				return templateCache, err
			}
		}

		templateCache[name] = ts
	}

	return templateCache, nil
}
