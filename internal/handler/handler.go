package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ravip18596/order-food-online/internal/model"
	"github.com/ravip18596/order-food-online/internal/repository"
)

type Handler struct {
	productRepo *repository.ProductRepository
	orderRepo   *repository.OrderRepository
	set         map[string][]int
}

func NewHandler(productRepo *repository.ProductRepository, orderRepo *repository.OrderRepository,
	set map[string][]int) *Handler {
	return &Handler{
		productRepo: productRepo,
		orderRepo:   orderRepo,
		set:         set,
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	// Basic health check
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")

	// Product routes
	r.HandleFunc("/product", h.ListProducts).Methods("GET")
	r.HandleFunc("/product", h.CreateProduct).Methods("POST")
	r.HandleFunc("/product/{productId}", h.GetProduct).Methods("GET")

	// Order routes
	r.HandleFunc("/order", h.PlaceOrder).Methods("POST")
}

// CreateProduct handles POST /product
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if product.Name == "" || product.Price <= 0 || product.Category == "" {
		http.Error(w, "Name, price, and category are required fields", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.productRepo.Create(&product)
	if err != nil {
		http.Error(w, "Error creating product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productRepo.GetAll()
	if err != nil {
		http.Error(w, "Error getting all products: "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GET /product/{productId}
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]

	product, err := h.productRepo.GetByID(productID)
	if err != nil {
		http.Error(w, "Error fetching product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// PlaceOrder handles POST /order
func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq model.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate order items
	if len(orderReq.Items) == 0 {
		http.Error(w, "At least one order item is required", http.StatusBadRequest)
		return
	}

	// Prepare order with items
	order := &model.Order{
		ID:        uuid.New().String(),
		Items:     make([]model.OrderItem, 0, len(orderReq.Items)),
		Products:  make([]model.Product, 0, len(orderReq.Items)),
		Discounts: 0,
	}

	// Calculate total and validate products
	total := 0.0
	for _, item := range orderReq.Items {
		if item.Quantity <= 0 {
			http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
			return
		}

		product, err := h.productRepo.GetByID(item.ProductID)
		if err != nil {
			http.Error(w, "Error fetching product: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if product == nil {
			http.Error(w, "Product not found: "+item.ProductID, http.StatusNotFound)
			return
		}

		// Add to order items
		order.Items = append(order.Items, model.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
		})

		// Add product details to order
		order.Products = append(order.Products, *product)

		// Calculate subtotal
		subtotal := product.Price * float64(item.Quantity)
		total += subtotal
	}

	// Apply coupon code if provided
	if orderReq.CouponCode != "" {
		if len(orderReq.CouponCode) >= 8 && len(orderReq.CouponCode) <= 10 {
			fileNos, ok := h.set[orderReq.CouponCode]
			if ok && len(fileNos) == 2 {
				fmt.Println("Coupon code " + orderReq.CouponCode + " is valid")
				order.Discounts = total * 0.1
			} else {
				fmt.Println("Coupon code " + orderReq.CouponCode + " is invalid")
			}
		} else {
			fmt.Println("Coupon code " + orderReq.CouponCode + " is invalid")
		}
	}

	order.Total = total - order.Discounts

	// Create the order in database
	createdOrder, err := h.orderRepo.Create(order)
	if err != nil {
		http.Error(w, "Error creating order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created order with 201 status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrder)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(model.HeartbeatResponse{
		Status: "OK",
		Code:   http.StatusOK,
	})
}
