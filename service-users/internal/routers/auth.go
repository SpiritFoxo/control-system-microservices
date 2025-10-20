package routers

import (
	"service-users/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, s *handlers.Server) {
	h := s.UserHandler

	r.POST("/login", h.LoginUser)
}
