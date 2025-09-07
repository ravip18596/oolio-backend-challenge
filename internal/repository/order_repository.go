package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ravip18596/order-food-online/internal/model"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type OrderDB struct {
	ID         string    `db:"id"`
	Total      float64   `db:"total"`
	Discounts  float64   `db:"discounts"`
	CouponCode *string   `db:"coupon_code"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type OrderItemDB struct {
	ID           string    `db:"id"`
	OrderID      string    `db:"order_id"`
	ProductID    string    `db:"product_id"`
	Quantity     int       `db:"quantity"`
	PricePerUnit float64   `db:"price_per_unit"`
	CreatedAt    time.Time `db:"created_at"`
}

func (r *OrderRepository) Create(order *model.Order) (*model.Order, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Generate order ID if not provided
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	// Insert order
	query := `
		INSERT INTO orders (id, total, discounts)
		VALUES (?, ?, ?)
	`

	_, err = tx.Exec(query, order.ID, order.Total, order.Discounts)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// Insert order items
	for _, item := range order.Items {
		itemID := uuid.New().String()
		query = `
			INSERT INTO order_items (id, order_id, product_id, quantity, price_per_unit)
			VALUES (?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(query, itemID, order.ID, item.ProductID, item.Quantity, 0) // TODO: Get actual price
		if err != nil {
			return nil, fmt.Errorf("error creating order item: %w", err)
		}
	}

	return order, nil
}

func (r *OrderRepository) GetByID(id string) (*model.Order, error) {
	// Get order
	var orderDB OrderDB
	query := `SELECT * FROM orders WHERE id = ?`
	err := r.db.Get(&orderDB, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching order: %w", err)
	}

	// Get order items
	var itemsDB []OrderItemDB
	query = `SELECT * FROM order_items WHERE order_id = ?`
	err = r.db.Select(&itemsDB, query, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching order items: %w", err)
	}

	// Convert to model
	order := &model.Order{
		ID:        orderDB.ID,
		Total:     orderDB.Total,
		Discounts: orderDB.Discounts,
		Items:     make([]model.OrderItem, len(itemsDB)),
	}

	for i, item := range itemsDB {
		order.Items[i] = model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	return order, nil
}

func (r *OrderRepository) List() ([]model.Order, error) {
	// Get all orders
	var ordersDB []OrderDB
	query := `SELECT * FROM orders ORDER BY created_at DESC`
	err := r.db.Select(&ordersDB, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching orders: %w", err)
	}

	if len(ordersDB) == 0 {
		return []model.Order{}, nil
	}

	// Get order IDs for batch fetching items
	orderIDs := make([]string, len(ordersDB))
	for i, order := range ordersDB {
		orderIDs[i] = order.ID
	}

	// Get all order items for these orders
	var itemsDB []OrderItemDB
	query, args, err := sqlx.In(`
		SELECT * FROM order_items 
		WHERE order_id IN (?) 
		ORDER BY created_at
	`, orderIDs)

	if err != nil {
		return nil, fmt.Errorf("error building order items query: %w", err)
	}

	query = r.db.Rebind(query)
	err = r.db.Select(&itemsDB, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error fetching order items: %w", err)
	}

	// Create a map of order ID to its items
	itemsByOrderID := make(map[string][]model.OrderItem)
	for _, item := range itemsDB {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// Combine orders with their items
	orders := make([]model.Order, len(ordersDB))
	for i, orderDB := range ordersDB {
		order := model.Order{
			ID:        orderDB.ID,
			Total:     orderDB.Total,
			Discounts: orderDB.Discounts,
			Items:     itemsByOrderID[orderDB.ID],
		}
		orders[i] = order
	}

	return orders, nil
}
