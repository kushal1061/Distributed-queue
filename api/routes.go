package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/submit-task", SubmitTask)
	router.GET("/task/status/:task_id", GetTaskStatus)
	router.GET("/tasks/status", GetAllTasks)
}
