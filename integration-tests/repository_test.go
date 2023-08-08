//go:build integration

package integration_tests

import (
	"context"
	"fmt"
	"github.com/nekidaz/todolist/internal/entity"
	"github.com/nekidaz/todolist/internal/usecase/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
	"time"
)

type RepositorySuite struct {
	suite.Suite
	repository repo.TodoRepository
	ctx        context.Context
}

func (s *RepositorySuite) SetupSuite() {
	connectionString := fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	repoName := os.Getenv("MONGO_NAME")
	repoCollection := "test"

	repoInstance, err := repo.NewRepository(connectionString, repoName, repoCollection)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %s", err)
	}

	s.repository = repoInstance
}

func (s *RepositorySuite) TearDownSuite() {
	// Очистка ресурсов и отключение от MongoDB
	if s.repository != nil {
		if err := s.repository.Close(); err != nil {
			log.Fatalf("Ошибка при закрытии подключения к базе данных: %s", err)
		}
	}
}

func (s *RepositorySuite) TestCreateAndRetrieveTodos() {
	numTodos := 3
	todos := make([]*entity.Todo, numTodos)

	for i := 0; i < numTodos; i++ {
		newTodo := &entity.Todo{
			Title:     fmt.Sprintf("Test Task %d", i+1),
			ActiveAt:  time.Now(),
			Completed: false,
		}

		createdTodo, err := s.repository.CreateNewTodo(s.ctx, newTodo)
		assert.NoError(s.T(), err, "Ошибка создания задачи")
		todos[i] = createdTodo
	}

	for _, createdTodo := range todos {
		retrievedTodo, err := s.repository.GetTaskByID(s.ctx, createdTodo.ID)
		assert.NoError(s.T(), err, "Не удалось получить задачу из базы данных")

		assert.Equal(s.T(), createdTodo.Title, retrievedTodo.Title)
		assert.Equal(s.T(), createdTodo.ActiveAt.UTC().Format(time.RFC3339), retrievedTodo.ActiveAt.UTC().Format(time.RFC3339))
		assert.Equal(s.T(), createdTodo.Completed, retrievedTodo.Completed)
	}
}

func (s *RepositorySuite) TestUpdateTodo() {
	// Создание тестовой задачи
	newTodo := &entity.Todo{
		Title:     "Test Task",
		ActiveAt:  time.Now(),
		Completed: false,
	}
	createdTodo, err := s.repository.CreateNewTodo(s.ctx, newTodo)
	assert.NoError(s.T(), err, "Ошибка создания задачи")

	// Обновление задачи
	updatedTodo := &entity.Todo{
		ID:        createdTodo.ID,
		Title:     "Updated Test Task",
		ActiveAt:  time.Now().Add(24 * time.Hour),
		Completed: true,
	}
	updatedTodo, err = s.repository.UpdateTodo(s.ctx, updatedTodo.ID, updatedTodo)
	assert.NoError(s.T(), err, "Ошибка обновления задачи")

	// Проверка обновленных данных
	retrievedTodo, err := s.repository.GetTaskByID(s.ctx, createdTodo.ID)
	assert.NoError(s.T(), err, "Не удалось получить задачу из базы данных")

	assert.Equal(s.T(), updatedTodo.Title, retrievedTodo.Title)
	assert.Equal(s.T(), updatedTodo.ActiveAt.UTC().Format(time.RFC3339), retrievedTodo.ActiveAt.UTC().Format(time.RFC3339))
	assert.Equal(s.T(), updatedTodo.Completed, retrievedTodo.Completed)
}

func (s *RepositorySuite) TestDeleteTodo() {
	// Создание тестовой задачи
	newTodo := &entity.Todo{
		Title:     "Test Task",
		ActiveAt:  time.Now(),
		Completed: false,
	}
	createdTodo, err := s.repository.CreateNewTodo(s.ctx, newTodo)
	assert.NoError(s.T(), err, "Ошибка создания задачи")

	// Удаление задачи
	err = s.repository.DeleteTodo(s.ctx, createdTodo.ID)
	assert.NoError(s.T(), err, "Ошибка удаления задачи")

	// Попытка получения удаленной задачи
	_, err = s.repository.GetTaskByID(s.ctx, createdTodo.ID)
	assert.Error(s.T(), err, "Ожидаемая ошибка при извлечении удаленной задачи")
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
