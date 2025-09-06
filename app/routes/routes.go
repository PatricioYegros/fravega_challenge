package routes

import (
	"challenge_pyegros/app/handlers/orders"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpRoutes(client *mongo.Client, orderHandler *orders.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(router chi.Router) {
		router.Post("/orders", orderHandler.CreateOrder)
		router.Post("/orders/{orderId}/events", orderHandler.UpdateEventOrder)
		router.Get("/orders/{orderId}", orderHandler.GetOrderByID)
		router.Get("/orders/search", orderHandler.GetOrderByFilters)
	})

	return r
}
