package cart

type Cart struct {
	ID     int        `json:"id"`
	UserID int        `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	ID       int     `json:"id"`
	BookID   int     `json:"book_id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Price    float64 `json:"price"`
	CoverURL string  `json:"cover_url"`
	Quantity int     `json:"quantity"`
}

type AddItemRequest struct {
	BookID int `json:"book_id"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity"`
}
