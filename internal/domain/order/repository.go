package order

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(order *Order) error
	FindByID(id uint) (*Order, error)
	ListByUser(userID uint, page, limit int) ([]Order, int, error)
	UpdateStatus(id uint, status OrderStatus) error
	FindByUserAndID(userID, orderID uint) (*Order, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(order *Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Create order items
		for _, item := range order.Items {
			item.OrderID = order.ID
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *repository) FindByID(id uint) (*Order, error) {
	var order Order
	err := r.db.Preload("Items").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *repository) ListByUser(userID uint, page, limit int) ([]Order, int, error) {
	var orders []Order
	var total int64

	query := r.db.Model(&Order{}).Where("user_id = ?", userID)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := query.Offset((page - 1) * limit).
		Limit(limit).
		Preload("Items").
		Order("created_at DESC").
		Find(&orders).Error

	return orders, int(total), err
}

func (r *repository) UpdateStatus(id uint, status OrderStatus) error {
	return r.db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *repository) FindByUserAndID(userID, orderID uint) (*Order, error) {
	var order Order
	err := r.db.Preload("Items").Where("user_id = ? AND id = ?", userID, orderID).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}