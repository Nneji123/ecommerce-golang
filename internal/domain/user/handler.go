package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/middleware"
)

// @Summary Register new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} User
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 409 {object} middleware.ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existing, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return echo.NewHTTPError(http.StatusConflict, "email already exists")
	}

	user := &User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     "user",
	}

	if err := h.repo.Create(user); err != nil {
		return err
	}

	// Send welcome email
	err = h.emailService.SendEmail(
		"Welcome to Our Store",
		"welcome.mjml",
		user.Email,
		map[string]interface{}{
			"Name": user.Name,
		},
		nil,
	)
	if err != nil {
		// Log error but don't fail registration
		h.logger.Error("Failed to send welcome email:", err)
	}

	return c.JSON(http.StatusCreated, user)
}

// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Generate JWT token
	token, err := h.auth.GenerateToken(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  *user,
	})
}