package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	customMiddleware "proyecto1-web-bck/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"proyecto1-web-bck/db"
	"proyecto1-web-bck/handlers"
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

	// series related routes
	r.Get("/series", handlers.GetSeries)
	r.Get("/series/{id}", handlers.GetSeriesByID)

	r.Post("/series", handlers.CreateSeries)
	r.Put("/series/{id}", handlers.UpdateSeries)
	r.Delete("/series/{id}", handlers.DeleteSeries)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
