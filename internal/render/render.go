package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
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

func RenderTemplate(w http.ResponseWriter, tmpl string, tmplD *models.TemplateData, r *http.Request) {
	var tmplCache map[string]*template.Template

	if appConfig.UseCache {
		tmplCache = appConfig.TemplateCache
	} else {
		tmplCache, _ = CreateTemplateCache()
	}

	t, ok := tmplCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache ", ok, t)
	}

	buf := new(bytes.Buffer)

	tmplD = AddDefaultData(tmplD, r)

	_ = t.Execute(buf, tmplD)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

}

//CreateTemplateCache creates template cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl.html", pathToTemplates))

	//println("pages: ", pages)

	if err != nil {
		//fmt.Println("template file path error")
		return myCache, err
	}
	//fmt.Println("alkaa for loop sivujen lapikaynti")
	for _, page := range pages {
		//fmt.Println("loopissa")
		name := filepath.Base(page)
		//fmt.Println("page filelistassa on", page, "ja name on ", name)

		tmplSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			//fmt.Println("Template setin luonti")
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
		//fmt.Println("tmplSet ", tmplSet)
	}

	//	for index, element := range myCache {
	//		fmt.Println(index, "=>", element)
	//	}
	//fmt.Println(" CreateTemplateCache() OK")
	return myCache, nil
}
