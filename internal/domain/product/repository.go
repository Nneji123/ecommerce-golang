package product

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(product *Product) error {
	return r.db.Create(product).Error
}

func (r *Repository) GetByID(id uint) (*Product, error) {
	var product Product
	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *Repository) Update(product *Product) error {
	return r.db.Save(product).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}

type ListProductsQuery struct {
	Page     int     `query:"page"`
	Limit    int     `query:"limit"`
	MinPrice float64 `query:"min_price"`
	MaxPrice float64 `query:"max_price"`
	Search   string  `query:"search"`
	SortBy   string  `query:"sort_by"`
	SortDir  string  `query:"sort_dir"`
}

func (r *Repository) List(query *ListProductsQuery) ([]Product, int64, error) {
	var products []Product
	var total int64

	db := r.db.Model(&Product{})

	// Apply filters
	if query.MinPrice > 0 {
		db = db.Where("price >= ?", query.MinPrice)
	}
	if query.MaxPrice > 0 {
		db = db.Where("price <= ?", query.MaxPrice)
	}
	if query.Search != "" {
		search := "%" + query.Search + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", search, search)
	}

	// Count total before pagination
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if query.SortBy != "" {
		direction := "ASC"
		if strings.ToUpper(query.SortDir) == "DESC" {
			direction = "DESC"
		}

		validColumns := map[string]bool{
			"name": true, "price": true, "created_at": true,
		}
		if validColumns[query.SortBy] {
			db = db.Order(query.SortBy + " " + direction)
		}
	} else {
		db = db.Order("created_at DESC")
	}

	page := query.Page
	if page < 1 {
		page = 1
	}

	limit := query.Limit
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
