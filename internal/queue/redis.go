package queue

import (
	"context"
	"log"
	"time"

	"github.com/ShabnamHaque/task-queue/config"
	"github.com/go-redis/redis/v8"
)

const QueueName = "tasks_queue"

func PushTask(taskID string) error { //enqueues a taskid into redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.RedisClient.RPush(ctx, QueueName, taskID).Err()
	if err != nil {
		log.Printf("🔴 Error pushing task %s to Redis: %v", taskID, err)
		return err
	}
	log.Printf("✅ Task %s added to queue", taskID)
	return nil
}
func PopTask() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the queue is empty
	queueLength, err := config.RedisClient.LLen(ctx, "tasks_queue").Result()
	if err != nil {
		log.Printf("🔴 Error checking queue length: %v", err)
		return "", err
	}
	if queueLength == 0 {
		log.Println("⚠️ Task queue is empty.")
		return "", nil
	}

	// Pop a task from the queue
	taskID, err := config.RedisClient.LPop(ctx, "tasks_queue").Result()
	if err == redis.Nil {
		return "", nil // Queue is empty
	} else if err != nil {
		log.Printf("🔴 Error popping task from Redis: %v", err)
		return "", err
	}

	log.Printf("🔄 Task %s dequeued for processing", taskID)
	return taskID, nil
}

/* Redis is an in-memory data structure that is used for faster access to data.
It is used to store data that needs to be accessed frequently and fast.
It is not used for storing large amounts of data. If you want to store and retrieve large amounts of
data you need to use a traditional database such as MongoDB or MYSQL. Redis provides a variety of data structures such as sets, strings,
 hashes, and lists.

redis lists - lpush and blpop used - In real, it works like linked lists under the hood,
which means adding or removing elements from the beginning or the end of the list is
an efficient operation here having a time complexity of O(1) that is constant.
*/
