package order

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"net"
	"testing"

	"order-service/config"
)

// Mock Repository
type mockRepo struct {
	createCalled bool
	fail         bool
}

func (m *mockRepo) Create(o *Order) error {
	m.createCalled = true
	if m.fail {
		return errors.New("db error")
	}
	return nil
}
func (m *mockRepo) FindByProductID(productId int) ([]Order, error) {
	return nil, nil
}

// Mock Redis
type mockRedis struct {
	delCalled bool
}

func (r *mockRedis) Get(key string, dest interface{}) error {
	return errors.New("not found") // default: cache miss
}
func (r *mockRedis) Set(key string, value interface{}, ttl int) error {
	return nil
}
func (r *mockRedis) Del(key string) error {
	r.delCalled = true
	return nil
}


// Mock RabbitMQ
type mockRabbit struct {
	published bool
}
func (r *mockRabbit) Publish(queue string, msg interface{}) error {
	r.published = true
	return nil
}

func TestCreateOrder_Success(t *testing.T) {
	// fake product-service server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":1,"name":"Phone","price":1000,"qty":10}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// parse host dan port dari server.URL
	u, _ := url.Parse(server.URL)
	host, port, _ := net.SplitHostPort(u.Host)

	cfg := config.Config{
		ProductServiceHost: host,
		ProductServicePort: port,
	}

	repo := &mockRepo{}
	rdb := &mockRedis{}
	rmq := &mockRabbit{}

	service := NewService(repo, rdb, rmq, cfg)

	order, err := service.CreateOrder(1, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.TotalPrice != 2000 {
		t.Errorf("expected total price 2000, got %v", order.TotalPrice)
	}
	if !repo.createCalled {
		t.Error("expected repo.Create to be called")
	}
	if !rdb.delCalled {
		t.Error("expected redis.Del to be called")
	}
	if !rmq.published {
		t.Error("expected rabbitmq.Publish to be called")
	}
	t.Log(" should create order successfully and trigger repo, redis, and rabbitmq")

}

func TestCreateOrder_ProductNotFound(t *testing.T) {
	// fake product-service server balikin 404
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	u, _ := url.Parse(server.URL)
	host, port, _ := net.SplitHostPort(u.Host)

	cfg := config.Config{
		ProductServiceHost: host,
		ProductServicePort: port,
	}

	repo := &mockRepo{}
	rdb := &mockRedis{}
	rmq := &mockRabbit{}

	service := NewService(repo, rdb, rmq, cfg)

	_, err := service.CreateOrder(99, 1)
	if err == nil || err.Error() != "product not found" {
		t.Fatalf("expected 'product not found' error, got %v", err)
	}

	t.Log("✓ should return error when product-service responds 404")
}
