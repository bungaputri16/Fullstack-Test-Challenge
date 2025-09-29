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
	// Pastikan Content-Type JSON
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Decode body
	var req struct {
		ProductID int `json:"productId"`
		Qty       int `json:"qty"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validasi field
	if req.ProductID <= 0 {
		http.Error(w, "productId must be greater than 0", http.StatusUnprocessableEntity)
		return
	}
	if req.Qty <= 0 {
		http.Error(w, "qty must be greater than 0", http.StatusUnprocessableEntity)
		return
	}

	// Panggil service
	order, err := h.service.CreateOrder(req.ProductID, req.Qty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) GetOrdersByProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(mux.Vars(r)["productId"])
	if err != nil || productId <= 0 {
		http.Error(w, "invalid productId", http.StatusBadRequest)
		return
	}

	orders, err := h.service.GetOrdersByProduct(productId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
