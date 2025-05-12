package handler

import (
	"html/template"
	"net/http"
	"path/filepath"

	"kvstore/internal/auth"
	"kvstore/internal/config"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router) {
	r.Get("/login", loginPage)
	r.Post("/login", loginHandler)
	r.Get("/register", registerPage)
	r.Post("/register", registerHandler)
	r.Get("/logout", logoutHandler)
}

// Login page
func loginPage(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(filepath.Join(config.Cfg.TemplatesDir, "login.html"))
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.Execute(w, nil)
}

// Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if !auth.ValidateUser(username, password) {
		tpl, err := template.ParseFiles(filepath.Join(config.Cfg.TemplatesDir, "login.html"))
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tpl.Execute(w, map[string]string{"Error": "Invalid username or password"})
		return
	}

	token := auth.CreateSession(username)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

// Register page
func registerPage(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles(filepath.Join(config.Cfg.TemplatesDir, "register.html")))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.Execute(w, nil)
}

// Register handler
func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	if password != confirm {
		renderRegisterPage(w, "Passwords do not match")
		return
	}

	err := auth.RegisterUser(username, password)
	if err != nil {
		renderRegisterPage(w, "Registration failed: "+err.Error())
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

// Render register page with error
func renderRegisterPage(w http.ResponseWriter, errMsg string) {
	tpl := template.Must(template.ParseFiles(filepath.Join(config.Cfg.TemplatesDir, "register.html")))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.Execute(w, map[string]string{"Error": errMsg})
}

// Logout handler
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil && cookie.Value != "" {
		auth.Logout(cookie.Value)
	}

	// Always clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}
