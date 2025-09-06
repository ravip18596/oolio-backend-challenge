package model

type HeartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

// Image represents product images in different sizes
type Image struct {
	Thumbnail string `json:"thumbnail,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Tablet    string `json:"tablet,omitempty"`
	Desktop   string `json:"desktop,omitempty"`
}

// Product represents a product in the system
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Image    Image   `json:"image"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

// OrderRequest represents the request body for placing an order
type OrderRequest struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items"`
}

// Order represents an order in the system
type Order struct {
	ID        string      `json:"id"`
	Total     float64     `json:"total"`
	Discounts float64     `json:"discounts"`
	Items     []OrderItem `json:"items"`
	Products  []Product   `json:"products"`
}

// ApiResponse represents a standard API response
type ApiResponse struct {
	Code    int    `json:"code,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}
