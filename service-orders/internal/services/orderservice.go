package services

import (
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

func (s *OrderService) GetOrders() ([]*OrderResponse, error) {

	var orders []models.Order
	if err := s.orderRepo.GetOrders(&orders); err != nil {
		return nil, err
	}

	var orderResponses []*OrderResponse
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
		orderResponses = append(orderResponses, orderResponse)
	}

	return orderResponses, nil
}
