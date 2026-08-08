package models

import "time"

type Store struct {
	ID          int			`json:"id"`
	Seller_ID   int64		`json:"seller_id"`
	Name        string		`json:"name"`
	Description string		`json:"description"`
	CreatedAt  time.Time 	`json:"created_at"`
}
