package render

import (
	"net/http"
	"testing"

	"github.com/t-Ikonen/bbbookingsystem/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "666")

	result := AddDefaultData(&td, r)

	if result.Flash != "666" {
		t.Error("flash value 666 not found in session")
	}

}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	appConfig.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter
	err = Template(&ww, "home.page.tmpl.html", &models.TemplateData{}, r)
	if err != nil {
		t.Error("Error wrtiting template to brovser", err)
	}

	err = Template(&ww, "nonexisting.tmpl.html", &models.TemplateData{}, r)
	if err == nil {
		t.Error("Should not be able to render non existing template", err)
	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(appConfig)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/jokuurli", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
