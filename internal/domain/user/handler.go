package user

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/common/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type Handler struct {
	repo         Repository
	validator    *validator.Validate
	auth         AuthService
	emailService EmailService
	logger       *zap.Logger
}

func NewHandler(repo Repository, validator *validator.Validate, auth AuthService, emailService EmailService, logger *zap.Logger) *Handler {
	return &Handler{
		repo:         repo,
		validator:    validator,
		auth:         auth,
		emailService: emailService,
		logger:       logger,
	}
}

// generateToken creates a random token for email verification or password reset
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	Register a new user with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"Registration details"
//	@Success		201		{object}	User
//	@Failure		400		{object}	middleware.ErrorResponse
//	@Failure		409		{object}	middleware.ErrorResponse
//	@Router			/auth/register [post]
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

	token, err := generateToken()
	if err != nil {
		return err
	}

	expiryTime := time.Now().Add(24 * time.Hour)

	user := &User{
		Email:                   req.Email,
		Password:                req.Password,
		Name:                    req.Name,
		Role:                    "user",
		EmailVerificationToken:  &token,
		EmailVerificationExpiry: &expiryTime,
		PasswordResetToken:      nil,
		PasswordResetExpiry:     nil,
	}

	if err := h.repo.Create(user); err != nil {
		return err
	}

	// Send verification email
	if err := h.emailService.SendVerificationEmail(user.Email, token); err != nil {
		h.logger.Error("Failed to send verification email", zap.Error(err))
	}

	return c.JSON(http.StatusCreated, user)
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and return JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Login credentials"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	middleware.ErrorResponse
//	@Failure		401		{object}	middleware.ErrorResponse
//	@Router			/auth/login [post]
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

	if !user.IsEmailVerified {
		return echo.NewHTTPError(http.StatusUnauthorized, "email not verified")
	}

	token, err := h.auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  *user,
	})
}

// ConfirmRegistration godoc
//
//	@Summary		Verify email address
//	@Description	Verify user's email address using verification token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		VerifyEmailRequest	true	"Verification token"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	middleware.ErrorResponse
//	@Router			/auth/confirm-registration [post]
func (h *Handler) ConfirmRegistration(c echo.Context) error {
	var req VerifyEmailRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.repo.FindByEmailVerificationToken(req.Token)
	if err != nil {
		return err
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid or expired token")
	}

	user.IsEmailVerified = true
	user.EmailVerificationToken = nil  // Set to nil instead of empty string
	user.EmailVerificationExpiry = nil // Clear the expiry time as well

	if err := h.repo.Update(user); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "email verified successfully",
	})
}

// RequestPasswordReset godoc
//
//	@Summary		Request password reset
//	@Description	Send password reset email to user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ResetPasswordRequest	true	"User email"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	middleware.ErrorResponse
//	@Router			/auth/password-reset-request [post]
func (h *Handler) RequestPasswordReset(c echo.Context) error {
	var req ResetPasswordRequest
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

	// Always return success to prevent email enumeration
	if user == nil {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "if your email is registered, you will receive a password reset link",
		})
	}

	token, err := generateToken()
	if err != nil {
		return err
	}

	expiryTime := time.Now().Add(1 * time.Hour)
	user.PasswordResetToken = &token
	user.PasswordResetExpiry = &expiryTime

	if err := h.repo.Update(user); err != nil {
		return err
	}

	if err := h.emailService.SendPasswordResetEmail(user.Email, token); err != nil {
		h.logger.Error("Failed to send password reset email", zap.Error(err))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "if your email is registered, you will receive a password reset link",
	})
}

// ConfirmPasswordReset godoc
//
//	@Summary		Confirm password reset
//	@Description	Reset user's password using reset token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ResetPasswordConfirmRequest	true	"Reset token and new password"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	middleware.ErrorResponse
//	@Router			/auth/confirm-password-reset [post]
func (h *Handler) ConfirmPasswordReset(c echo.Context) error {
	var req ResetPasswordConfirmRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.repo.FindByPasswordResetToken(req.Token)
	if err != nil {
		return err
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid or expired token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.PasswordResetToken = nil  // Set to nil instead of empty string
	user.PasswordResetExpiry = nil // Set to nil instead of zero time

	if err := h.repo.Update(user); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "password reset successfully",
	})
}

// UserDetail godoc
//
//	@Summary		Get user details
//	@Description	Retrieve the logged-in user's details
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	User
//	@Failure		401	{object}	middleware.ErrorResponse
//	@Router			/user/detail [get]
func (h *Handler) UserDetail(c echo.Context) error {
	claims, ok := c.Get("userClaims").(*models.Claims)
	if !ok {
		h.logger.Error("Handler: Claims not found in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or missing user claims")
	}

	response := struct {
		UserID uint   `json:"user_id"`
		Email  string `json:"email"`
		Role   string `json:"role"`
	}{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}

	h.logger.Info("Handler: User details retrieved",
		zap.Uint("user_id", claims.UserID),
		zap.String("email", claims.Email))

	return c.JSON(http.StatusOK, response)
}
