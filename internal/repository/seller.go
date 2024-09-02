package repository

import "e-commerce/internal/models"

func (p *Postgres) FindSellerByEmail(email string) (*models.Seller, error) {
	seller := &models.Seller{}

	if err := p.DB.Where("email = ?", email).First(&seller).Error; err != nil {
		return nil, err
	}
	return seller, nil
}

// Create a user in the database
func (p *Postgres) CreateSeller(seller *models.Seller) error {
	if err := p.DB.Create(seller).Error; err != nil {
		return err
	}
	return nil
}

// Update a user in the database
func (p *Postgres) UpdateSeller(seller *models.Seller) error {
	if err := p.DB.Save(seller).Error; err != nil {
		return err
	}
	return nil
}

// create a product in the database
func (p *Postgres) CreateProduct(product *models.Product) error {
	if err := p.DB.Create(product).Error; err != nil {
		return err
	}
	return nil
}

// list orders
func (p *Postgres) ListOrders(sellerID uint) ([]*models.Order, error) {
	orders := []*models.Order{}

	if err := p.DB.Where("seller_id = ?", sellerID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (p *Postgres) GetProductsBySellerID(sellerID uint, products *[]models.Product) error {
	return p.DB.Where("seller_id = ?", sellerID).Find(products).Error
}

// GetOrdersByProductID retrieves orders associated with a specific product ID,
// including the related order items and products.
func (p *Postgres) GetOrdersByProductID(productID uint, orders *[]models.Order) error {
	return p.DB.
		Preload("Items").         // Preload the OrderItems related to the order
		Preload("Items.Product"). // Preload the associated Product for each OrderItem
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Where("order_items.product_id = ?", productID).
		Find(orders).Error
}

func (p *Postgres) GetOrderByID(orderID uint) (*models.Order, error) {
	order := &models.Order{}

	if err := p.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (p *Postgres) UpdateOrder(order *models.Order) error {
	if err := p.DB.Save(order).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) DeleteProduct(product *models.Product) error {
	if err := p.DB.Delete(product).Error; err != nil {
		return err
	}
	return nil
}
