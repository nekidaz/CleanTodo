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

	r.GET("/api/todo-list/tasks/:ID", todoController.GetTaskByID)
	r.GET("/api/todo-list/tasks", todoController.GetTasksByStatusHandler)

	r.POST("/api/todo-list/tasks", todoController.CreateNewTodoHandler)

	r.DELETE("/api/todo-list/tasks/:ID", todoController.DeleteTodoHandler)

	r.PUT("/api/todo-list/tasks/:ID", todoController.UpdateTodoHandler)

	r.PATCH("/api/todo-list/tasks/:ID/done", todoController.MarkAsCompletedHandler)

	r.Run(":8080")
}
