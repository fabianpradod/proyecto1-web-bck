package models

import "time"

type Series struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Genre     string    `json:"genre"`
	Status    string    `json:"status"`
	Episodes  int       `json:"episodes"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}
