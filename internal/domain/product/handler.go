package product

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/labstack/echo/v4"
// 	"github.com/nneji123/ecommerce-golang/internal/middleware"
// )

// type Handler struct {
// 	repo      Repository
// 	validator *validator.Validate
// 	logger    *zap.Logger
// }

// func NewHandler(repo Repository, validator *validator.Validate, logger *zap.Logger) *Handler {
// 	return &Handler{
// 		repo:      repo,
// 		validator: validator,
// 		logger:    logger,
// 	}
// }

// // @Summary Create product
// // @Description Create a new product (admin only)
// // @Tags products
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param request body ProductRequest true "Product details"
// // @Success 201 {object} Product
// // @Failure 400 {object} middleware.ErrorResponse
// // @Failure 401,403 {object} middleware.ErrorResponse
// // @Router /admin/products [post]
// func (h *Handler) CreateProduct(c echo.Context) error {
// 	var req ProductRequest
// 	if err := c.Bind(&req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	if err := h.validator.Struct(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	product := &Product{
// 		Name:        req.Name,
// 		Description: req.Description,
// 		Price:       req.Price,
// 		Stock:       req.Stock,
// 	}

// 	if err := h.repo.Create(product, req.Categories); err != nil {
// 		h.logger.Error("Failed to create product", zap.Error(err))
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create product")
// 	}

// 	return c.JSON(http.StatusCreated, product)
// }

// // @Summary Update product
// // @Description Update an existing product (admin only)
// // @Tags products
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Product ID"
// // @Param request body ProductRequest true "Product details"
// // @Success 200 {object} Product
// // @Failure 400,404 {object} middleware.ErrorResponse
// // @Failure 401,403 {object} middleware.ErrorResponse
// // @Router /admin/products/{id} [put]
// func (h *Handler) UpdateProduct(c echo.Context) error {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
// 	}

// 	var req ProductRequest
// 	if err := c.Bind(&req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	if err := h.validator.Struct(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	product, err := h.repo.FindByID(uint(id))
// 	if err != nil {
// 		h.logger.Error("Failed to find product", zap.Error(err))
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to find product")
// 	}
// 	if product == nil {
// 		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
// 	}

// 	product.Name = req.Name
// 	product.Description = req.Description
// 	product.Price = req.Price
// 	product.Stock = req.Stock

// 	if err := h.repo.Update(product, req.Categories); err != nil {
// 		h.logger.Error("Failed to update product", zap.Error(err))
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update product")
// 	}

// 	return c.JSON(http.StatusOK, product)
// }
