// todo_service.go

package service

import (
	"RegionLabTZ/models"
	"RegionLabTZ/repositories"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TodoService interface {
	CreateNewTodo(ctx context.Context, title string, activeAt time.Time) (*models.Todo, error)
	UpdateTodo(ctx context.Context, id primitive.ObjectID, title string, activeAt time.Time) (*models.Todo, error)
	DeleteTodo(ctx context.Context, id primitive.ObjectID) error
	MarkAsCompleted(ctx context.Context, id primitive.ObjectID) error
	GetTasksByStatus(ctx context.Context, status string) ([]*models.Todo, error)
	GetAllTasks(ctx context.Context) ([]*models.Todo, error)
	GetTaskByID(ctx context.Context, id primitive.ObjectID) (*models.Todo, error)
}

type todoService struct {
	repo repositories.TodoRepository
}

func NewTodoService(repo repositories.TodoRepository) TodoService {
	return &todoService{
		repo: repo,
	}
}

func (s *todoService) CreateNewTodo(ctx context.Context, title string, activeAt time.Time) (*models.Todo, error) {
	todo := models.NewTodo(title, activeAt)
	return s.repo.CreateNewTodo(ctx, todo)
}

func (s *todoService) UpdateTodo(ctx context.Context, id primitive.ObjectID, title string, activeAt time.Time) (*models.Todo, error) {
	todo := models.NewTodo(title, activeAt)
	return s.repo.UpdateTodo(ctx, id, todo)
}

func (s *todoService) DeleteTodo(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteTodo(ctx, id)
}

func (s *todoService) MarkAsCompleted(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.MarkAsCompleted(ctx, id)
}

func (s *todoService) GetTasksByStatus(ctx context.Context, status string) ([]*models.Todo, error) {
	return s.repo.GetTasksByStatus(ctx, status)
}

func (s *todoService) GetAllTasks(ctx context.Context) ([]*models.Todo, error) {
	return s.repo.GetAllTasks(ctx)
}

func (s *todoService) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*models.Todo, error) {
	return s.repo.GetTaskByID(ctx, id)
}
