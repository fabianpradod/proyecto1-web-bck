package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"proyecto1-web-bck/db"
	customMiddleware "proyecto1-web-bck/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db.Connect()

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(customMiddleware.CORS)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
