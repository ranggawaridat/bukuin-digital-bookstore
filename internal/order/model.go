package order

import "time"

type Order struct {
	ID            int         `json:"id"`
	UserID        int         `json:"user_id"`
	TotalAmount   float64     `json:"total_amount"`
	Status        string      `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	PaymentToken  string      `json:"payment_token,omitempty"`
	PaymentURL    string      `json:"payment_url,omitempty"`
	PaymentMethod string      `json:"payment_method,omitempty"`
	TransactionID string      `json:"transaction_id,omitempty"`
	PaidAt        *time.Time  `json:"paid_at,omitempty"`
	Items         []OrderItem `json:"items"`
}

type OrderItem struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	BookID   int     `json:"book_id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
	FilePath string  `json:"file_path,omitempty"`
}
