package main

import (
	"log"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/controller"
	"github.com/amirdashtii/go_auth/controller/middleware"
	"github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/gin-gonic/gin"
)

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
	}
	appLogger, err := logger.NewFileLogger(loggerConfig)
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}

	// Initialize router
	r := gin.New() // Use gin.New() instead of gin.Default() to have more control
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(appLogger))

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
