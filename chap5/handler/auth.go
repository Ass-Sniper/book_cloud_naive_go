package handler

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kay/kvstore/auth"
)

func RegisterAuthRoutes(r chi.Router) {
	r.Get("/login", loginPage)
	r.Post("/login", loginHandler)
	r.Get("/register", registerPage)
	r.Post("/register", registerHandler)
	r.Get("/logout", logoutHandler)
}

// 登录页面
func loginPage(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "模板加载失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.Execute(w, nil)
}

// 登录处理
func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if !auth.ValidateUser(username, password) {
		tpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "模板加载失败", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tpl.Execute(w, map[string]string{"Error": "用户名或密码错误"})
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

// 注册页面
func registerPage(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("templates/register.html"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.Execute(w, nil)
}

// 注册处理
func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	if password != confirm {
		renderRegisterPage(w, "两次输入的密码不一致")
		return
	}

	err := auth.RegisterUser(username, password)
	if err != nil {
		renderRegisterPage(w, "注册失败: "+err.Error())
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

// 抽出错误渲染逻辑
func renderRegisterPage(w http.ResponseWriter, errMsg string) {
	tpl := template.Must(template.ParseFiles("templates/register.html"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.Execute(w, map[string]string{"Error": errMsg})
}

// 注销处理
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil && cookie.Value != "" {
		auth.Logout(cookie.Value)
	}

	// 永远清除 Cookie，无论是否存在有效 session
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}
