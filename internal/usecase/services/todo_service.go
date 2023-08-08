package services

import (
	"context"
	"github.com/nekidaz/todolist/internal/entity"
	"github.com/nekidaz/todolist/internal/usecase/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TodoService interface {
	CreateNewTodo(ctx context.Context, title string, activeAt time.Time) (*entity.Todo, error)
	UpdateTodo(ctx context.Context, id primitive.ObjectID, title string, activeAt time.Time) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, id primitive.ObjectID) error
	MarkAsCompleted(ctx context.Context, id primitive.ObjectID) error
	GetAllTasks(ctx context.Context) ([]*entity.Todo, error)
	GetTasksByStatus(ctx context.Context, status string) ([]*entity.Todo, error)
}

type todoService struct {
	repo repo.TodoRepository
}

func NewTodoService(repo repo.TodoRepository) TodoService {
	return &todoService{
		repo: repo,
	}
}

func (s *todoService) CreateNewTodo(ctx context.Context, title string, activeAt time.Time) (*entity.Todo, error) {
	todo := entity.NewTodo(title, activeAt)
	return s.repo.CreateNewTodo(ctx, todo)
}

func (s *todoService) UpdateTodo(ctx context.Context, id primitive.ObjectID, title string, activeAt time.Time) (*entity.Todo, error) {
	todo := entity.NewTodo(title, activeAt)
	return s.repo.UpdateTodo(ctx, id, todo)
}

func (s *todoService) DeleteTodo(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteTodo(ctx, id)
}

func (s *todoService) MarkAsCompleted(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.MarkAsCompleted(ctx, id)
}

func (s *todoService) GetAllTasks(ctx context.Context) ([]*entity.Todo, error) {
	return s.repo.GetAllTasks(ctx)
}

func (s *todoService) GetTasksByStatus(ctx context.Context, status string) ([]*entity.Todo, error) {
	return s.repo.GetTasksByStatus(ctx, status)
}
