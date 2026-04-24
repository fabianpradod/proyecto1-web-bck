package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"proyecto1-web-bck/db"
	"proyecto1-web-bck/models"

	"github.com/go-chi/chi/v5"
)

func CreateRating(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid series ID")
		return
	}

	var rating models.Rating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if rating.Score < 1 || rating.Score > 10 {
		writeError(w, http.StatusBadRequest, "Score must be between 1 and 10")
		return
	}

	rating.SeriesID = id
	err = db.DB.QueryRow(
		`INSERT INTO ratings (series_id, score) VALUES ($1, $2) RETURNING id, created_at`,
		rating.SeriesID, rating.Score,
	).Scan(&rating.ID, &rating.CreatedAt)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error saving rating")
		return
	}

	writeJSON(w, http.StatusCreated, rating)
}

func GetRating(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid series ID")
		return
	}

	var summary models.RatingSummary
	err = db.DB.QueryRow(
		`SELECT series_id, COALESCE(AVG(score), 0), COUNT(*) FROM ratings WHERE series_id = $1 GROUP BY series_id`,
		id,
	).Scan(&summary.SeriesID, &summary.Average, &summary.Count)

	if err != nil {
		// no ratings yet return default summary
		summary = models.RatingSummary{SeriesID: id, Average: 0, Count: 0}
	}

	writeJSON(w, http.StatusOK, summary)
}
