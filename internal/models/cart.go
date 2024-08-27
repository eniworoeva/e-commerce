package models

import "gorm.io/gorm"

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
