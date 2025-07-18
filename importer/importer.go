package importer

import (
	"encoding/csv"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"salesapi/models"
	"strconv"
	"strings"
	"time"
)

type Importer interface {
	StartWorker(manager *RefreshManager)
	RefreshHandler(manager *RefreshManager) http.HandlerFunc
}

type service struct {
	db     *gorm.DB
	config *Config
	logger *zap.Logger
}

// handler.go
func (s *service) RefreshHandler(manager *RefreshManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !manager.TryLock() {
			http.Error(w, "Refresh already in progress", http.StatusConflict)
			return
		}

		go func() {
			manager.TryLock()
			defer manager.Unlock()
			s.logger.Info("Manual refresh started")
			err := s.refreshData()
			if err != nil {
				s.logger.Error("Manual refresh failed", zap.Error(err))
			}
		}()

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Refresh started"))
	}
}

func (s *service) StartWorker(manager *RefreshManager) {
	// Initial run
	go s.triggerRefresh(manager)

	// Background worker. This triggers data refresh every 24 hour
	ticker := time.NewTicker(time.Duration(s.config.Interval) * time.Hour)
	go func() {
		for range ticker.C {
			s.triggerRefresh(manager)
		}
	}()
}

func NewImporterService(db *gorm.DB, config *Config, logger *zap.Logger) Importer {
	return &service{
		db:     db,
		config: config,
		logger: logger,
	}
}

func (s *service) refreshData() error {
	file, err := os.Open(s.config.CSVPath)
	if err != nil {
		s.logger.Error("Failed to open CSV file", zap.String("path", s.config.CSVPath), zap.Error(err))
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		s.logger.Error("Failed to read CSV rows", zap.Error(err))
		return err
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		orderID := row[0]
		productID := row[1]
		customerID := row[2]
		productName := row[3]
		category := row[4]
		regionName := row[5]
		dateOfSale, _ := time.Parse("2006-01-02", row[6])
		quantitySold, _ := strconv.Atoi(row[7])
		unitPrice, _ := strconv.ParseFloat(row[8], 64)
		discount, _ := strconv.ParseFloat(row[9], 64)
		shippingCost, _ := strconv.ParseFloat(row[10], 64)
		paymentMethodName := row[11]
		customerName := row[12]
		customerEmail := row[13]
		customerAddress := strings.Trim(row[14], "\"")

		region := s.getOrCreateRegion(regionName)
		paymentMethod := s.getOrCreatePaymentMethod(paymentMethodName)
		customer := s.getOrCreateCustomer(customerID, customerName, customerEmail, customerAddress)
		product := s.getOrCreateProduct(productID, productName, category)

		s.db.Save(&models.Order{
			OrderID:         orderID,
			RegionID:        region.ID,
			DateOfSale:      dateOfSale,
			CustomerID:      customer.CustomerID,
			PaymentMethodID: paymentMethod.ID,
		})

		s.db.Save(&models.OrderItem{
			OrderID:      orderID,
			ProductID:    product.ProductID,
			QuantitySold: quantitySold,
			UnitPrice:    unitPrice,
			Discount:     discount,
			ShippingCost: shippingCost,
		})
	}

	s.logger.Info("CSV import completed successfully", zap.String("timestamp", time.Now().Format(time.RFC3339)))

	return nil
}

func (s *service) triggerRefresh(manager *RefreshManager) {
	if !manager.TryLock() {
		s.logger.Info("Refresh already in progress — skipping")
		return
	}

	s.logger.Info("Starting data refresh worker...")
	err := s.refreshData()
	if err != nil {
		s.logger.Error("Refresh failed", zap.Error(err))
	}
	manager.Unlock()
}

// Get or create reusable records
func (s *service) getOrCreateRegion(name string) models.Region {
	var region models.Region
	s.db.FirstOrCreate(&region, models.Region{Name: name})
	return region
}

func (s *service) getOrCreatePaymentMethod(name string) models.PaymentMethod {
	var method models.PaymentMethod
	s.db.FirstOrCreate(&method, models.PaymentMethod{Name: name})
	return method
}

func (s *service) getOrCreateCustomer(id, name, email, address string) models.Customer {
	var customer models.Customer
	s.db.FirstOrCreate(&customer, models.Customer{
		CustomerID: id,
	}, models.Customer{
		CustomerName:    name,
		CustomerEmail:   email,
		CustomerAddress: address,
	})
	return customer
}

func (s *service) getOrCreateProduct(id, name, category string) models.Product {
	var product models.Product
	s.db.FirstOrCreate(&product, models.Product{
		ProductID: id,
	}, models.Product{
		ProductName: name,
		Category:    category,
	})
	return product
}
