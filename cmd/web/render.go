package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

const templatePath = "./cmd/web/templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	AnyData       map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
}

func (app *Config) failedToRender(w http.ResponseWriter, err error) {
	app.ErrorLog.Print(err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, templateName string, td *TemplateData) {
	templateSlice := []string{
		fmt.Sprintf("%s/%s", templatePath, templateName),
		fmt.Sprintf("%s/base.layout.gohtml", templatePath),
		fmt.Sprintf("%s/header.partial.gohtml", templatePath),
		fmt.Sprintf("%s/navbar.partial.gohtml", templatePath),
		fmt.Sprintf("%s/footer.partial.gohtml", templatePath),
		fmt.Sprintf("%s/alerts.partial.gohtml", templatePath),
	}

	if td == nil {
		td = &TemplateData{}
	}

	templ, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.failedToRender(w, err)
		return
	}

	if err = templ.Execute(w, app.AddDefaultData(td, r)); err != nil {
		app.failedToRender(w, err)
		return
	}
	app.InfoLog.Println(http.StatusOK, r.URL)
}

func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	if app.IsAuthenticated(r) {
		td.Authenticated = true
	}
	td.Now = time.Now()
	return td
}

func (app *Config) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
