package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"auth-proxy/internal/auth"
	"auth-proxy/internal/config"
	"auth-proxy/internal/middleware"
	"auth-proxy/internal/proxy"
	"auth-proxy/internal/token"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	authProvider := auth.NewStaticProvider(cfg.Users)

	tokenService := token.NewService(
		cfg.Auth.JWTSecret,
		cfg.Auth.CookieName,
		cfg.Auth.CookieSecure,
		cfg.Auth.CookieMaxAge,
		cfg.Auth.TokenDuration,
	)

	proxyHandler, err := proxy.NewHandler(cfg.Backend.URL)
	if err != nil {
		log.Fatalf("Failed to create proxy handler: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login", loginHandler(authProvider, tokenService))
	mux.HandleFunc("/logout", logoutHandler(tokenService))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.HandleFunc("/login-page", loginPageHandler())

	protectedHandler := middleware.Auth(tokenService)(proxyHandler)
	mux.Handle("/", protectedHandler)

	log.Printf("Starting server on %s", cfg.Server.Listen)
	log.Printf("Proxying to %s", cfg.Backend.URL)
	if err := http.ListenAndServe(cfg.Server.Listen, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func loginPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/login.html")
	}
}

func loginHandler(provider auth.Provider, tokenService *token.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Redirect(w, r, "/login-page", http.StatusFound)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		user, err := provider.Authenticate(r.Context(), req.Username, req.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
			return
		}

		tokenStr, err := tokenService.GenerateToken(user.Username, user.Roles)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		tokenService.SetCookie(w, tokenStr)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"redirect": "/"})
	}
}

func logoutHandler(tokenService *token.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenService.ClearCookie(w)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func usage() {
	fmt.Println("Usage: server [config.yaml]")
	fmt.Println("Environment variable CONFIG_PATH can also be used to specify config file")
	os.Exit(1)
}
