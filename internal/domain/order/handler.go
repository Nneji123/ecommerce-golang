package order

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/common/models"
	"go.uber.org/zap"
)

type Handler struct {
    repo    *Repository
    logger  *zap.Logger
}

func NewHandler(repo *Repository, logger *zap.Logger) *Handler {
    return &Handler{
        repo:    repo,
        logger:  logger,
    }
}

//	@Summary		Create order
//	@Description	Place a new order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			order	body		Order	true	"Order object"
//	@Success		201		{object}	Order
//	@Failure		400		{object}	middleware.ErrorResponse
//	@Router			/orders [post]
func (h *Handler) Create(c echo.Context) error {
    var order Order
    if err := c.Bind(&order); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    // Get user ID from claims
    claims := c.Get("userClaims").(*models.Claims)
    order.UserID = claims.UserID

    if err := h.repo.Create(&order); err != nil {
        h.logger.Error("Failed to create order", zap.Error(err))
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create order")
    }

    return c.JSON(http.StatusCreated, order)
}

//	@Summary		List user orders
//	@Description	List all orders for the authenticated user
//	@Tags			orders
//	@Produce		json
//	@Success		200	{array}	Order
//	@Router			/orders [get]
func (h *Handler) ListUserOrders(c echo.Context) error {
    claims := c.Get("userClaims").(*models.Claims)
    page, _ := strconv.Atoi(c.QueryParam("page"))
    if page < 1 {
        page = 1
    }
    limit := 10

    orders, total, err := h.repo.ListByUser(claims.UserID, page, limit)
    if err != nil {
        h.logger.Error("Failed to list orders", zap.Error(err))
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list orders")
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "orders": orders,
        "total":  total,
        "page":   page,
        "limit":  limit,
    })
}

//	@Summary		Cancel order
//	@Description	Cancel an order (only if pending)
//	@Tags			orders
//	@Param			id	path		int	true	"Order ID"
//	@Success		200	{object}	Order
//	@Failure		400	{object}	middleware.ErrorResponse
//	@Router			/orders/{id}/cancel [post]
func (h *Handler) CancelOrder(c echo.Context) error {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
    }

    order, err := h.repo.GetByID(uint(id))
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Order not found")
    }

    claims := c.Get("userClaims").(*models.Claims)
    if order.UserID != claims.UserID {
        return echo.NewHTTPError(http.StatusForbidden, "Not authorized to cancel this order")
    }

    if order.Status != StatusPending {
        return echo.NewHTTPError(http.StatusBadRequest, "Only pending orders can be cancelled")
    }

    if err := h.repo.UpdateStatus(order.ID, StatusCancelled); err != nil {
        h.logger.Error("Failed to cancel order", zap.Error(err))
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to cancel order")
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "Order cancelled successfully"})
}

//	@Summary		Update order status
//	@Description	Update order status (Admin only)
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Order ID"
//	@Param			status	body		string	true	"New status"
//	@Success		200		{object}	Order
//	@Failure		400		{object}	middleware.ErrorResponse
//	@Router			/orders/{id}/status [put]
func (h *Handler) UpdateStatus(c echo.Context) error {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
    }

    var statusUpdate struct {
        Status OrderStatus `json:"status"`
    }

    if err := c.Bind(&statusUpdate); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    // Validate status
    if !IsValidOrderStatus(statusUpdate.Status) {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid order status")
    }

    if err := h.repo.UpdateStatus(uint(id), statusUpdate.Status); err != nil {
        h.logger.Error("Failed to update order status", zap.Error(err))
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update order status")
    }

    return c.JSON(http.StatusOK, map[string]string{
        "message": "Order status updated successfully",
        "status":  string(statusUpdate.Status),
    })
}