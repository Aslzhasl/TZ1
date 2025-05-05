package main

import (
	"TZ/internal"
	"TZ/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {

	db.InitPostgres()

	r := gin.Default()
	internal.SetupRoutes(r)
	r.Run(":8080")
}
