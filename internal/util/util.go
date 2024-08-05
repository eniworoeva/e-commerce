package util

import (
	"net/http"
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
