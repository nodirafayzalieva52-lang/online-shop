package models

import "time"

type Product struct {
	ID          int		   `json:"id"`
	StoreID     int		   `json:"user_id"`
	CategoryID  int        `json:"category_id"`
	Name		string	   `json:"name"`
	Title       string	   `json:"title"`
	Description string	   `json:"description"`
	Price       float64	   `json:"price"`
	Stock       int		   `json:"stock"`
	Created_At  time.Time  `json:"created_at"`
	Updated_At  time.Time  `json:"updated_at"`
}
