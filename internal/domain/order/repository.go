package order

import (
    "gorm.io/gorm"
    "errors"
)

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Create(order *Order) error {
    return r.db.Create(order).Error
}

func (r *Repository) GetByID(id uint) (*Order, error) {
    var order Order
    if err := r.db.Preload("Items.Product").First(&order, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("order not found")
        }
        return nil, err
    }
    return &order, nil
}

func (r *Repository) ListByUser(userID uint, page, limit int) ([]Order, int64, error) {
    var orders []Order
    var total int64

    query := r.db.Model(&Order{}).Where("user_id = ?", userID)
    query.Count(&total)

    offset := (page - 1) * limit
    err := query.Preload("Items.Product").Offset(offset).Limit(limit).Find(&orders).Error
    return orders, total, err
}

func (r *Repository) UpdateStatus(orderID uint, status OrderStatus) error {
    return r.db.Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
}