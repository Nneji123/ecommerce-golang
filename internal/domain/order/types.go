package order

import (
    "time"
    "gorm.io/gorm"
    "github.com/nneji123/ecommerce-golang/internal/domain/product"
)

type OrderStatus string

const (
    StatusPending   OrderStatus = "pending"
    StatusConfirmed OrderStatus = "confirmed"
    StatusShipped   OrderStatus = "shipped"
    StatusDelivered OrderStatus = "delivered"
    StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
    ID            uint           `gorm:"primaryKey" json:"id"`
    UserID        uint           `gorm:"not null" json:"user_id"`
    Status        OrderStatus    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
    TotalAmount   float64        `gorm:"not null" json:"total_amount"`
    Items         []OrderItem    `json:"items"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
    ID        uint    `gorm:"primaryKey" json:"id"`
    OrderID   uint    `gorm:"not null" json:"order_id"`
    ProductID uint    `gorm:"not null" json:"product_id"`
    Product   product.Product `json:"product"`
    Quantity  int     `gorm:"not null" json:"quantity"`
    Price     float64 `gorm:"not null" json:"price"`
}
