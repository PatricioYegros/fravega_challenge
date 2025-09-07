package orders

import (
	"challenge_pyegros/app/models"
	ports "challenge_pyegros/app/ports/orders"

	"github.com/redis/go-redis/v9"
)

type UseCase struct {
	r     ports.OrdersUseCase
	redis *redis.Client
}

func NewUseCase(r ports.OrdersUseCase, redis *redis.Client) *UseCase {
	return &UseCase{
		r:     r,
		redis: redis,
	}
}

func (u *UseCase) CreateOrder(order models.Order) (*models.ResponseCreate, error) {
	return u.r.CreateOrder(order)
}

func (u *UseCase) UpdateEventOrder(orderID int64, event models.Event) (*models.ResponseUpdate, error) {
	return u.r.UpdateEventOrder(orderID, event)
}

func (u *UseCase) GetOrderByID(orderID int64) (*models.ResponseGet, error) {
	return u.r.GetOrderByID(orderID)
}

func (u *UseCase) GetOrderByFilters(filters models.Filters) ([]models.Order, error) {
	return u.r.GetOrderByFilters(filters)
}
