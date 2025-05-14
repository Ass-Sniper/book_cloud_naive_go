package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kvstore/internal/auth"
	"kvstore/internal/config"
	"kvstore/internal/handler"
	"kvstore/internal/i18n"
	"kvstore/internal/logger"
	"kvstore/internal/store"
	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// loadConfiguration loads the configuration file and returns an error if any.
func loadConfiguration() error {
	configPath := flag.String("config", "config/config.json", "path to config file")
	flag.Parse()
	if err := config.LoadConfig(*configPath); err != nil {
		return err
	}
	return nil
}

// loadUsers loads user data and returns an error if any.
func loadUsers() error {
	if err := auth.LoadUsers(); err != nil {
		return err
	}
	return nil
}

// initializeStore initializes the store and returns it, or an error if any.
func initializeStore() (*store.Store, error) {
	dbStore, err := store.NewStore(config.Cfg.DBFile)
	if err != nil {
		return nil, err
	}
	return dbStore, nil
}

// startTTLGC starts a goroutine for TTL garbage collection.
func startTTLGC(dbStore *store.Store) {
	go dbStore.StartTTLGC(30 * time.Second) // Every 30 seconds to clean expired keys
}

// configureRouter configures the routes and middlewares for the HTTP server.
func configureRouter(dbStore *store.Store) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Register public routes
	handler.RegisterAuthRoutes(r)

	// Health check route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Static file handling for /public
	r.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir(config.Cfg.PublicDir))))

	// Routes that require authentication
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		handler.RegisterKVRoutes(r, dbStore)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, config.Cfg.IndexFile)
		})
	})

	return r
}

// startPprof starts the pprof server if enabled in the configuration.
func startPprof() {
	if config.Cfg.EnablePprof {
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
}

// startServer starts the HTTP server and listens on the configured address.
func startServer(r *chi.Mux) *http.Server {
	server := &http.Server{
		Addr:    config.Cfg.HTTPAddr,
		Handler: r,
	}

	go func() {
		logger.Log.Infof("KV store running at http://%s", config.Cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return server
}

// shutdownServer gracefully shuts down the server with a timeout.
func shutdownServer(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for stop signal
	<-stop

	logger.Log.Info("Shutting down server...")

	// Setting up timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatalf("Server Shutdown Failed: %v", err)
	}
	logger.Log.Info("Server exited")
}

func main() {
	// Load configuration
	if err := loadConfiguration(); err != nil {
		logger.Log.Fatalf("Failed to load config: %v", err)
	}

	if err := i18n.InitializeTranslator(); err != nil {
		logger.Log.Fatalf("Failed to initialize translator: %v", err)
	}

	// Load users for authentication
	if err := loadUsers(); err != nil {
		logger.Log.Fatalf("Failed to load users: %v", err)
	}

	// Initialize the key-value store
	dbStore, err := initializeStore()
	if err != nil {
		logger.Log.Fatalf("Failed to open DB: %v", err)
	}
	defer dbStore.Close()

	// Start TTL garbage collection
	startTTLGC(dbStore)

	// Configure the router
	r := configureRouter(dbStore)

	// Start pprof server if enabled
	startPprof()

	// Start the main HTTP server
	server := startServer(r)

	// Gracefully shutdown the server
	shutdownServer(server)
}
