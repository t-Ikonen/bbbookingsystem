package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/t-Ikonen/bbbookingsystem/internal/config"
	"github.com/t-Ikonen/bbbookingsystem/internal/models"
)

var functions = template.FuncMap{
	"shortDate":  ShortDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

var appConfig *config.AppConfig

var pathToTemplates = "./templates"

//fomat time
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

//iterate return a slice of ints , starting at 1, going to count
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

//adds 2 ints, returns sum
func Add(a, b int) int {
	return a + b
}

//NewRender sets the package for the template package
func NewRenderer(a *config.AppConfig) {
	appConfig = a
}

//ShortDate formats long date format to just date without time aka HumanDate
func ShortDate(t time.Time) string {
	return t.Format("02-01-2006")
}

//AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = appConfig.Session.PopString(r.Context(), "flash")
	td.Error = appConfig.Session.PopString(r.Context(), "error")
	td.Warning = appConfig.Session.PopString(r.Context(), "warning")
	if appConfig.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}

	td.CSRFToken = nosurf.Token(r)
	return td
}

//Template renders templates using html.tmpl to full htmp pages
func Template(w http.ResponseWriter, tmpl string, tmplD *models.TemplateData, r *http.Request) error {
	var tmplCache map[string]*template.Template

	if appConfig.UseCache {
		tmplCache = appConfig.TemplateCache
	} else {
		tmplCache, _ = CreateTemplateCache()
	}

	t, ok := tmplCache[tmpl]
	if !ok {
		//log.Fatal("Could not get template from template cache ", ok, t)
		return errors.New("cannot get templates from cache")
	}
	buf := new(bytes.Buffer)
	tmplD = AddDefaultData(tmplD, r)
	_ = t.Execute(buf, tmplD)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}
	return nil
}

//CreateTemplateCache creates template cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl.html", pathToTemplates))
	if err != nil {

		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmplSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl.html", pathToTemplates))
		if err != nil {
			fmt.Println("Matches error")
			return myCache, err
		}
		if len(matches) > 0 {
			tmplSet, err = tmplSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl.html", pathToTemplates))
			if err != nil {

				return myCache, err
			}
		}
		myCache[name] = tmplSet

	}

	return myCache, nil
}
