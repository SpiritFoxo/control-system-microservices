package routers

import (
	"github.com/SpiritFoxo/control-system-microservices/service-orders/internal/handlers"
	"github.com/SpiritFoxo/control-system-microservices/shared/middleware"
	"github.com/SpiritFoxo/control-system-microservices/shared/userroles"
	"github.com/gin-gonic/gin"
)

func SetupOrdersRoutes(r *gin.RouterGroup, s *handlers.Server) {
	h := s.OrderHandler

	r.GET("/:orderId", h.GetOrderByID)
	r.POST("/", middleware.RoleMiddleware(userroles.RoleEngineer, userroles.RoleObserver), h.CreateOrder)
	r.PATCH("/:orderId", middleware.RoleMiddleware(userroles.RoleManager), h.UpdateOrder)
	r.DELETE("/:orderId", middleware.RoleMiddleware(), h.DeleteOrder)
}
