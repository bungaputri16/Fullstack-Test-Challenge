package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductID int `json:"productId"`
		Qty       int `json:"qty"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	order, err := h.service.CreateOrder(req.ProductID, req.Qty)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) GetOrdersByProduct(w http.ResponseWriter, r *http.Request) {
	productId, _ := strconv.Atoi(mux.Vars(r)["productId"])

	orders, err := h.service.GetOrdersByProduct(productId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(orders)
}
