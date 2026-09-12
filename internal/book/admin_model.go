package book

type CreateBookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CoverURL    string `json:"cover_url"`
	FilePath    string `json:"file_path"`
}
