package models

import (
	"gorm.io/gorm"
)

type BlacklistTokens struct {
	gorm.Model
	Token string `json:"token" gorm:"not null"`
}
