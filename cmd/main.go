package main

import (
	"RegionLabTZ/controllers"
	"RegionLabTZ/helpers"
	"RegionLabTZ/repositories"
	service "RegionLabTZ/services"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/swag/example/basic/docs"
	"log"
)

var config helpers.Config
var err error

func init() {

	config, err = helpers.ConfigSetup()

	if err != nil {
		log.Fatalf("Error setting up configuration: %s", err)
	}
}

func main() {

	config, err := helpers.ConfigSetup()
	if err != nil {
		log.Fatalf("Error setting up configuration: %s", err)
	}
	repo, err := repositories.NewRepository(config.DBConnectionString, config.DBName, config.CollectionName)

	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Создание сервиса и контроллера
	todoService := service.NewTodoService(repo)
	todoController := controllers.NewTodoController(todoService)

	// Создание маршрутов и запуск сервера
	r := gin.Default()

	api := r.Group("/api/todo-list")

	{
		api.GET("/tasks/:ID", todoController.GetTaskByID)
		api.GET("/tasks", todoController.GetTasksByStatusHandler)
		api.GET("/tasks/all", todoController.GetAllTasks)

		api.POST("/tasks", todoController.CreateNewTodoHandler)
		api.DELETE("/tasks/:ID", todoController.DeleteTodoHandler)
		api.PUT("/tasks/:ID", todoController.UpdateTodoHandler)
		api.PATCH("/tasks/:ID/done", todoController.MarkAsCompletedHandler)

		r.Run(":8080")
	}
}
