package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId uint   `gorm:"not null"`
	Status string `gorm:"not null"`
	Cost   int    `gorm:"not null"`
	Items  []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderId  uint   `gorm:"not null"`
	Order    Order  `gorm:"foreignKey:OrderId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity int    `gorm:"not null"`
	Name     string `gorm:"not null"`
}
