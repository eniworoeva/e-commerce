package repository

import (
	"e-commerce/internal/models"
	"e-commerce/internal/ports"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

// NewDB create/returns a new instance of our Database
func NewDB(DB *gorm.DB) ports.Repository {
	return &Postgres{
		DB: DB,
	}
}

// Initialize opens the database, create tables if not created and populate it if its empty and returns a DB
func Initialize(dbURI string) (*gorm.DB, error) {

	conn, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = conn.AutoMigrate(&models.User{}, &models.Seller{}, &models.BlacklistTokens{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.IndividualItemInCart{})
	if err != nil {
		return nil, err
	}
	log.Println("Database connection successful")
	return conn, nil
}
