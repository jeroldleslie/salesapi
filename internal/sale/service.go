package sale

import (
	"encoding/json"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"salesapi/internal/sale/models"
	"strconv"
)

type Service interface {
	GetTopProducts(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (s *service) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	limitStr := r.URL.Query().Get("limit")
	category := r.URL.Query().Get("category")
	region := r.URL.Query().Get("region")

	if from == "" || to == "" {
		http.Error(w, "Missing required parameters: from and to", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	s.logger.Info("Parsed query parameters",
		zap.String("from", from),
		zap.String("to", to),
		zap.String("category", category),
		zap.String("region", region),
		zap.Int("limit", limit),
	)

	var products []models.TopProduct

	query := s.db.Table("order_items").
		Select("products.product_id, products.product_name, products.category, regions.name as region_name, SUM(order_items.quantity_sold) AS total_sold").
		Joins("JOIN products ON products.product_id = order_items.product_id").
		Joins("JOIN orders ON orders.order_id = order_items.order_id").
		Joins("JOIN regions ON orders.region_id = regions.id").
		Where("orders.date_of_sale BETWEEN ? AND ?", from, to)

	if category != "" {
		query = query.Where("products.category = ?", category)
	}

	if region != "" {
		query = query.Where("regions.name = ?", region)
	}

	err = query.
		Group("products.product_id, products.product_name, products.category, regions.name").
		Order("total_sold DESC").
		Limit(limit).
		Scan(&products).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func NewSaleAnalysisService(db *gorm.DB, logger *zap.Logger) Service {
	return &service{
		db:     db,
		logger: logger,
	}
}
