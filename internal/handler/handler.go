package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ravip18596/order-food-online/internal/model"
)

type Handler struct {
	// Add any dependencies here (e.g., service layer)
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	// Basic health check
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")
	// Product routes
	r.HandleFunc("/product", h.ListProducts).Methods("GET")
	r.HandleFunc("/product/{productId}", h.GetProduct).Methods("GET")

	// Order routes
	r.HandleFunc("/order", h.PlaceOrder).Methods("POST")
}

// ListProducts handles GET /product
func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement actual product listing logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

// GetProduct handles GET /product/{productId}
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// TODO: Implement actual product retrieval logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    productID,
		"name":  "Sample Product",
		"price": 9.99,
	})
}

// PlaceOrder handles POST /order
func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement order placement logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        "0000-0000-0000-0000",
		"total":     0.0,
		"discounts": 0.0,
		"items":     []interface{}{},
		"products":  []interface{}{},
	})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(model.HeartbeatResponse{
		Status: "OK",
		Code:   http.StatusOK,
	})
}
