package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	middleware.InitConfig(cfg)
	r := gin.Default()
	routers.Setup(r, cfg)
	r.Run(cfg.Addr)
}
