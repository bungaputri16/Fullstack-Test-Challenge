package main

import (
	"log"
	"net/http"

	"order-service/config"
	"order-service/order"
	"order-service/rabbitmq"
	"order-service/redis"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	db := config.InitDB(cfg)
	rdb := redis.NewRedis(cfg)
	rmq := rabbitmq.NewRabbit(cfg)

	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo, rdb, rmq, cfg)
	orderHandler := order.NewHandler(orderService)

	r := mux.NewRouter()
	r.Use(order.RequestID)        
	r.Use(order.ValidateJSONContent) 

	r.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	r.HandleFunc("/orders/product/{productId}", orderHandler.GetOrdersByProduct).Methods("GET")

	log.Println("Order-service running on port", cfg.ServicePort)
	http.ListenAndServe(":"+cfg.ServicePort, r)

	
}
