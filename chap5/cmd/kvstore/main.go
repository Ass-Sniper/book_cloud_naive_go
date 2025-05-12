package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"kvstore/internal/auth"
	"kvstore/internal/config"
	"kvstore/internal/handler"
	"kvstore/internal/logger" // 引入日志模块
	"kvstore/internal/store"

	"net/http/pprof"
	_ "net/http/pprof"
)

func main() {

	// 加载配置
	// 定义命令行参数，默认为 config/config.json
	configPath := flag.String("config", "config/config.json", "path to config file")
	flag.Parse()
	err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = auth.LoadUsers()
	if err != nil {
		logger.Log.Fatalf("Failed to load users: %v", err)
	}

	logger.Log.Info("Starting the KV Store Server...")

	dbStore, err := store.NewStore(config.Cfg.DBFile)
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
				http.FileServer(http.Dir(config.Cfg.PublicDir)).ServeHTTP(w, r)
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
			http.ServeFile(w, r, config.Cfg.IndexFile)
		})
	})

	if config.Cfg.EnablePprof {
		// 启动 pprof 服务（默认监听 :6060）
		go func() {
			logger.Log.Infof("pprof running at http://%s/debug/pprof/", config.Cfg.PprofAddr)
			pprofMux := http.NewServeMux()
			pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
			pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

			if err := http.ListenAndServe(config.Cfg.PprofAddr, pprofMux); err != nil {
				logger.Log.Fatalf("Failed to start pprof: %v", err)
			}
		}()
	}

	// 创建 HTTP 服务器对象
	server := &http.Server{
		Addr:    config.Cfg.HTTPAddr,
		Handler: r,
	}

	// 启动主服务（默认监听 :8080）
	go func() {
		logger.Log.Infof("KV store running at http://%s", config.Cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 捕获终止信号 (SIGINT, SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// 等待停止信号
	<-stop

	logger.Log.Info("Shutting down server...")

	// 设置超时进行优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭服务
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatalf("Server Shutdown Failed: %v", err)
	}
	logger.Log.Info("Server exited")

}
