package routers

import (
	"service-users/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup, s *handlers.Server) {
	h := s.UserHandler

	r.POST("/users/register", h.RegisterUser)
	r.GET("/:userId", h.GetUserByID)
	r.PUT("/:userId", h.UpdateUser)
	r.GET("/users", h.GetUsers)
}
