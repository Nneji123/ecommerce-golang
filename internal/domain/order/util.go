package order

import "errors"

func IsValidOrderStatus(status OrderStatus) bool {
	validStatuses := map[OrderStatus]bool{
		StatusPending:   true,
		StatusConfirmed: true,
		StatusShipped:   true,
		StatusDelivered: true,
		StatusCancelled: true,
	}
	return validStatuses[status]
}

func CalculateOrderTotal(items []OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func ValidateOrder(order *Order) error {
	if len(order.Items) == 0 {
		return errors.New("order must contain at least one item")
	}

	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return errors.New("item quantity must be greater than 0")
		}
		if item.Price <= 0 {
			return errors.New("item price must be greater than 0")
		}
	}

	return nil
}
