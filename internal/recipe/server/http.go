package server

import (
    "net/http"
    "fmt"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "cook/internal/recipe/config"
)

func startHTTP() error {
    cfg, err := config.Load()
    if err != nil {
        return err
    }
    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    r.Get("/api/v1/recipes", func(w http.ResponseWriter, _ *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write([]byte("[]"))
    })

    r.Post("/api/v1/query", func(w http.ResponseWriter, _ *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write([]byte(`{"answer":"TODO","sources":[]}`))
    })

    addr := fmt.Sprintf(":%d", cfg.Server.Port)
    return http.ListenAndServe(addr, r)
}