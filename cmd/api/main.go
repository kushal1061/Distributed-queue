package main

import (
	"fmt"
	"log"

	api "github.com/ShabnamHaque/task-queue/api"

	"github.com/ShabnamHaque/task-queue/config"
)

func main() {
	fmt.Println("🚀 Starting API Server...")

	config.LoadEnv()
	config.InitMongoDB()
	config.InitRedis()
	err := api.Run()
	if err != nil {
		log.Fatalf("❌ API Server failed: %v", err)
	}
}
