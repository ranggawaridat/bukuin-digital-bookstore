package book

import "time"

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CoverURL    string    `json:"cover_url"`
	FilePath    string    `json:"file_path"`
	CreatedAt   time.Time `json:"created_at"`
}
