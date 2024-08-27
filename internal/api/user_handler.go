package api

import (
	"e-commerce/internal/middleware"
	"e-commerce/internal/models"
	"e-commerce/internal/util"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Create User
func (u *HTTPHandler) CreateUser(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBind(&user); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	//check if user already exists
	_, err := u.Repository.FindUserByEmail(user.Email)
	if err == nil {
		util.Response(c, "User already exists", 400, "Bad request body", nil)
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	user.Password = hashedPassword

	err = u.Repository.CreateUser(user)
	if err != nil {
		util.Response(c, "User not created", 500, err.Error(), nil)
		return
	}
	util.Response(c, "User created", 200, nil, nil)

}

// Login User
func (u *HTTPHandler) LoginUser(c *gin.Context) {
	var loginRequest *models.LoginRequestUser
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

	user, err := u.Repository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "Email does not exist", 404, err.Error(), nil)
		return
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		util.Response(c, "Invalid password", 400, err.Error(), nil)
		return
	}

	accessClaims, refreshClaims := middleware.GenerateClaims(user.Email)

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

	/* notifications, err := u.Repository.GetNotificationsByUserID(user.ID)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	} */

	c.Header("access_token", *accessToken)
	c.Header("refresh_token", *refreshToken)

	util.Response(c, "Login successful", 200, gin.H{
		"user": user,
		//"notification_details": notifications,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// get all products
func (u *HTTPHandler) GetAllProducts(c *gin.Context) {
	//get useer id from context
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	products, err := u.Repository.GetAllProducts()
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Products fetched", 200, gin.H{
		"products": products,
	}, nil)
}

// add product to cart
func (u *HTTPHandler) AddProductToCart(c *gin.Context) {
	//get user id from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	//bind request to struct
	var cart *models.IndividualItemInCart
	if err := c.ShouldBind(&cart); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	//validate request
	product, err := u.Repository.GetProductByID(cart.ProductID)
	if err != nil {
		util.Response(c, "Product not found", 404, err.Error(), nil)
		return
	}

	//check if product quantity is less
	if cart.Quantity > product.Quantity {
		util.Response(c, "Product quantity is less", 400, nil, nil)
		return
	}

	cart.UserID = user.ID

	err = u.Repository.AddProductToCart(cart)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Product added to cart", 200, nil, nil)
}

// get all products in cart
func (u *HTTPHandler) ViewCart(c *gin.Context) {
	// Get user ID from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	// Get cart items by user ID
	cartItems, err := u.Repository.GetCartsByUserID(user.ID)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}

	// If cart is empty, return early
	if len(cartItems) == 0 {
		util.Response(c, "Your cart is empty", 404, nil, nil)
		return
	}

	// Prepare the response structure
	var cartTotal models.CartTotal
	cartTotal.Cart = make([]*models.CartItem, len(cartItems))

	// Calculate the total price and prepare the cart items
	var total float64
	for i, cartItem := range cartItems {
		product, err := u.Repository.GetProductByID(cartItem.ProductID)
		if err != nil {
			util.Response(c, "Error fetching product details", 500, err.Error(), nil)
			return
		}
		cartTotal.Cart[i] = &models.CartItem{
			CartID:   cartItem.ID,
			Product:  product,
			Quantity: cartItem.Quantity,
		}
		total += float64(cartItem.Quantity) * product.Price
	}

	// Set the total price
	cartTotal.Total = total

	// Return the cart items and total price
	util.Response(c, "Cart fetched successfully", 200, gin.H{
		"cart":  cartTotal.Cart,
		"total": cartTotal.Total,
	}, nil)
}

// place order
func (u *HTTPHandler) PlaceOrder(c *gin.Context) {
	// Get user ID from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	// Get products in cart
	cartItems, err := u.Repository.GetCartsByUserID(user.ID)
	if err != nil {
		util.Response(c, "Error fetching cart items", 500, err.Error(), nil)
		return
	}

	if len(cartItems) == 0 {
		util.Response(c, "Cart is empty", 400, "No items in the cart", nil)
		return
	}

	// Calculate total and prepare order items
	var total float64
	var orderItems []models.OrderItem
	for _, cartItem := range cartItems {
		product, err := u.Repository.GetProductByID(cartItem.ProductID)
		if err != nil {
			util.Response(c, "Error fetching product details", 500, err.Error(), nil)
			return
		}

		// Check if the product is out of stock
		if cartItem.Quantity > product.Quantity {
			util.Response(c, "Product out of stock", 400, "Product is out of stock", nil)
			return
		}

		// Calculate total price
		total += float64(cartItem.Quantity) * product.Price

		// Prepare order item
		orderItems = append(orderItems, models.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
		})
	}

	// Prepare the order
	order := &models.Order{
		UserID: user.ID,
		Total:  total,
		Status: "PLACED",
		Items:  orderItems,
	}

	// Save the order and clear the cart within a transaction
	err = u.Repository.CreateOrder(order)
	if err != nil {
		util.Response(c, "Error creating order", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Order placed successfully", 200, nil, nil)
}

// edit cart
func (u *HTTPHandler) EditCart(c *gin.Context) {
	// Get user id from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	// Bind request to struct
	var cart *models.IndividualItemInCart
	if err := c.ShouldBind(&cart); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	// Get cart by user id
	shoppingCart, err := u.Repository.GetCartItemByProductID(cart.ProductID)
	if err != nil {
		util.Response(c, "Cart not found", 404, err.Error(), nil)
		return
	}

	// Validate request
	product, err := u.Repository.GetProductByID(cart.ProductID)
	if err != nil {
		util.Response(c, "Product not found", 404, err.Error(), nil)
		return
	}

	// Check if product quantity is less
	if product.Quantity < cart.Quantity {
		util.Response(c, "Product quantity is less", 400, nil, nil)
		return
	}

	// Update cart
	cart.UserID = user.ID
	cart.ID = shoppingCart.ID

	err = u.Repository.AddProductToCart(cart)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Cart updated", 200, nil, nil)
}

// delete product from cart
func (u *HTTPHandler) DeleteProductFromCart(c *gin.Context) {
	// Get user id from context
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	// Get product by id
	productID := c.Param("id")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		util.Response(c, "Invalid product ID", 400, err.Error(), nil)
		return
	}

	// Validate request
	shoppingCart, err := u.Repository.GetCartItemByProductID(uint(productIDInt))
	if err != nil {
		util.Response(c, "Product not found", 404, err.Error(), nil)
		return
	}

	err = u.Repository.DeleteProductFromCart(shoppingCart)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Product deleted from cart", 200, nil, nil)
}

// get products by id
func (u *HTTPHandler) GetProductByID(c *gin.Context) {
	// Get user id from context
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	// Get product by id
	productID := c.Param("id")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		util.Response(c, "Invalid product ID", 400, err.Error(), nil)
		return
	}

	product, err := u.Repository.GetProductByID(uint(productIDInt))
	if err != nil {
		util.Response(c, "Product not found", 404, err.Error(), nil)
		return
	}
	util.Response(c, "Product fetched", 200, gin.H{
		"product": product,
	}, nil)
}

// view orders
func (u *HTTPHandler) ViewOrders(c *gin.Context) {
	// Get user id from context
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "Error getting user from context", 500, err.Error(), nil)
		return
	}

	orders, err := u.Repository.GetOrdersByUserID(user.ID)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Orders fetched", 200, gin.H{
		"orders": orders,
	}, nil)
}
