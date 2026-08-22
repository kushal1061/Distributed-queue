package main

import (
	"fmt"
	"time"

	"github.com/ShabnamHaque/task-queue/config"
	"github.com/ShabnamHaque/task-queue/internal/worker"
)

func main() {
	fmt.Println("🚀 Starting Worker Service...")
	config.LoadEnv()
	config.InitMongoDB()
	config.InitRedis()

	for i := 1; i <= 3; i++ { // Start 3 workers
		go func(workerID int) {
			worker.StartWorker(workerID)
		}(i)
		time.Sleep(1 * time.Second) // Small delay to stagger worker startup
	}
	select {}
}
