package ports

import (
	"challenge_pyegros/app/models"
)

//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=./$GOFILE -destination=./mocks/$GOFILE -package mocks

type OrdersRepository interface {
	CreateOrder(order models.Order) (*models.ResponseCreate, error)
	UpdateEventOrder(orderID int64, event models.Event) (*models.ResponseUpdate, error)
	GetOrderByID(orderID int64) (*models.ResponseGet, error)
	GetOrderByFilters(filters models.Filters) ([]models.Order, error)
}
