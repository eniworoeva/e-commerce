package ports

import "e-commerce/internal/models"

type Repository interface {
	FindUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
	FindAllUsers() ([]models.User, error)
	FindSellerByEmail(email string) (*models.Seller, error)
	CreateUser(user *models.User) error
	CreateSeller(Seller *models.Seller) error
	UpdateUser(user *models.User) error
	UpdateSeller(user *models.Seller) error
	BlacklistToken(token *models.BlacklistTokens) error
	TokenInBlacklist(token *string) bool
	GetAllProducts() ([]models.Product, error)
	GetProductByID(productID uint) (*models.Product, error)
	AddProductToCart(cart *models.IndividualItemInCart) error
	GetCartsByUserID(userID uint) ([]*models.IndividualItemInCart, error)
	CreateOrder(order *models.Order) error
	CreateProduct(product *models.Product) error
	DeleteProductFromCart(cart *models.IndividualItemInCart) error
	GetOrdersByUserID(userID uint) ([]*models.Order, error)
	GetCartItemByProductID(productID uint) (*models.IndividualItemInCart, error)
	ListOrders(sellerID uint) ([]*models.Order, error)
	GetProductsBySellerID(sellerID uint, products *[]models.Product) error
	GetOrdersByProductID(productID uint, orders *[]models.Order) error
	GetOrderByID(orderID uint) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteProduct(product *models.Product) error
	GetOrderItemsByOrderID(orderID uint) ([]*models.OrderItem, error)
}
