package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"time"

	"go-app/config"
	"go-app/database"
	"go-app/domain"
	"go-app/internal/logging"
	httpRepo "go-app/internal/repository/http"
	"go-app/internal/repository/postgres"
	"go-app/internal/rest"
	"go-app/internal/rest/middleware"
	"go-app/service"

	_ "go-app/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	config.LoadEnv()
}

func main() {
	// Initialize logging configuration
	config.SetupLogging()

	dbPool, err := database.SetupPgxPool()
	if err != nil {
		logging.LogError(context.Background(), err, "database_setup")
		os.Exit(1)
	}
	defer dbPool.Close()

	e := echo.New()
	e.HideBanner = true

	e.Logger.SetOutput(os.Stdout)
	e.Logger.SetLevel(0)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	e.Use(middleware.RequestIDMiddleware())
	e.Use(middleware.SlogLoggerMiddleware())
	appEnv := os.Getenv(string(config.AppEnvKey))
	insecureEnv := []string{"local", "development"}
	if !slices.Contains(insecureEnv, appEnv) {
		e.Use(middleware.Cors())
		e.Use(middleware.SecurityHeadersMiddleware())
	}
	e.Use(middleware.CompressionMiddleware())
	e.Use(middleware.RateLimitMiddleware(10.0, 20))
	e.Use(middleware.TimeoutMiddleware(180 * time.Second))

	// Register the routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, domain.Response{
			Code:    200,
			Message: "All is well!",
		})
	})

	journalRepo := postgres.NewJournalRepository(dbPool)
	embeddingHttp := httpRepo.NewEmbeddingHTTPRepository()
	journalService := service.NewJournalService(journalRepo, embeddingHttp)

	// Swagger
	enableSwagger := os.Getenv("ENABLE_SWAGGER")
	if enableSwagger == "true" {
		// @securityDefinitions.apikey BearerAuth
		// @in header
		// @name Authorization
		// @description Enter your bearer token in the format **Bearer <token>**
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	apiV1 := e.Group("/api/v1")
	usersGroup := apiV1.Group("/journals")

	rest.NewJournalHandler(usersGroup, journalService)

	// Get host from environment variable, default to 127.0.0.1 if not set
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	// Get port from environment variable, default to 8000 if not set
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	// Server address and port to listen on
	serverAddr := fmt.Sprintf("%s:%s", host, port)

	go func() {
		logging.LogInfo(ctx, "Server starting", slog.String("address", serverAddr))
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			logging.LogError(ctx, err, "server_start")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logging.LogInfo(ctx, "Shutting down server gracefully...")
	if err := e.Shutdown(ctx); err != nil {
		logging.LogError(ctx, err, "server_shutdown")
	}
}
