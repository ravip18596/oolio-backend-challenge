package model

type HeartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type Image struct {
	Thumbnail string `json:"thumbnail,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Tablet    string `json:"tablet,omitempty"`
	Desktop   string `json:"desktop,omitempty"`
}

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Image    Image   `json:"image"`
}

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderRequest struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items"`
}

type Order struct {
	ID        string      `json:"id"`
	Total     float64     `json:"total"`
	Discounts float64     `json:"discounts"`
	Items     []OrderItem `json:"items"`
	Products  []Product   `json:"products"`
}

type ApiResponse struct {
	Code    int    `json:"code,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}
