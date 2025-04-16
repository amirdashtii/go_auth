package main

import (
	"log"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/controller"
	"github.com/gin-gonic/gin"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	r := gin.Default()
	controller.NewAuthRoutes(r)
	controller.NewUserRoutes(r)
	controller.NewAdminRoutes(r)

	log.Printf("Server is running on port %s\n", config.Server.Port)
	log.Printf("http://localhost:%s\n", config.Server.Port)

	r.Run(":" + config.Server.Port)
}
