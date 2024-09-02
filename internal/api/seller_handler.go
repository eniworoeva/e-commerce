package api

import (
	"e-commerce/internal/middleware"
	"e-commerce/internal/models"
	"e-commerce/internal/util"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Create Seller
func (u *HTTPHandler) CreateSeller(c *gin.Context) {
	var seller *models.Seller
	if err := c.ShouldBind(&seller); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	_, err := u.Repository.FindSellerByEmail(seller.Email)
	if err == nil {
		util.Response(c, "User already exists", 400, "Bad request body", nil)
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(seller.Password)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	seller.Password = hashedPassword

	err = u.Repository.CreateSeller(seller)
	if err != nil {
		util.Response(c, "Seller not created", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Seller created", 200, nil, nil)

}

// Login Seller
func (u *HTTPHandler) LoginSeller(c *gin.Context) {
	var loginRequest *models.LoginRequestSeller
	err := c.ShouldBind(&loginRequest)
	if err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	loginRequest.Email = strings.TrimSpace(loginRequest.Email)
	loginRequest.Password = strings.TrimSpace(loginRequest.Password)

	if loginRequest.Email == "" {
		util.Response(c, "Email must not be empty", 400, nil, nil)
		return
	}
	if loginRequest.Password == "" {
		util.Response(c, "Password must not be empty", 400, nil, nil)
		return
	}

	Seller, err := u.Repository.FindSellerByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "Email does not exist", 404, err.Error(), nil)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(Seller.Password), []byte(loginRequest.Password))
	if err != nil {
		util.Response(c, "Invalid password", 400, err.Error(), nil)
		return
	}

	accessClaims, refreshClaims := middleware.GenerateClaims(Seller.Email)

	secret := os.Getenv("JWT_SECRET")

	accessToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating access token", 500, err.Error(), nil)
		return
	}

	refreshToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating refresh token", 500, err.Error(), nil)
		return
	}

	c.Header("access_token", *accessToken)
	c.Header("refresh_token", *refreshToken)

	util.Response(c, "Login successful", 200, gin.H{
		"Seller":        Seller,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// create a new product
func (u *HTTPHandler) CreateProduct(c *gin.Context) {
	seller, err := u.GetSellerFromContext(c)
	if err != nil {
		util.Response(c, "Invalid token", 401, err.Error(), nil)
		return
	}

	var product *models.Product
	if err := c.ShouldBind(&product); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	product.SellerID = seller.ID

	err = u.Repository.CreateProduct(product)
	if err != nil {
		util.Response(c, "Product not created", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Product created", 200, nil, nil)
}

// list orders
func (u *HTTPHandler) ListOrders(c *gin.Context) {
	// Get seller ID from context (assuming you store it there after authentication)
	seller, err := u.GetSellerFromContext(c)
	if err != nil {
		util.Response(c, "Error getting seller from context", 500, err.Error(), nil)
		return
	}

	// Fetch all products belonging to the seller
	var products []models.Product
	if err := u.Repository.GetProductsBySellerID(seller.ID, &products); err != nil {
		util.Response(c, "Error fetching seller's products", 500, err.Error(), nil)
		return
	}

	// Collect all order IDs associated with the seller's products
	var orders []models.Order
	for _, product := range products {
		var productOrders []models.Order
		if err := u.Repository.GetOrdersByProductID(product.ID, &productOrders); err != nil {
			util.Response(c, "Error fetching orders for product", 500, err.Error(), nil)
			return
		}
		orders = append(orders, productOrders...)
	}

	// Remove duplicate orders (optional, depending on the structure of your database queries)
	uniqueOrders := util.RemoveDuplicateOrders(orders)

	// Send the response
	util.Response(c, "Orders fetched successfully", 200, uniqueOrders, nil)
}

// Accept the order
func (u *HTTPHandler) AcceptOrder(c *gin.Context) {
	_, err := u.GetSellerFromContext(c)
	if err != nil {
		util.Response(c, "Invalid token", 401, err.Error(), nil)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		util.Response(c, "Order ID not provided", 400, nil, nil)
		return
	}

	//convert id to uint
	orderIDUint, err := util.ConvertStringToUint(orderID)
	if err != nil {
		util.Response(c, "Invalid order ID", 400, err.Error(), nil)
		return
	}

	// Get the order from the database
	order, err := u.Repository.GetOrderByID(orderIDUint)
	if err != nil {
		util.Response(c, "Order not found", 404, err.Error(), nil)
		return
	}

	//check if order is already accepted
	if order.Status != "ACCEPTED" {
		util.Response(c, "Order already accepted", 400, nil, nil)
		return
	}

	// Update the order status to accepted
	order.Status = "ACCEPTED"
	if err := u.Repository.UpdateOrder(order); err != nil {
		util.Response(c, "Error updating order", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Order accepted", 200, nil, nil)
}

// Decline the order
func (u *HTTPHandler) DeclineOrder(c *gin.Context) {
	_, err := u.GetSellerFromContext(c)
	if err != nil {
		util.Response(c, "Invalid token", 401, err.Error(), nil)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		util.Response(c, "Order ID not provided", 400, nil, nil)
		return
	}

	//convert id to uint
	orderIDUint, err := util.ConvertStringToUint(orderID)
	if err != nil {
		util.Response(c, "Invalid order ID", 400, err.Error(), nil)
		return
	}

	// Get the order from the database
	order, err := u.Repository.GetOrderByID(orderIDUint)
	if err != nil {
		util.Response(c, "Order not found", 404, err.Error(), nil)
		return
	}

	//check if order is already declined
	if order.Status != "DECLINED" {
		util.Response(c, "Order already declined", 400, nil, nil)
		return
	}

	// Update the order status to declined
	order.Status = "DECLINED"
	if err := u.Repository.UpdateOrder(order); err != nil {
		util.Response(c, "Error updating order", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Order declined", 200, nil, nil)
}

// delete product
func (u *HTTPHandler) DeleteProduct(c *gin.Context) {
	seller, err := u.GetSellerFromContext(c)
	if err != nil {
		util.Response(c, "Invalid token", 401, err.Error(), nil)
		return
	}

	productID := c.Param("id")
	if productID == "" {
		util.Response(c, "Product ID not provided", 400, nil, nil)
		return
	}

	//convert id to uint
	productIDUint, err := util.ConvertStringToUint(productID)
	if err != nil {
		util.Response(c, "Invalid product ID", 400, err.Error(), nil)
		return
	}

	// Get the product from the database
	product, err := u.Repository.GetProductByID(productIDUint)
	if err != nil {
		util.Response(c, "Product not found", 404, err.Error(), nil)
		return
	}

	// Check if the product belongs to the seller
	if product.SellerID != seller.ID {
		util.Response(c, "Product does not belong to seller", 400, nil, nil)
		return
	}

	// Delete the product
	if err := u.Repository.DeleteProduct(product); err != nil {
		util.Response(c, "Error deleting product", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Product deleted", 200, nil, nil)
}
