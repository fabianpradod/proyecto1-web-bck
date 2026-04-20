package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"proyecto1-web-bck/db"
	"proyecto1-web-bck/models"

	"github.com/go-chi/chi/v5"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func GetSeries(w http.ResponseWriter, r *http.Request) {
	// query params
	q := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// defaults
	if sort == "" {
		sort = "id"
	}
	validSorts := map[string]bool{"id": true, "name": true, "genre": true, "episodes": true, "created_at": true}
	if !validSorts[sort] {
		sort = "id"
	}
	if order != "desc" {
		order = "asc"
	}

	page := 1
	limit := 10
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	offset := (page - 1) * limit

	query := `SELECT id, name, genre, status, episodes, image_url, created_at FROM series WHERE 1=1`
	args := []any{}
	argIdx := 1

	if q != "" {
		query += ` AND name ILIKE $` + strconv.Itoa(argIdx)
		args = append(args, "%"+q+"%")
		argIdx++
	}

	query += ` ORDER BY ` + sort + ` ` + order
	query += ` LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error fetching series")
		return
	}
	defer rows.Close()

	series := []models.Series{}
	for rows.Next() {
		var s models.Series
		if err := rows.Scan(&s.ID, &s.Name, &s.Genre, &s.Status, &s.Episodes, &s.ImageURL, &s.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "Error scanning series")
			return
		}
		series = append(series, s)
	}

	writeJSON(w, http.StatusOK, series)
}

func GetSeriesByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var s models.Series
	err = db.DB.QueryRow(
		`SELECT id, name, genre, status, episodes, image_url, created_at FROM series WHERE id = $1`, id,
	).Scan(&s.ID, &s.Name, &s.Genre, &s.Status, &s.Episodes, &s.ImageURL, &s.CreatedAt)

	if err != nil {
		writeError(w, http.StatusNotFound, "Series not found")
		return
	}

	writeJSON(w, http.StatusOK, s)
}

func CreateSeries(w http.ResponseWriter, r *http.Request) {
	var s models.Series
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if s.Name == "" {
		writeError(w, http.StatusBadRequest, "Name is required")
		return
	}

	err := db.DB.QueryRow(
		`INSERT INTO series (name, genre, status, episodes, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		s.Name, s.Genre, s.Status, s.Episodes, s.ImageURL,
	).Scan(&s.ID, &s.CreatedAt)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error creating series")
		return
	}

	writeJSON(w, http.StatusCreated, s)
}

func UpdateSeries(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var s models.Series
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if s.Name == "" {
		writeError(w, http.StatusBadRequest, "Name is required")
		return
	}

	result, err := db.DB.Exec(
		`UPDATE series SET name=$1, genre=$2, status=$3, episodes=$4, image_url=$5 WHERE id=$6`,
		s.Name, s.Genre, s.Status, s.Episodes, s.ImageURL, id,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error updating series")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "Series not found")
		return
	}

	s.ID = id
	writeJSON(w, http.StatusOK, s)
}

func DeleteSeries(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	result, err := db.DB.Exec(`DELETE FROM series WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error deleting series")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "Series not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
