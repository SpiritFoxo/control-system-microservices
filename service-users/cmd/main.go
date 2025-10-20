package main

import (
	"log"
	"service-users/internal/config"
	"service-users/internal/handlers"
	"service-users/internal/models"
	"service-users/internal/routers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DbInit(cfg *config.Config) *gorm.DB {
	db, err := models.Setup(cfg)
	if err != nil {
		log.Println("Connection error")
	}
	return db
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	cfg := config.Load()
	db := DbInit(cfg)
	server := handlers.NewServer(db, cfg)

	api := r.Group("api/v1")

	admin := api.Group("/admin")
	routers.RegisterAdminRoutes(admin, server)

	auth := api.Group("/auth")
	routers.RegisterAuthRoutes(auth, server)

	return r
}

func main() {
	r := SetupRouter()
	r.Run(":8082")
}
