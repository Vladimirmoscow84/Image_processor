package model

import "time"

type Image struct {
	ID            int       `json:"id" db:"id"`
	OriginalPath  string    `json:"original_path" db:"original_path"`
	ProcessedPath string    `json:"processed_path,omitempty" db:"processed_path"`
	ThumbnailPath string    `json:"thumbnail_path,omitempty" db:"thumbnail_path"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
