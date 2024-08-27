package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID uint        `json:"user_id"`
	Items  []OrderItem `json:"items"`
	Total  float64     `json:"total"`
	Status OrderStatus `json:"status"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id" gorm:"uniqueIndex:idx_order_product"`
	ProductID uint    `json:"product_id" gorm:"uniqueIndex:idx_order_product"`
	Quantity  int     `json:"quantity"`
	Product   Product `json:"product"`
}

type OrderStatus string

const (
	PLACED    OrderStatus = "PLACED"
	ACCEPTED  OrderStatus = "ACCEPTED"
	COMPLETED OrderStatus = "COMPLETED"
	DECLINED  OrderStatus = "DECLINED"
)
