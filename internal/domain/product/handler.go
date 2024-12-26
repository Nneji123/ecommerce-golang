package product

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	repo   *Repository
	logger *zap.Logger
}

func NewHandler(repo *Repository, logger *zap.Logger) *Handler {
	return &Handler{
		repo:   repo,
		logger: logger,
	}
}

// @Summary		Create product
// @Description	Create a new product (Admin only)
// @Tags			products
// @Accept			json
// @Produce		json
// @Param			product	body		Product	true	"Product object"
// @Success		201		{object}	Product
// @Failure		400		{object}	middleware.ErrorResponse
// @Router			/products [post]
func (h *Handler) Create(c echo.Context) error {
	var product Product
	if err := c.Bind(&product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.repo.Create(&product); err != nil {
		h.logger.Error("Failed to create product", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create product")
	}

	return c.JSON(http.StatusCreated, product)
}

// @Summary		Get product
// @Description	Get product by ID
// @Tags			products
// @Produce		json
// @Param			id	path		int	true	"Product ID"
// @Success		200	{object}	Product
// @Failure		404	{object}	middleware.ErrorResponse
// @Router			/products/{id} [get]
func (h *Handler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	product, err := h.repo.GetByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return c.JSON(http.StatusOK, product)
}

// @Summary		Update product
// @Description	Update product by ID (Admin only)
// @Tags			products
// @Accept			json
// @Produce		json
// @Param			id		path		int		true	"Product ID"
// @Param			product	body		Product	true	"Product object"
// @Success		200		{object}	Product
// @Failure		404		{object}	middleware.ErrorResponse
// @Router			/products/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	product, err := h.repo.GetByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	if err := c.Bind(product); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.repo.Update(product); err != nil {
		h.logger.Error("Failed to update product", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update product")
	}

	return c.JSON(http.StatusOK, product)
}

// @Summary		Delete product
// @Description	Delete product by ID (Admin only)
// @Tags			products
// @Param			id	path	int	true	"Product ID"
// @Success		204	"No Content"
// @Failure		404	{object}	middleware.ErrorResponse
// @Router			/products/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		h.logger.Error("Failed to delete product", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete product")
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary		List products
// @Description	Get a paginated list of products with optional filters
// @Tags			products
// @Produce		json
// @Param			page		query		int		false	"Page number"
// @Param			limit		query		int		false	"Items per page"
// @Param			min_price	query		number	false	"Minimum price"
// @Param			max_price	query		number	false	"Maximum price"
// @Param			search		query		string	false	"Search term"
// @Param			sort_by		query		string	false	"Sort by field (name, price, created_at)"
// @Param			sort_dir	query		string	false	"Sort direction (asc, desc)"
// @Success		200			{object}	middleware.PaginatedResponse
// @Router			/products [get]
func (h *Handler) List(c echo.Context) error {
	var query ListProductsQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 10
	}

	if query.SortBy != "" {
		validColumns := map[string]bool{
			"name": true, "price": true, "created_at": true,
		}
		if !validColumns[query.SortBy] {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid sort column")
		}
	}
	if query.SortDir != "" {
		query.SortDir = strings.ToUpper(query.SortDir)
		if query.SortDir != "ASC" && query.SortDir != "DESC" {
			query.SortDir = "ASC"
		}
	}

	products, total, err := h.repo.List(&query)
	if err != nil {
		h.logger.Error("Failed to list products",
			zap.Error(err),
			zap.Int("page", query.Page),
			zap.Int("limit", query.Limit))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list products")
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))

	response := map[string]interface{}{
		"products": products,
		"pagination": map[string]interface{}{
			"current_page": query.Page,
			"total_pages":  totalPages,
			"total_items":  total,
			"limit":        query.Limit,
		},
		"filters": map[string]interface{}{
			"min_price": query.MinPrice,
			"max_price": query.MaxPrice,
			"search":    query.Search,
			"sort_by":   query.SortBy,
			"sort_dir":  query.SortDir,
		},
	}

	return c.JSON(http.StatusOK, response)
}
