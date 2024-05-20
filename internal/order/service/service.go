package service

import (
	"WB/internal/order/storage/database/dto"
	"context"
	"go.uber.org/zap"

	"WB/internal/order/model"
	"WB/internal/order/storage"
)

type OrderService interface {
	AddOrder(ctx context.Context, order model.Order) (*model.Order, error)
	GetUserOrder(ctx context.Context, orderID string) (*dto.OrderFromCache, error)
}

type OrderServiceApp struct {
	storage storage.OrderStorage
	logger  *zap.Logger
}

func New(storage storage.OrderStorage) *OrderServiceApp {
	return &OrderServiceApp{
		storage: storage,
	}
}
