package services

import (
	"errors"

	"github.com/SpiritFoxo/control-system-microservices/service-orders/internal/config"
	"github.com/SpiritFoxo/control-system-microservices/service-orders/internal/models"
	"github.com/SpiritFoxo/control-system-microservices/service-orders/internal/repositories"
)

type OrderService struct {
	orderRepo *repositories.OrderRepository
	cfg       *config.Config
}

func NewOrderService(orderRepo *repositories.OrderRepository, cfg *config.Config) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		cfg:       cfg,
	}
}

type OrderResponse struct {
	ID         uint                `json:"id"`
	UserID     uint                `json:"user_id"`
	Status     string              `json:"status"`
	Cost       int                 `json:"cost"`
	OrderItems []OrderItemResponse `json:"order_items"`
}

type OrderItemResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type CreateOrderInput struct {
	UserID     uint             `json:"user_id" binding:"required"`
	Status     string           `json:"status" binding:"required,oneof=created in_progress completed cancelled"`
	OrderItems []OrderItemInput `json:"order_items" binding:"required,min=1"`
	Cost       int              `json:"cost" binding:"required,min=0"`
}

type OrderListInput struct {
	Page   int    `form:"page" json:"page"`
	Limit  int    `form:"limit" json:"limit"`
	UserID uint   `form:"userId" json:"userId"`
	Status string `form:"status" json:"status"`
}

type OrderListResponse struct {
	Orders     []*OrderResponse `json:"orders"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
}

type OrderItemInput struct {
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type UpdateOrderInput struct {
	Status string `json:"status" binding:"required,oneof=created in_progress completed cancelled"`
}

func (s *OrderService) GetOrderByID(id uint) (*OrderResponse, error) {
	order, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	orderResponse := &OrderResponse{
		ID:         order.ID,
		UserID:     order.UserId,
		Status:     order.Status,
		Cost:       order.Cost,
		OrderItems: make([]OrderItemResponse, len(order.Items)),
	}

	for i, item := range order.Items {
		orderResponse.OrderItems[i] = OrderItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
		}
	}

	return orderResponse, nil
}

func (s *OrderService) CreateOrder(input *CreateOrderInput) (*OrderResponse, error) {
	var orderItems []models.OrderItem

	order := &models.Order{
		UserId: input.UserID,
		Status: input.Status,
		Cost:   input.Cost,
		Items:  orderItems,
	}

	if err := s.orderRepo.CreateOrder(order); err != nil {
		return nil, err
	}

	orderResponse := &OrderResponse{
		ID:         order.ID,
		UserID:     order.UserId,
		Status:     order.Status,
		Cost:       order.Cost,
		OrderItems: make([]OrderItemResponse, len(order.Items)),
	}
	for i, item := range order.Items {
		orderResponse.OrderItems[i] = OrderItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
		}
	}

	return orderResponse, nil
}

func (s *OrderService) UpdateOrder(id uint, input *UpdateOrderInput) (*OrderResponse, error) {

	order, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	order.Status = input.Status
	if err := s.orderRepo.UpdateOrder(order); err != nil {
		return nil, err
	}

	updatedOrder, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	orderResponse := &OrderResponse{
		ID:         updatedOrder.ID,
		UserID:     updatedOrder.UserId,
		Status:     updatedOrder.Status,
		Cost:       updatedOrder.Cost,
		OrderItems: make([]OrderItemResponse, len(updatedOrder.Items)),
	}
	for i, item := range updatedOrder.Items {
		orderResponse.OrderItems[i] = OrderItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
		}
	}

	return orderResponse, nil
}

func (s *OrderService) DeleteOrder(id uint) error {

	order, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return err
	}

	if err := s.orderRepo.DeleteOrder(order); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetOrders(input OrderListInput) (*OrderListResponse, error) {

	if input.Page < 1 {
		return nil, errors.New("invalid page number")
	}
	if input.Limit < 1 {
		return nil, errors.New("invalid limit value")
	}

	orders, total, err := s.orderRepo.GetOrders(input.Page, input.Limit, input.UserID, input.Status)
	if err != nil {
		return nil, err
	}

	response := make([]*OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponse := &OrderResponse{
			ID:         order.ID,
			UserID:     order.UserId,
			Status:     order.Status,
			Cost:       order.Cost,
			OrderItems: make([]OrderItemResponse, len(order.Items)),
		}
		for i, item := range order.Items {
			orderResponse.OrderItems[i] = OrderItemResponse{
				ID:       item.ID,
				Name:     item.Name,
				Quantity: item.Quantity,
			}
		}
		response = append(response, orderResponse)
	}

	totalPages := int((total + int64(input.Limit) - 1) / int64(input.Limit))

	return &OrderListResponse{
		Orders:     response,
		Total:      total,
		Page:       input.Page,
		Limit:      input.Limit,
		TotalPages: totalPages,
	}, nil
}
