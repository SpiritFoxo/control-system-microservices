package handlers

import (
	"errors"
	"net/http"
	"service-users/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func response(c *gin.Context, status int, success bool, data interface{}, err error) {
	if err != nil {
		c.JSON(status, gin.H{
			"success": success,
			"error": gin.H{
				"code":    status,
				"message": err.Error(),
			},
		})
		return
	}
	c.JSON(status, gin.H{
		"success": success,
		"data":    data,
	})
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input services.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response(c, http.StatusBadRequest, false, nil, err)
		return
	}

	user, err := h.service.RegisterUser(input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" || err.Error() == "invalid email format" ||
			err.Error() == "password must be at least 8 characters" ||
			err.Error() == "email, password, and name are required" ||
			err.Error()[:12] == "invalid role" {
			status = http.StatusBadRequest
		}
		response(c, status, false, nil, err)
		return
	}

	response(c, http.StatusCreated, true, user, nil)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	type LoginInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response(c, http.StatusBadRequest, false, nil, err)
		return
	}

	token, user, err := h.service.LoginUser(input.Email, input.Password)
	if err != nil {
		response(c, http.StatusUnauthorized, false, nil, err)
		return
	}

	response(c, http.StatusOK, true, gin.H{
		"token": token,
		"user":  user,
	}, nil)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response(c, http.StatusBadRequest, false, nil, errors.New("invalid user ID"))
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		response(c, http.StatusNotFound, false, nil, err)
		return
	}

	response(c, http.StatusOK, true, user, nil)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("userId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response(c, http.StatusBadRequest, false, nil, errors.New("invalid user ID"))
		return
	}

	var input services.EditUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response(c, http.StatusBadRequest, false, nil, err)
		return
	}

	user, err := h.service.UpdateUser(uint(id), input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error()[:12] == "invalid role" {
			status = http.StatusBadRequest
		}
		response(c, status, false, nil, err)
		return
	}

	response(c, http.StatusOK, true, user, nil)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	emailFilter := c.Query("email")
	roleFilter := c.Query("role")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response(c, http.StatusBadRequest, false, nil, errors.New("invalid page number"))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response(c, http.StatusBadRequest, false, nil, errors.New("invalid limit value"))
		return
	}

	input := services.UserListInput{
		Page:        page,
		Limit:       limit,
		EmailFilter: emailFilter,
		RoleFilter:  roleFilter,
	}

	result, err := h.service.GetUsers(input)
	if err != nil {
		response(c, http.StatusInternalServerError, false, nil, err)
		return
	}

	response(c, http.StatusOK, true, gin.H{
		"users": result.Users,
		"pagination": gin.H{
			"page":       result.Page,
			"limit":      result.Limit,
			"total":      result.Total,
			"totalPages": result.TotalPages,
		},
	}, nil)
}
