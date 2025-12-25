package render

import (
	"net/http"
	"testing"

	"booking-app/internal/models"
)

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc
	r, _ := getSession()

	var ww myWriter
	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error(err)
	}

	err = Template(&ww, r, "non-existence.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that doesnot exist")
	}
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "this is nice message")

	result := AddDefaultData(&td, r)

	if result == nil {
		t.Error("failed to add default data")
		return
	}

	if result.Flash != "this is nice message" {
		t.Error("flash value not found in session")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func TestNewRenderTemplate(t *testing.T) {
	NewRenderer(&testApp)
}
