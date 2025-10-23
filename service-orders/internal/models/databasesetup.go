package models

import (
	"log"

	"github.com/SpiritFoxo/control-system-microservices/service-orders/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Can not connect to the database:", err)
	}

	if err := db.AutoMigrate(&Order{}, &OrderItem{}); err != nil {
		return nil, err
	}

	return db, nil
}
