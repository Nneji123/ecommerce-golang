package main

import (
	"context"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/swaggo/echo-swagger"
	"os/signal"
	"runtime"
	"syscall"

	_ "github.com/nneji123/ecommerce-golang/docs"

	"github.com/MadAppGang/httplog/echolog"
	"github.com/SporkHubr/echo-http-cache"
	"github.com/SporkHubr/echo-http-cache/adapter/redis"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nneji123/ecommerce-golang/internal/config"
	"github.com/nneji123/ecommerce-golang/internal/db"
	"github.com/nneji123/ecommerce-golang/internal/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

//	@title			Ecommerce API
//	@version		1.0
//	@description	Leadz Aura API is a service for leads generation and outreach..
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Ecommerce Support
//	@contact.url	http://Ecommerce.com
//	@contact.email	contact@Ecommerce.com

//	@license.name	MIT
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	// Connect to the database
	_, err = db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}()

	ringOpt := &redis.RingOptions{
		Addrs: map[string]string{
			"server": cfg.RedisAddr,
		},
	}
	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(redis.NewAdapter(ringOpt)),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Instatiate Echo Instance
	e := echo.New()

	// Middleware
	e.Use(echolog.LoggerWithName("ECHO NATIVE"))
	e.Use(middleware.Recover())
	e.Use(cacheClient.Middleware())
	// e.Use(middleware.AddTrailingSlash())
	e.Use(middlewares.CorsWithConfig(cfg))
	e.Use(middlewares.RateLimiterMiddleware(getRateLimitedRoutes()))

	// Test Routes
	e.GET("/ping", handleGetPing)
	e.GET("/", handleGetRoot)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Create a new asynq client for enqueuing tasks.
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisAddr})
	defer client.Close()

	// Print server details
	printServerDetails(cfg)

	// Start server in a separate goroutine
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

// List of open routes
func getOpenRoutes() []string {
	return []string{"/ping", "/send-email", "/swagger/*", "/healthcheck"}
}

// Define list of rate-limited routes
func getRateLimitedRoutes() []string {
	return []string{"/ping", "/healthcheck"}
}

// HandlePing
func handleGetPing(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}

func handleGetRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Server is running")
}

func printServerDetails(cfg config.Config) {
	myFigure := figure.NewFigure("Leadz Aura", "", true)
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
