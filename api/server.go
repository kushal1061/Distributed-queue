package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// Run starts the API server
func Run() error {
	router := gin.Default()
	SetupRoutes(router)
	port := "8080"
	fmt.Printf("🚀 API Server is running on port %s\n", port)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("❌ Failed to start API server: %v", err)
	}
	return err
}
