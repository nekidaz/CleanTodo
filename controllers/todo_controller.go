// controllers.go

package controllers

import (
	"RegionLabTZ/helpers"
	"RegionLabTZ/models"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": helpers.ErrParseActiveAt})
		return
	}

	// Создаем новую задачу через сервис
	todo, err := c.todoService.CreateNewTodo(ctx, requestBody.Title, activeAtTime)

	if err != nil {
		switch {
		case errors.Is(err, helpers.ErrTodoExists):
			ctx.Status(http.StatusNoContent)
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Задача с названием " + todo.Title + " успешно создана "})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": helpers.ErrParseActiveAt})
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

	ctx.Status(http.StatusNoContent)
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
	id, tasks, errReturned := c.processRequestID(ctx)
	if errReturned {
		return
	}
	ctx.JSON(http.StatusOK, tasks[id])
}

func (c *TodoController) processRequestID(ctx *gin.Context) (id int, tasks []*models.Todo, errReturned bool) {
	idStr := ctx.Param("ID")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		defer ctx.JSON(http.StatusBadRequest, gin.H{"error": helpers.ErrInvalidID})
		return 0, nil, true
	}

	id = id - 1

	tasks, err = c.todoService.GetAllTasks(ctx)
	if err != nil {
		defer ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 0, nil, true
	}

	if id < 0 || id >= len(tasks) {
		defer ctx.JSON(http.StatusNotFound, gin.H{"error": helpers.ErrTaskNotFound})
		return 0, nil, true
	}

	return id, tasks, false
}
