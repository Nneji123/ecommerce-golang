package order

import (
	"time"

	"gorm.io/gorm"
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
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	Status        OrderStatus    `json:"status" gorm:"not null;default:'pending'"`
	Total         float64        `json:"total" gorm:"not null"`
	Items         []OrderItem    `json:"items"`
	ShippingAddr  string         `json:"shipping_address" gorm:"not null"`
	PaymentMethod string         `json:"payment_method" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
}

type CreateOrderRequest struct {
	Items         []OrderItemRequest `json:"items" validate:"required,dive"`
	ShippingAddr  string             `json:"shipping_address" validate:"required"`
	PaymentMethod string             `json:"payment_method" validate:"required"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required,oneof=pending confirmed shipped delivered cancelled"`
}
