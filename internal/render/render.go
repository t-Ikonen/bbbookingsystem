package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/t-Ikonen/bbbookingsystem/internal/config"
	"github.com/t-Ikonen/bbbookingsystem/internal/models"
)

var functions = template.FuncMap{}

var appConfig *config.AppConfig

var pathToTemplates = "./templates"

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = appConfig.Session.PopString(r.Context(), "flash")
	td.Error = appConfig.Session.PopString(r.Context(), "error")
	td.Warning = appConfig.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	return td
}

//NewTemplates sets the package for the template package
func NewTemplates(a *config.AppConfig) {
	appConfig = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string, tmplD *models.TemplateData, r *http.Request) error {
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
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl.html", pathToTemplates))
		if err != nil {
			fmt.Println("Matches")
			return myCache, err
		}
		if len(matches) > 0 {
			tmplSet, err = tmplSet.ParseGlob(fmt.Sprintf("%s/*.page.tmpl.html", pathToTemplates))
			if err != nil {

				return myCache, err
			}
		}
		myCache[name] = tmplSet

	}

	return myCache, nil
}
