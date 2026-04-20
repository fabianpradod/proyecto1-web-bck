package models

import "time"

type Rating struct {
	ID        int       `json:"id"`
	SeriesID  int       `json:"series_id"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

type RatingSummary struct {
	SeriesID int     `json:"series_id"`
	Average  float64 `json:"average"`
	Count    int     `json:"count"`
}
