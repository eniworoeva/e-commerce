package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	SellerID    uint    `json:"merchant_id"`
	Title       string  `json:"title"`
	ImageUrl    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Overview    string  `json:"overview"`
	Description string  `json:"description"`
	Status      bool    `json:"status"`
	Stock       int     `json:"stock"`
}

type IndividualItemInCart struct {
	gorm.Model
	UserID    uint  `json:"user_id"`
	ProductID uint  `json:"product_id"`
	Quantity  int   `json:"quantity"`
	OrderID   *uint `json:"order_id" gorm:"default:null"`
}

type CartItem struct {
	CartID   uint     `json:"cart_id"`
	Product  *Product `json:"product"`
	Quantity int      `json:"quantity"`
}

type CartTotal struct {
	Cart  []*CartItem `json:"cart"`
	Total float64     `json:"total"`
}

type Order struct {
	gorm.Model
	UserID uint                   `json:"user_id"`
	Items  []IndividualItemInCart `json:"items" gorm:"foreignKey:OrderID"`
	Total  float64                `json:"total"`
	Status string                 `json:"status"`
}
