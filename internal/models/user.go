package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string `json:"first_name"`
	Lastname    string `json:"last_name"`
	Password    string `json:"password"`
	DateOfBirth string `json:"date_of_birth"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
}

type LoginRequestUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
