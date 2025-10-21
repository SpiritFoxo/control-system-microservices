package routers

import (
	"service-users/internal/handlers"
	"service-users/internal/middleware"
	"service-users/internal/models"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup, s *handlers.Server) {
	h := s.UserHandler

	r.POST("/users/register", middleware.RoleMiddleware(models.RoleSuperadmin, models.RoleAdmin), h.RegisterUser)
	r.GET("/:userId", middleware.RoleMiddleware(models.RoleSuperadmin, models.RoleAdmin), h.GetUserByID)
	r.PUT("/:userId", middleware.RoleMiddleware(models.RoleSuperadmin, models.RoleAdmin), h.UpdateUser)
	r.GET("/users", middleware.RoleMiddleware(models.RoleSuperadmin, models.RoleAdmin), h.GetUsers)
}
