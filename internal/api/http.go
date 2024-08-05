package api

import (
	"e-commerce/internal/models"
	"e-commerce/internal/ports"
	"fmt"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	Repository ports.Repository
}

func NewHTTPHandler(repository ports.Repository) *HTTPHandler {
	return &HTTPHandler{
		Repository: repository,
	}
}

func (u *HTTPHandler) GetUserFromContext(c *gin.Context) (*models.User, error) {
	contextUser, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("error getting user from context")
	}
	user, ok := contextUser.(*models.User)
	if !ok {
		return nil, fmt.Errorf("an error occurred")
	}
	return user, nil
}

func (u *HTTPHandler) GetSellerFromContext(c *gin.Context) (*models.Seller, error) {
	contextSeller, exists := c.Get("Seller")
	if !exists {
		return nil, fmt.Errorf("error getting Seller from context")
	}
	seller, ok := contextSeller.(*models.Seller)
	if !ok {
		return nil, fmt.Errorf("an error occurred")
	}
	return seller, nil
}

func (u *HTTPHandler) GetTokenFromContext(c *gin.Context) (string, error) {
	tokenI, exists := c.Get("access_token")
	if !exists {
		return "", fmt.Errorf("error getting access token")
	}
	tokenstr := tokenI.(string)
	return tokenstr, nil
}
