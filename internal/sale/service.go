package sale

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"salesapi/internal/sale/entities"
	"salesapi/internal/sale/repo"
	"strconv"
)

type Service interface {
	GetTopProducts(w http.ResponseWriter, r *http.Request)
}

type service struct {
	repo   repo.Repository
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

	products, err := s.repo.GetTopProducts(entities.GetTopProductsRequest{
		From:     from,
		To:       to,
		Limit:    limit,
		Category: category,
		Region:   region,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func NewSaleAnalysisService(repo repo.Repository, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}
