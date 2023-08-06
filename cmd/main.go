// main.go

package main

import (
	"RegionLabTZ/controllers"
	"RegionLabTZ/repositories"
	service "RegionLabTZ/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к MongoDB
	dbConnectionString := "mongodb://localhost:27017" // Замените на свой URL MongoDB
	dbName := "reloglab"                              // Замените на свое имя базы данных
	collectionName := "todolists"                     // Замените на имя коллекции

	repo, err := repositories.NewRepository(dbConnectionString, dbName, collectionName)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Создание сервиса и контроллера
	todoService := service.NewTodoService(repo)
	todoController := controllers.NewTodoController(todoService)

	// Создание маршрутов и запуск сервера
	r := gin.Default()

	// Обработчики для CRUD операций
	r.GET("/api/todo-list/tasks/:ID", todoController.GetTaskByID)
	r.GET("/api/todo-list/tasks", todoController.GetAllTask)

	r.POST("/api/todo-list/tasks", todoController.CreateNewTodoHandler)

	r.DELETE("/api/todo-list/tasks/:ID", todoController.DeleteTodoHandler)

	r.PUT("/api/todo-list/tasks/:ID", todoController.UpdateTodoHandler)
	//тут патч подходит лучше так как меняем 1 поле

	r.PATCH("/api/todo-list/tasks/:ID/done", todoController.MarkAsCompletedHandler)

	r.Run(":8080")
}
