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

	repoOrders := orderRepository.NewRepository(client)
	useCaseOrders := orderUseCase.NewUseCase(repoOrders)
	orderHandler := orderHandler.NewHandler(useCaseOrders)

	r := routes.SetUpRoutes(client, orderHandler)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
