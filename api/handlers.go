package api

import (
	"log"
	"net/http"

	"github.com/ShabnamHaque/task-queue/internal/models"
	queue "github.com/ShabnamHaque/task-queue/internal/queue"
	repository "github.com/ShabnamHaque/task-queue/internal/repository"

	"github.com/gin-gonic/gin"
)

func SubmitTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := repository.InsertTask(&task)

	if err != nil {
		log.Printf("cannot add task to mongodb : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert task into db"})
		return
	}

	err = queue.PushTask(task.ID)
	if err != nil {
		log.Printf("cannot add task to queue : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert task into db"})
		return

	}

	log.Println("Task submitted to queue and mongodb ", task.ID)
	c.JSON(http.StatusAccepted, gin.H{"message": "Task submitted", "task_id": task.ID})
}

func GetTaskStatus(c *gin.Context) { // GetTaskStatus retrieves the status of a task from mongodb
	taskID := c.Param("task_id")
	log.Println("task status for id ", taskID, " required")
	var task *models.Task
	if taskID == "" {
		log.Printf("error task id not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id is required"})
		return
	}
	task, err := repository.GetTaskByID(taskID)
	if err != nil {
		log.Printf("error %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"task_id": task.ID,
		"status":  task.Status,
	})
}
func GetAllTasks(c *gin.Context) {
	var tasks []models.Task
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	tasks, err := repository.GetAllTasks()
	if err != nil {
		log.Printf("error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to fetch all tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks, "error": nil})
}
