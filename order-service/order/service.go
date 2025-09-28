package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-service/config"
	"order-service/rabbitmq"
	"order-service/redis"
	"strconv"
)

type Service struct {
	repo *Repository
	rdb  *redis.Client
	rmq  *rabbitmq.Client
	cfg  config.Config
}

func NewService(repo *Repository, rdb *redis.Client, rmq *rabbitmq.Client, cfg config.Config) *Service {
	return &Service{repo, rdb, rmq, cfg}
}

func (s *Service) CreateOrder(productId, qty int) (*Order, error) {
	url := fmt.Sprintf("http://%s:%s/products/%d",
		s.cfg.ProductServiceHost, s.cfg.ProductServicePort, productId)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("product not found")
	}
	defer resp.Body.Close()

	var product struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Qty   int     `json:"qty"`
	}
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &product)

	totalPrice := product.Price * float64(qty)

	order := &Order{
		ProductID:  productId,
		TotalPrice: totalPrice,
		Status:     "CREATED",
	}

	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	event := map[string]interface{}{"productId": productId, "qty": qty}
	s.rmq.Publish("order.created", event)

	key := "orders:product:" + strconv.Itoa(productId)
	_ = s.rdb.Del(key)

	return order, nil
}

func (s *Service) GetOrdersByProduct(productId int) ([]Order, error) {
	// key := "orders:product:" + strconv.Itoa(productId)

	// var cached []Order
	// if err := s.rdb.Get(key, &cached); err == nil {
	// 	// Pastikan slice tidak nil
	// 	if cached == nil {
	// 		cached = []Order{}
	// 	}
	// 	return cached, nil
	// }

	orders, err := s.repo.FindByProductID(productId)
	if err != nil {
		return nil, err
	}
	
	if orders == nil {
		orders = []Order{}
		}
		
	// _ = s.rdb.Set(key, orders, 60)
	
	return orders, nil
}
