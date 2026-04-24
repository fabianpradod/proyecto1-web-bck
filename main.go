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

	// ratings related routes
	r.Post("/series/{id}/rating", handlers.CreateRating)
	r.Get("/series/{id}/rating", handlers.GetRating)

	// docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
<html>
<head>
	<title>Series Tracker API</title>
	<meta charset="utf-8"/>
	<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
SwaggerUIBundle({
	url: "/openapi.yaml",
	dom_id: '#swagger-ui',
	presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
	layout: "BaseLayout"
})
</script>
</body>
</html>`
		fmt.Fprint(w, html)
	})

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		http.ServeFile(w, r, "docs/openapi.yaml")
	})

	// this must always be last
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
