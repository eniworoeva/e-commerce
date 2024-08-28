package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	SellerID    uint    `json:"seller_id"`
	Title       string  `json:"title"`
	ImageUrl    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Overview    string  `json:"overview"`
	Description string  `json:"description"`
	Status      bool    `json:"status"`
	Orders      []Order `json:"orders" gorm:"many2many:order_items;"`
}
