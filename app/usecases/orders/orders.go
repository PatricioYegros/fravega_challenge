package orders

import (
	"challenge_pyegros/app/models"
	ports "challenge_pyegros/app/ports/orders"
)

type UseCase struct {
	r ports.OrdersUseCase
}

func NewUseCase(r ports.OrdersUseCase) *UseCase {
	return &UseCase{r: r}
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
