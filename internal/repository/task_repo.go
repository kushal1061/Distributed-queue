package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ShabnamHaque/task-queue/config"
	"github.com/ShabnamHaque/task-queue/internal/models"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get the MongoDB collection
func GetTaskCollection() *mongo.Collection {
	return config.MongoDB.Collection("tasks")
}

// InsertTask inserts a new task into the database
func InsertTask(task *models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := GetTaskCollection()
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt
	task.Status = "queued"
	task.ID = uuid.New().String() //assign uuid

	_, err := collection.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

// GetTaskByID retrieves a task by its ID
func GetTaskByID(taskID string) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := GetTaskCollection()
	var task models.Task
	err := collection.FindOne(ctx, bson.M{"_id": taskID}).Decode(&task)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("task not found")
	} else if err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTaskStatus updates the status of a task
func UpdateTaskStatus(taskID string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := GetTaskCollection()
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": taskID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
} 
func UpdateTaskCompletion(taskID string, status string) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := GetTaskCollection()
	var task models.Task
	err := collection.FindOne(context.TODO(), bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err // Task not found
		}
		return err // Other error
	}
	// Calculate time taken to complete the task
	completionTime := int(time.Since(task.CreatedAt).Seconds())

	update := bson.M{
		"$set": bson.M{
			"status":          status,
			"completion_time": completionTime, // Store the difference in seconds
			"updated_at":      time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": taskID}, update)
	return err
}

// GetAllTasks retrieves all tasks from the database
func GetAllTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GetTaskCollection()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
