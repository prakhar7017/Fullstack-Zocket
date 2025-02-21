package main

import (
	"fmt"
	"log"
	"taskmanagement/config"
	models "taskmanagement/model"
	"taskmanagement/routes"
	"taskmanagement/services"

	"github.com/gin-contrib/cors"
)

func main() {
	config.InitDB()
	services.InitAI()

	err := config.DB.AutoMigrate(&models.User{}, &models.Task{}, &models.TaskHistory{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	r := routes.SetupRouter()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	fmt.Println("Database migration completed successfully!")
	fmt.Println("Server is running on port 8080")
	r.Run(":8080")
}
