#!/bin/bash

# 项目名称
PROJECT_NAME="kvstore"

# 创建项目目录结构
echo "创建项目目录结构..."

mkdir -p $PROJECT_NAME/{auth,handler,store,public}
touch $PROJECT_NAME/{go.mod,go.sum}
touch $PROJECT_NAME/main.go

# 创建 auth 目录下的文件
touch $PROJECT_NAME/auth/{session.go,middleware.go}

# 创建 handler 目录下的文件
touch $PROJECT_NAME/handler/{kv.go,auth.go}

# 创建 store 目录下的文件
touch $PROJECT_NAME/store/boltdb.go

# 创建 public 目录下的 index.html 文件
cat <<EOL > $PROJECT_NAME/public/index.html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>KV 存储系统</title>
  <style>
    body { font-family: Arial; margin: 40px; }
    table { border-collapse: collapse; width: 60%; }
    th, td { border: 1px solid #ccc; padding: 8px; }
    input, button { margin: 5px; padding: 6px; }
  </style>
</head>
<body>
  <h1>KV 存储系统 Web UI</h1>

  <h2>➕ 添加键值对</h2>
  <form id="addForm">
    Key: <input type="text" id="key" required>
    Value: <input type="text" id="value" required>
    TTL(秒): <input type="number" id="ttl" value="0" min="0">
    <button type="submit">添加</button>
  </form>

  <h2>📋 当前存储</h2>
  <button onclick="loadKeys()">🔄 刷新列表</button>
  <table id="kvTable">
    <thead>
      <tr><th>Key</th><th>Value</th><th>操作</th></tr>
    </thead>
    <tbody></tbody>
  </table>

  <script>
    const API = '/kv';

    async function loadKeys() {
      const tbody = document.querySelector('#kvTable tbody');
      tbody.innerHTML = '';
      const keys = await fetch(API).then(r => r.json());
      for (let key of keys) {
        const res = await fetch(`${API}/${key}`);
        if (!res.ok) continue;
        const val = await res.json();
        const row = document.createElement('tr');
        row.innerHTML = `
          <td>${key}</td>
          <td>${val.value}</td>
          <td><button onclick="deleteKey('${key}')">删除</button></td>
        `;
        tbody.appendChild(row);
      }
    }

    async function deleteKey(key) {
      await fetch(`${API}/${key}`, { method: 'DELETE' });
      loadKeys();
    }

    document.getElementById('addForm').addEventListener('submit', async e => {
      e.preventDefault();
      const key = document.getElementById('key').value;
      const value = document.getElementById('value').value;
      const ttl = parseInt(document.getElementById('ttl').value);
      await fetch(`${API}/${key}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ value, ttl })
      });
      e.target.reset();
      loadKeys();
    });

    loadKeys();
  </script>
</body>
</html>
EOL

# 创建 main.go 文件
cat <<EOL > $PROJECT_NAME/main.go
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"kvstore/handler"
	"kvstore/store"
	"kvstore/auth"
)

func main() {
	dbStore, err := store.NewStore("kvstore.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer dbStore.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// 注册登录路由（不需要登录）
	handler.RegisterAuthRoutes(r)

	// 需认证保护
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		handler.RegisterKVRoutes(r, dbStore)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "public/index.html")
		})
	})

	log.Println("KV store running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
EOL

# 创建 auth/session.go 文件
cat <<EOL > $PROJECT_NAME/auth/session.go
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

var (
	users = map[string]string{
		"admin": "123456", // 预设用户名密码
	}

	sessions = map[string]string{}
	mu       sync.Mutex
)

func ValidateUser(username, password string) bool {
	pass, ok := users[username]
	return ok && pass == password
}

func CreateSession(username string) string {
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)

	mu.Lock()
	defer mu.Unlock()
	sessions[token] = username

	// 可扩展为过期策略（此处略）
	return token
}

func GetUsernameByToken(token string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()
	u, ok := sessions[token]
	return u, ok
}
EOL

# 创建 auth/middleware.go 文件
cat <<EOL > $PROJECT_NAME/auth/middleware.go
package auth

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if _, ok := GetUsernameByToken(cookie.Value); !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
EOL

# 创建 handler/kv.go 文件
cat <<EOL > $PROJECT_NAME/handler/kv.go
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"kvstore/store"
)

type KVRequest struct {
	Value string `json:"value"`
	TTL   int64  `json:"ttl"`
}

type KVResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func RegisterKVRoutes(r chi.Router, s *store.Store) {
	r.Route("/kv", func(r chi.Router) {
		r.Get("/", listKeysHandler(s))
		r.Put("/{key}", setHandler(s))
		r.Get("/{key}", getHandler(s))
		r.Delete("/{key}", deleteHandler(s))
	})
}

func getHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		val, ok := s.Get(key)
		if !ok {
			http.Error(w, "Key not found or expired", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(KVResponse{Key: key, Value: val})
	}
}

func setHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		var req KVRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		err := s.Set(key, req.Value, req.TTL)
		if err != nil {
			http.Error(w, "Failed to set value", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(KVResponse{Key: key, Value: req.Value})
	}
}

func deleteHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if err := s.Delete(key); err != nil {
			http.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func listKeysHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, err := s.ListKeys()
		if err != nil {
			http.Error(w, "Failed to list keys", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(keys)
	}
}
EOL

# 创建 handler/auth.go 文件
cat <<EOL > $PROJECT_NAME/handler/auth.go
package handler

import (
	"html/template"
	"net/http"

	"kvstore/auth"
)

func RegisterAuthRoutes(r chi.Router) {
r.Get("/login", loginPage)
r.Post("/login", loginHandler)
}

var loginTpl = `<html><body>

<h2>登录 KV 系统</h2> <form method="POST" action="/login"> 用户名: <input name="username"><br> 密码: <input name="password" type="password"><br> <button type="submit">登录</button> </form> </body></html>`
func loginPage(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/html")
w.Write([]byte(loginTpl))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
username := r.FormValue("username")
password := r.FormValue("password")
if !auth.ValidateUser(username, password) {
http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
return
}
token := auth.CreateSession(username)
http.SetCookie(w, &http.Cookie{
Name: "session",
Value: token,
Path: "/",
HttpOnly: true,
})
http.Redirect(w, r, "/", http.StatusFound)
}
EOL

echo "项目目录结构创建完成！"

