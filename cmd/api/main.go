package main

import (
	"context"
	"fmt"
	"os/signal"
	"runtime"
	"syscall"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/MadAppGang/httplog/echolog"
	"github.com/common-nighthawk/go-figure"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	_ "github.com/nneji123/ecommerce-golang/docs"
	"github.com/nneji123/ecommerce-golang/internal/common/email"
	"github.com/nneji123/ecommerce-golang/internal/domain/order"
	"github.com/nneji123/ecommerce-golang/internal/domain/product"
	"github.com/nneji123/ecommerce-golang/internal/domain/user"

	"github.com/nneji123/ecommerce-golang/internal/config"
	"github.com/nneji123/ecommerce-golang/internal/db"
)

//	@title			GoCommerce API
//	@version		1.0
//	@description	GoCommerce API is a service for leads generation and outreach.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	GoCommerce Support
//	@contact.url	http://GoCommerce.com
//	@contact.email	contact@GoCommerce.com

//	@license.name	MIT
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host	localhost:8080
// @BasePath
// @schemes	http
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	database, err := db.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Defer database connection close
	sqlDB, err := database.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", zap.Error(err))
	}
	defer sqlDB.Close()

	validate := validator.New()

	e := echo.New()

	e.Use(echolog.LoggerWithName("ECHO NATIVE"))
	e.Use(middleware.Recover())

	// Test Routes
	e.GET("/ping", handleGetPing)
	e.GET("/", handleGetRoot)
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	userRepo := user.NewRepository(database)
	emailNotificationService, err := email.NewEmailNotificationService("smtp", &cfg)
	if err != nil {
		logger.Fatal("Failed to initialize email service", zap.Error(err))
	}
	emailService := user.NewEmailService(emailNotificationService, &cfg)
	authService := user.NewJWTService(cfg.JWTSecret)

	// Initialize handlers
	userHandler := user.NewHandler(
		userRepo,
		validate,
		authService,
		emailService,
		logger,
	)

	user.RegisterRoutes(e, userHandler)

	if err := database.AutoMigrate(&user.User{}, &product.Product{}, &order.Order{}, &order.OrderItem{}); err != nil {
		logger.Fatal("Failed to auto-migrate user model", zap.Error(err))
	}

	printServerDetails(cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: e,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Graceful shutdown
	shutdownServer(srv)
}

// HandlePing
func handleGetPing(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}

func handleGetRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Server is running")
}

func printServerDetails(cfg config.Config) {
	myFigure := figure.NewFigure("GoCommerce", "", true)
	myFigure.Print()

	fmt.Println("\nServer Information:")
	fmt.Printf("Version:       %s\n", "1.0.0")
	fmt.Printf("Environment:   %s\n", "Development")
	fmt.Printf("Platform:      %s\n", runtime.GOOS)
	fmt.Printf("Architecture:  %s\n", runtime.GOARCH)
	fmt.Printf("Server Port:   %s\n", cfg.ServerPort)

	// Get CPU information
	cpuInfo, err := cpu.Info()
	if err != nil {
		fmt.Println("Failed to get CPU information:", err)
		return
	}
	fmt.Println("\nCPU:")
	for _, info := range cpuInfo {
		fmt.Printf("  Model: %s\n", info.ModelName)
		fmt.Printf("  Cores: %d\n", info.Cores)
	}

	// Get memory information
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Failed to get memory information:", err)
		return
	}
	fmt.Println("\nMemory:")
	fmt.Printf("  Total: %d MB\n", memInfo.Total/1024/1024)
	fmt.Printf("  Free:  %d MB\n", memInfo.Free/1024/1024)
	fmt.Printf("  Used:  %d MB\n\n", memInfo.Used/1024/1024)
}

func shutdownServer(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server gracefully stopped")
}
