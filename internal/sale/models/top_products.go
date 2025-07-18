package models

type TopProduct struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Region      string `json:"region" gorm:"column:region_name"`
	TotalSold   int    `json:"total_sold"`
}
