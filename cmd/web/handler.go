package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", &TemplateData{})
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", &TemplateData{})
}

func (app *Config) FailToLogin(w http.ResponseWriter, r *http.Request, err error) {
	app.Session.Put(r.Context(), "error", "Invalid credentials.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	app.ErrorLog.Println(err)
}

func (app *Config) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())

	// parse post form data
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// get email and password from form post
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.FailToLogin(w, r, err)
		return
	}

	// Check password
	validPassword, err := user.PasswordMatches(password)
	if !validPassword || err != nil {
		app.FailToLogin(w, r, err)
		return
	}

	// log user in
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	app.Session.Put(r.Context(), "flash", "Successfully logged in!")
	// redirect the user
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", &TemplateData{})
}

func (app *Config) PostRegister(w http.ResponseWriter, r *http.Request) {
	// create a user

	// send an activation email

	// subscribe to the newsletter
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url

	// generate an invoice

	// send an email with attachments

	// send an email with the invoice attached
}
