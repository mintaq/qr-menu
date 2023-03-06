package models

// Cart struct to describe cart object.
type Cart struct {
	UserToken  string    `json:"user_token"`
	Items      []Product `json:"items"`
	TotalPrice float64   `json:"total_price"`
	Note       string    `json:"note"`
	ItemsCount uint      `json:"items_count"`
}
