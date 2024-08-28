package util

import (
	"e-commerce/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Response is customized to help return all responses need
func Response(c *gin.Context, message string, status int, data interface{}, errs []string) {
	responsedata := gin.H{
		"message":   message,
		"data":      data,
		"errors":    errs,
		"status":    http.StatusText(status),
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	c.IndentedJSON(status, responsedata)
}

// HashPassword takes a plaintext password and returns the hashed password or an error
func HashPassword(password string) (string, error) {

	// Use bcrypt to generate a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ConvertStringToUint(s string) (uint, error) {
	int, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint(int), nil
}

// Helper function to remove duplicate orders
func RemoveDuplicateOrders(orders []models.Order) []models.Order {
	seen := make(map[uint]bool)
	var uniqueOrders []models.Order

	for _, order := range orders {
		if _, exists := seen[order.ID]; !exists {
			seen[order.ID] = true
			uniqueOrders = append(uniqueOrders, order)
		}
	}

	return uniqueOrders
}
