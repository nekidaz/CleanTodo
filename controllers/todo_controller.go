// controllers.go

package controllers

import (
	"RegionLabTZ/repositories"
	service "RegionLabTZ/services"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TodoController struct {
	todoService service.TodoService
}

func NewTodoController(todoService service.TodoService) *TodoController {
	return &TodoController{
		todoService: todoService,
	}
}

func (c *TodoController) CreateNewTodoHandler(ctx *gin.Context) {
	var requestBody struct {
		Title    string `json:"title" binding:"required,max=200"`
		ActiveAt string `json:"activeAt" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразовываем строку ActiveAt в формат времени time.Time
	activeAtTime, err := time.Parse("2006-01-02", requestBody.ActiveAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse ActiveAt"})
		return
	}

	// Создаем новую задачу через сервис
	todo, err := c.todoService.CreateNewTodo(ctx, requestBody.Title, activeAtTime)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrTodoExists):
			ctx.Status(http.StatusNoContent)
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Задача с названием " + todo.Title + " успешно создана "})
}

func (c *TodoController) UpdateTodoHandler(ctx *gin.Context) {
	idStr := ctx.Param("ID")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id = id - 1

	tasks, err := c.todoService.GetAllTasks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if id < 0 || id >= len(tasks) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var requestBody struct {
		Title    string `json:"title" binding:"required,max=200"`
		ActiveAt string `json:"activeAt" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activeAtTime, err := time.Parse("2006-01-02", requestBody.ActiveAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse ActiveAt"})
		return
	}

	todo, err := c.todoService.UpdateTodo(ctx, tasks[id].ID, requestBody.Title, activeAtTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (c *TodoController) DeleteTodoHandler(ctx *gin.Context) {
	idStr := ctx.Param("ID")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id = id - 1

	tasks, err := c.todoService.GetAllTasks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if id < 0 || id >= len(tasks) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	err = c.todoService.DeleteTodo(ctx, tasks[id].ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *TodoController) MarkAsCompletedHandler(ctx *gin.Context) {
	idStr := ctx.Param("ID")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id = id - 1

	tasks, err := c.todoService.GetAllTasks(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if id < 0 || id >= len(tasks) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	err = c.todoService.MarkAsCompleted(ctx, tasks[id].ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *TodoController) GetAllTask(ctx *gin.Context) {
	tasks, err := c.todoService.GetAllTasks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ActiveAt.Before(tasks[j].ActiveAt)
	})
	ctx.JSON(http.StatusOK, tasks)
}

func (c *TodoController) GetTaskByID(ctx *gin.Context) {
	// Получаем ID задачи из параметра в URL
	idStr := ctx.Param("ID")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id = id - 1

	tasks, err := c.todoService.GetAllTasks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if id < 0 || id >= len(tasks) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, tasks[id])
}
