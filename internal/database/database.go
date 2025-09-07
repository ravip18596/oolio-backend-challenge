package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func InitDB() error {
	// Create data directory if it doesn't exist
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return fmt.Errorf("error creating data directory: %w", err)
	}

	// Open SQLite database
	dbPath := filepath.Join("data", "orders.db")
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Enable foreign keys
	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return fmt.Errorf("error enabling foreign keys: %w", err)
	}

	DB = db
	log.Println("Connected to SQLite database")

	return createTables()
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		category TEXT NOT NULL,
		image_thumbnail TEXT,
		image_mobile TEXT,
		image_tablet TEXT,
		image_desktop TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		total REAL NOT NULL,
		discounts REAL NOT NULL DEFAULT 0,
		coupon_code TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS order_items (
		id TEXT PRIMARY KEY,
		order_id TEXT NOT NULL,
		product_id TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		price_per_unit REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id)
	);

	CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
	CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);

	CREATE TRIGGER IF NOT EXISTS update_products_updated_at
	AFTER UPDATE ON products
	BEGIN
		UPDATE products SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

	CREATE TRIGGER IF NOT EXISTS update_orders_updated_at
	AFTER UPDATE ON orders
	BEGIN
		UPDATE orders SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
