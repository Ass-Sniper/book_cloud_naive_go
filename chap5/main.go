package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kay/kvstore/auth"
	"github.com/kay/kvstore/handler"
	"github.com/kay/kvstore/logger" // 引入日志模块
	"github.com/kay/kvstore/store"

	"net/http/pprof"
	_ "net/http/pprof"
)

func main() {

	err := auth.LoadUsers()
	if err != nil {
		logger.Log.Fatalf("Failed to load users: %v", err)
	}

	logger.Log.Info("Starting the KV Store Server...")

	dbStore, err := store.NewStore("kvstore.db")
	if err != nil {
		logger.Log.Fatalf("Failed to open DB: %v", err)
	}
	defer dbStore.Close()

	// 启动 TTL 垃圾回收
	go dbStore.StartTTLGC(30 * time.Second) // 每30秒清理一次过期的键值对

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// 注册登录路由（不需要登录）
	handler.RegisterAuthRoutes(r)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Handle("/public/*",
		http.StripPrefix("/public",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 设置 1 天浏览器缓存
				w.Header().Set("Cache-Control", "public, max-age=86400")
				http.FileServer(http.Dir("./public")).ServeHTTP(w, r)
			}),
		),
	)

	// 需认证保护的路由
	r.Group(func(r chi.Router) {

		// 认证中间件
		r.Use(auth.AuthMiddleware)

		// 注册 KV 路由
		handler.RegisterKVRoutes(r, dbStore)

		// 根路径
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "public/index.html")
		})
	})

	// 启动 pprof 服务（监听 :6060）
	go func() {
		logger.Log.Info("pprof running at http://0.0.0.0:6060/debug/pprof/")
		pprofMux := http.NewServeMux()
		pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
		pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		if err := http.ListenAndServe("0.0.0.0:6060", pprofMux); err != nil {
			logger.Log.Fatalf("Failed to start pprof: %v", err)
		}
	}()

	logger.Log.Info("KV store running at http://0.0.0.0:8080")
	err = http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		logger.Log.Fatalf("Failed to start server: %v", err)
	}

}
