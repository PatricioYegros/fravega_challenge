package main

import (
	"challenge_pyegros/app/database"
	orderHandler "challenge_pyegros/app/handlers/orders"
	orderRepository "challenge_pyegros/app/repositories/orders"
	"challenge_pyegros/app/routes"
	orderUseCase "challenge_pyegros/app/usecases/orders"
	"context"
	"log"
	"net/http"
)

// @title           Orders API
// @version         1.0
// @description     API for create and manage orders.

// @contact.name   Patricio Yegros
// @contact.url    github.com/PatricioYegros/fravega_challenge
// @contact.email  patricioyegros@hotmail.com

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	client, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	rdb := database.ConnectRedis()

	repoOrders := orderRepository.NewRepository(client, rdb)
	useCaseOrders := orderUseCase.NewUseCase(repoOrders, rdb)
	orderHandler := orderHandler.NewHandler(useCaseOrders)

	r := routes.SetUpRoutes(client, orderHandler)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
