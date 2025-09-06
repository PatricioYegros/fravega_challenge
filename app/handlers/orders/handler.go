package orders

import ports "challenge_pyegros/app/ports/orders"

type Handler struct {
	u ports.OrdersUseCase
}

func NewHandler(u ports.OrdersUseCase) *Handler {
	return &Handler{u: u}
}
