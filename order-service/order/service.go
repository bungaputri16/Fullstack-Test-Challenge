package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-service/config"
	"strconv"
)

// -----------------------------
// Interfaces untuk dependency
// -----------------------------
type RepositoryInterface interface {
	Create(o *Order) error
	FindByProductID(productId int) ([]Order, error)
}

type RedisInterface interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, ttl int) error
	Del(key string) error
}

type RabbitInterface interface {
	Publish(queue string, msg interface{}) error
}

// -----------------------------
// Struct utama service
// -----------------------------
type Service struct {
	repo RepositoryInterface
	rdb  RedisInterface
	rmq  RabbitInterface
	cfg  config.Config
}

func NewService(repo RepositoryInterface, rdb RedisInterface, rmq RabbitInterface, cfg config.Config) *Service {
	return &Service{repo, rdb, rmq, cfg}
}

// -----------------------------
// Struct Product
// -----------------------------
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

// -----------------------------
// Helper: ambil product dari product-service
// -----------------------------
func (s *Service) getProduct(productId int) (*Product, error) {
	url := fmt.Sprintf("http://%s:%s/products/%d",
		s.cfg.ProductServiceHost, s.cfg.ProductServicePort, productId)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("product not found")
	}
	defer resp.Body.Close()

	var product Product
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, errors.New("failed to parse product data")
	}

	return &product, nil
}

// -----------------------------
// CreateOrder: membuat order
// -----------------------------
func (s *Service) CreateOrder(productId, qty int) (*Order, error) {
	// Ambil product info
	product, err := s.getProduct(productId)
	if err != nil {
		return nil, err
	}

	// Hitung total price
	totalPrice := product.Price * float64(qty)

	// Buat order object
	order := &Order{
		ProductID:  productId,
		TotalPrice: totalPrice,
		Status:     "CREATED",
	}

	// Insert ke database
	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	// Event data & cache key
	event := map[string]interface{}{"productId": productId, "qty": qty}
	cacheKey := "orders:product:" + strconv.Itoa(productId)

	// Publish event & hapus cache secara asynchronous
	go func() {
		_ = s.rmq.Publish("order.created", event)
		_ = s.rdb.Del(cacheKey)
	}()

	return order, nil
}

// -----------------------------
// GetOrdersByProduct: ambil orders berdasarkan productId dengan caching Redis
// -----------------------------
func (s *Service) GetOrdersByProduct(productId int) ([]Order, error) {
	cacheKey := "orders:product:" + strconv.Itoa(productId)

	var cached []Order
	if err := s.rdb.Get(cacheKey, &cached); err == nil {
		return cached, nil
	}

	orders, err := s.repo.FindByProductID(productId)
	if err != nil {
		return nil, err
	}

	if orders == nil {
		orders = []Order{}
	}

	_ = s.rdb.Set(cacheKey, orders, 60) // TTL 60 detik

	return orders, nil
}
