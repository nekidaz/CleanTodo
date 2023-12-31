package controllers

import (
	"errors"
	"github.com/nekidaz/todolist/internal/entity"
	"github.com/nekidaz/todolist/internal/usecase/services"
	errors2 "github.com/nekidaz/todolist/pkg/errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TodoController struct {
	todoService services.TodoService
}

func NewTodoController(todoService services.TodoService) *TodoController {
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

	// тут парсим ибо формат не такой получаем как в тз
	activeAtTime, err := time.Parse("2006-01-02", requestBody.ActiveAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors2.ErrParseActiveAt})
		return
	}

	// создаем задачу через сервис
	todo, err := c.todoService.CreateNewTodo(ctx, requestBody.Title, activeAtTime)

	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrTodoExists):
			ctx.JSON(http.StatusNoContent, gin.H{"error": errors2.ErrAlreadyExist})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Задача с заголовком: " + todo.Title + " успешно создана"})
}

func (c *TodoController) UpdateTodoHandler(ctx *gin.Context) {

	id, tasks, errReturned := c.processRequestID(ctx)
	if errReturned {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors2.ErrParseActiveAt})
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
	id, tasks, errReturned := c.processRequestID(ctx)

	if errReturned {
		return
	}

	err := c.todoService.DeleteTodo(ctx, tasks[id].ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *TodoController) MarkAsCompletedHandler(ctx *gin.Context) {
	id, tasks, errReturned := c.processRequestID(ctx)
	if errReturned {
		return
	}

	err := c.todoService.MarkAsCompleted(ctx, tasks[id].ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *TodoController) GetTaskByID(ctx *gin.Context) {
	// Получаем ID задачи из параметра в URL
	id, tasks, errReturned := c.processRequestID(ctx)
	if errReturned {
		return
	}
	ctx.JSON(http.StatusOK, tasks[id])
}

func (c *TodoController) GetAllTasks(ctx *gin.Context) {
	tasks, err := c.todoService.GetAllTasks(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func (c *TodoController) GetTasksByStatusHandler(ctx *gin.Context) {
	status := ctx.Query("status")

	if status == "" {
		status = "active"
	}

	tasks, err := c.todoService.GetTasksByStatus(ctx, status)

	if err != nil {
		// если не все ок то показываем кастомную ошибку
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ну тут если задач нет то пустой массив
	if len(tasks) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
		return
	}
	for i, task := range tasks {
		if task.ActiveAt.Weekday() == time.Saturday || task.ActiveAt.Weekday() == time.Sunday {
			tasks[i].Title = "ВЫХОДНОЙ - " + task.Title
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// это для того чтобы получать данные в виде массива так как выполнять разные операции будет легчо выполнять по айдишкику в массиве
func (c *TodoController) processRequestID(ctx *gin.Context) (id int, tasks []*entity.Todo, errReturned bool) {
	idStr := ctx.Param("ID")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		defer ctx.JSON(http.StatusBadRequest, gin.H{"error": errors2.ErrInvalidID})
		return 0, nil, true
	}

	id = id - 1

	tasks, err = c.todoService.GetAllTasks(ctx)
	if err != nil {
		defer ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 0, nil, true
	}

	if id < 0 || id >= len(tasks) {
		defer ctx.JSON(http.StatusNotFound, gin.H{"error": errors2.ErrTaskNotFound})
		return 0, nil, true
	}

	return id, tasks, false
}
