package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nekidaz/todolist/config"
	"github.com/nekidaz/todolist/internal/controllers"
	"github.com/nekidaz/todolist/internal/usecase/repo"
	"github.com/nekidaz/todolist/internal/usecase/services"
	"log"
)

func main() {

	config, err := config.ConfigSetup()
	if err != nil {
		log.Fatalf("Ошибка при настройке конфигурации: %s", err)
	}

	repo, err := repo.NewRepository(config.DBConnectionString, config.DBName, config.CollectionName)

	if err != nil {
		log.Fatalf("Ошибка при подключении к MongoDB: %v", err)
	}

	// Создание сервиса и контроллера
	todoService := services.NewTodoService(repo)
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

	}

	r.Run(":8080")
}
