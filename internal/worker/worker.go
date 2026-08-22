package worker

import (
	"fmt"
	"log"
	"sync"
	"time"

	redis "github.com/ShabnamHaque/task-queue/internal/queue"
	"github.com/ShabnamHaque/task-queue/internal/repository"
)

func StartWorker(workerCount int) {
	var wg sync.WaitGroup

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			PopTasks(workerID)
		}(i)
	}

	wg.Wait()
}

// PopTasks fetches tasks from Redis and updates MongoDB status
func PopTasks(workerID int) {
	fmt.Printf("🚀 Worker %d started...\n", workerID)

	for {
		// Fetch a task from Redis
		taskID, err := redis.PopTask()
		if err != nil {
			log.Printf("❌ Worker %d: Failed to fetch task from queue: %v", workerID, err)
			continue
		}
		if taskID == "" {
			log.Printf("⚠️ Worker %d: No tasks in queue, waiting...", workerID)
			time.Sleep(5 * time.Second) // Avoid busy-waiting
			continue
		}
		task, err := repository.GetTaskByID(taskID)
		if err != nil {
			log.Printf("❌ Worker %d: Failed to retrieve task %s: %v", workerID, taskID, err)
			continue
		}
		fmt.Printf("👷 Worker %d: Processing Task: %+v\n", workerID, task)

		// Update task status to "in-progress" in MongoDB
		task.Status = "in-progress"
		task.UpdatedAt = time.Now()
		err = repository.UpdateTaskStatus(task.ID, "in-progress")
		if err != nil {
			log.Printf("❌ Worker %d: Failed to update task %s to in-progress: %v", workerID, task.ID, err)
			continue
		}

		// Simulate task processing time
		time.Sleep(time.Duration(task.TimeToComplete) * time.Second)

		// Simulate failure (10% chance of failure)
		// Simulate failure (10% chance of failure)
		if time.Now().Unix()%10 == 0 {
			task.Status = "failed"
			err = repository.UpdateTaskStatus(task.ID, "failed")
			if err != nil {
				log.Printf("❌ Worker %d: Failed to update task %s to failed: %v", workerID, task.ID, err)
			}
			// Requeue the failed task using InsertTask function
			insertErr := repository.InsertTask(task)
			insertErr = redis.PushTask(task.ID)
			if insertErr != nil {
				log.Printf("❌ Worker %d: Failed to requeue task %s: %v", workerID, task.ID, insertErr)
			} else {
				log.Printf("🔄 Worker %d: Task %s requeued for retry", workerID, task.ID)
			}
			continue
		}

		// Mark task as completed in MongoDB
		err = repository.UpdateTaskCompletion(task.ID, "success")
		if err != nil {
			log.Printf("❌ Worker %d: Failed to update task %s to success: %v", workerID, task.ID, err)
			continue
		}

		fmt.Printf("✅ Worker %d: Task %s completed in %f seconds\n", workerID, task.ID, time.Since(task.CreatedAt).Seconds())
	}
}
