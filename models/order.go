package models

import "time"

type Order struct {
	OrderID         string `gorm:"primaryKey"`
	RegionID        uint
	Region          Region
	DateOfSale      time.Time
	CustomerID      string
	Customer        Customer
	PaymentMethodID uint
	PaymentMethod   PaymentMethod
}

type OrderItem struct {
	OrderID      string `gorm:"primaryKey"`
	ProductID    string `gorm:"primaryKey"`
	QuantitySold int
	UnitPrice    float64
	Discount     float64
	ShippingCost float64
}
