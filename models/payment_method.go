package models

type PaymentMethod struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}
