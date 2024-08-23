package models

import (
	"gorm.io/gorm"
)

type Seller struct {
	gorm.Model
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Password      string    `json:"password"`
	DateOfBirth   string    `json:"date_of_birth"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	StoreName     string    `json:"store_name"`
	StoreCategory string    `json:"store_category"`
	Products      []Product `json:"products"`
}

type LoginRequestSeller struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
