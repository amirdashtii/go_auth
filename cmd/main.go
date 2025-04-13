package main

import (
	"log"

	"github.com/amirdashtii/go_auth/controller"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := gin.Default()
	controller.NewAuthRoutes(r)
	controller.NewUserRoutes(r)
	controller.NewAdminRoutes(r)

	log.Println("Server is running on port 8080")
	log.Println("http://localhost:8080")
	
	r.Run(":8080")
}
