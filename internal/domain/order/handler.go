package order

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/labstack/echo/v4"
// 	"github.com/nneji123/ecommerce-golang/internal/middleware"
// 	"github.com/nneji123/ecommerce-golang/internal/domain/product"
// )

// type Handler struct {
// 	repo          Repository
// 	productRepo   product.Repository
// 	emailService  email.Service
// 	validator     *validator.Validate
// 	logger        *zap.Logger
// }

// func NewHandler(
// 	repo Repository,
// 	productRepo product.Repository,
// 	emailService email.Service,
// 	validator *validator.Validate,
// 	logger *zap.Logger,
// ) *Handler {
// 	return &Handler{
// 		repo:          repo,
// 		productRepo:   productRepo,
// 		emailService:  emailService,
// 		validator:     validator,
// 		logger:        logger,
// 	}
// }

// // @Summary Create order
// // @Description Place a new order
// // @Tags orders
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param request body CreateOrderRequest true "Order details"
// // @Success 201 {object} Order
// // @Failure 400 {object} middleware.ErrorResponse
// // @Failure 401 {object} middleware.ErrorResponse
// // @Failure 422 {object} middleware.ErrorResponse
// // @Router /orders [post]
// func (h *Handler) CreateOrder(c echo.Context) error {
// 	var req CreateOrderRequest
// 	if err := c.Bind(&req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	if err := h.validator.Struct(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	// Get user ID from JWT claims
// 	claims := c.Get("user").(*middleware.Claims)
// 	userID := claims.UserID

// 	// Start transaction
// 	return h.repo.db.Transaction(func(tx *gorm.DB) error {
// 		// Get all products and validate stock
// 		var productIDs []uint
// 		for _, item := range req.Items {
// 			productIDs = append(productIDs, item.ProductID)
// 		}

// 		products, err := h.productRepo.FindByIDs(productIDs)
// 		if err != nil {
// 			h.logger.Error("Failed to fetch products", zap.Error(err))
// 			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process order")
// 		}

// 		// Create map for easy lookup
// 		productMap := make(map[uint]product.Product)
// 		for _, p := range products {
// 			productMap[p.ID] = p
// 		}

// 		// Calculate total and validate stock
// 		var total float64
// 		var orderItems []OrderItem
// 		for _, item := range req.Items {
// 			prod, exists := productMap[item.ProductID]
// 			if !exists {
// 				return echo.NewHTTPError(http.StatusUnprocessableEntity, "Product not found")
// 			}

// 			if prod.Stock < item.Quantity {
// 				return echo.NewHTTPError(http.StatusUnprocessableEntity, "Insufficient stock")
// 			}
