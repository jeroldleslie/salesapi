package repo

import (
	"gorm.io/gorm"
	"salesapi/internal/sale/models"
)

type GetTopProductsRequest struct {
	From     string
	To       string
	Limit    int
	Category string
	Region   string
}

type Repository interface {
	GetTopProducts(req GetTopProductsRequest) ([]models.TopProduct, error)
}

type repository struct {
	db *gorm.DB
}

func NewSaleAnalysisRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (s *repository) GetTopProducts(req GetTopProductsRequest) ([]models.TopProduct, error) {
	var products []models.TopProduct

	query := s.db.Table("order_items").
		Select("products.product_id, products.product_name, products.category, regions.name as region_name, SUM(order_items.quantity_sold) AS total_sold").
		Joins("JOIN products ON products.product_id = order_items.product_id").
		Joins("JOIN orders ON orders.order_id = order_items.order_id").
		Joins("JOIN regions ON orders.region_id = regions.id").
		Where("orders.date_of_sale BETWEEN ? AND ?", req.From, req.To)

	if req.Category != "" {
		query = query.Where("products.category = ?", req.Category)
	}

	if req.Region != "" {
		query = query.Where("regions.name = ?", req.Region)
	}

	err := query.
		Group("products.product_id, products.product_name, products.category, regions.name").
		Order("total_sold DESC").
		Limit(req.Limit).
		Scan(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
