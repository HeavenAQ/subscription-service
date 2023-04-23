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
		fmt.Sprintf("%s/base.layout.html", templatePath),
		fmt.Sprintf("%s/header.partial.html", templatePath),
		fmt.Sprintf("%s/navbar.partial.html", templatePath),
		fmt.Sprintf("%s/footer.partial.html", templatePath),
		fmt.Sprintf("%s/alerts.partial.html", templatePath),
	}

	if td == nil {
		td = &TemplateData{}
	}

	templ, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.failedToRender(w, err)
		return
	}

	if err = templ.Execute(w, app.AddDefaultData); err != nil {
		app.failedToRender(w, err)
		return
	}
}

func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	if app.isAuthenticated(r) {
		td.Authenticated = true
	}
	td.Now = time.Now()
	return td
}

func (app *Config) isAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
