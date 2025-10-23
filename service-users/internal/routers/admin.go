package routers

import (
	"github.com/SpiritFoxo/control-system-microservices/service-users/internal/handlers"

	"github.com/SpiritFoxo/control-system-microservices/service-users/internal/middleware"
	"github.com/SpiritFoxo/control-system-microservices/service-users/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup, s *handlers.Server) {
	h := s.UserHandler

	r.POST("/users/register", middleware.RoleMiddleware(models.RoleSuperadmin, models.RoleAdmin), h.RegisterUser)
	r.GET("/users/:userId", middleware.RoleMiddleware(models.RoleSuperadmin, models.RoleAdmin), h.GetUserByID)
	r.PUT("/users/:userId", middleware.RoleMiddleware(models.RoleSuperadmin, models.RoleAdmin), h.UpdateUser)
	r.GET("/users", middleware.RoleMiddleware(models.RoleSuperadmin, models.RoleAdmin), h.GetUsers)
}
