package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ravip18596/order-food-online/internal/model"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

type ProductDB struct {
	ID             string    `db:"id"`
	Name           string    `db:"name"`
	Price          float64   `db:"price"`
	Category       string    `db:"category"`
	ImageThumbnail string    `db:"image_thumbnail"`
	ImageMobile    string    `db:"image_mobile"`
	ImageTablet    string    `db:"image_tablet"`
	ImageDesktop   string    `db:"image_desktop"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (r *ProductRepository) Create(product *model.Product) (*model.Product, error) {
	if product == nil {
		return nil, errors.New("product cannot be nil")
	}

	// Generate new UUID if not provided
	if product.ID == "" {
		product.ID = uuid.New().String()
	}

	query := `
		INSERT INTO products (
			id, name, price, category, 
			image_thumbnail, image_mobile, image_tablet, image_desktop
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Category,
		product.Image.Thumbnail,
		product.Image.Mobile,
		product.Image.Tablet,
		product.Image.Desktop,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) GetByID(id string) (*model.Product, error) {
	var dbProduct ProductDB
	query := `SELECT * FROM products WHERE id = ?`

	err := r.db.Get(&dbProduct, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.Product{
		ID:       dbProduct.ID,
		Name:     dbProduct.Name,
		Price:    dbProduct.Price,
		Category: dbProduct.Category,
		Image: model.Image{
			Thumbnail: dbProduct.ImageThumbnail,
			Mobile:    dbProduct.ImageMobile,
			Tablet:    dbProduct.ImageTablet,
			Desktop:   dbProduct.ImageDesktop,
		},
	}, nil
}

func (r *ProductRepository) GetAll() ([]*model.Product, error) {
	query := `SELECT * FROM products`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var dbProduct ProductDB

		err = rows.Scan(
			&dbProduct.ID,
			&dbProduct.Name,
			&dbProduct.Price,
			&dbProduct.Category,
			&dbProduct.ImageThumbnail,
			&dbProduct.ImageMobile,
			&dbProduct.ImageTablet,
			&dbProduct.ImageDesktop,
			&dbProduct.CreatedAt,
			&dbProduct.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, &model.Product{
			ID:       dbProduct.ID,
			Name:     dbProduct.Name,
			Price:    dbProduct.Price,
			Category: dbProduct.Category,
			Image: model.Image{
				Thumbnail: dbProduct.ImageThumbnail,
				Mobile:    dbProduct.ImageMobile,
				Tablet:    dbProduct.ImageTablet,
				Desktop:   dbProduct.ImageDesktop,
			},
		})
	}

	return products, nil
}
