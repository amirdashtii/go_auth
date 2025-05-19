package main

import (
	"log"
	"os"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/controller"
	"github.com/amirdashtii/go_auth/controller/middleware"
	_ "github.com/amirdashtii/go_auth/docs"
	"github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Go Auth API
// @version         1.0
// @description     A simple authentication service in Go.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize logger
	loggerConfig := ports.LoggerConfig{
		Level:       "info",
		Environment: config.Environment,
		ServiceName: "go_auth",
		Output:      os.Stdout,
	}
	appLogger := logger.NewZerologLogger(loggerConfig)

	// Initialize router
	r := gin.New() // Use gin.New() instead of gin.Default() to have more control
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(appLogger))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
	controller.NewAuthRoutes(r)
	controller.NewUserRoutes(r)
	controller.NewAdminRoutes(r)

	appLogger.Info("Server is starting", ports.F("port", config.Server.Port))
	appLogger.Info("Server URL", ports.F("url", "http://localhost:"+config.Server.Port))

	if err := r.Run(":" + config.Server.Port); err != nil {
		appLogger.Fatal("Server failed to start", ports.F("error", err))
	}
}
