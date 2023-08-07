package test

import (
	"RegionLabTZ/models"
	"RegionLabTZ/repositories"
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
	"time"
)

var pool *dockertest.Pool
var repo repositories.TodoRepository
var container *dockertest.Resource

// тут операции по типу создание и удаление контейнера
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

// создание контейнера бд и подключеник к нему
func setup() {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	container, err = pool.Run("mongo", "4.4", nil)
	if err != nil {
		log.Fatalf("Could not start MongoDB container: %s", err)
	}

	if err := pool.Retry(func() error {
		uri := fmt.Sprintf("mongodb://localhost:%s", container.GetPort("27017/tcp"))
		repo, err = repositories.NewRepository(uri, "test_db", "test_collection")
		if err != nil {
			return err
		}

		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			return err
		}
		defer client.Disconnect(context.Background())

		return client.Ping(context.Background(), nil)
	}); err != nil {
		log.Fatalf("Could not connect to MongoDB: %s", err)
	}
}

// удаление конейнера
func teardown() {
	if err := pool.Purge(container); err != nil {
		log.Fatalf("Could not purge MongoDB container: %s", err)
	}
}

// создаем задачи и проверяем (мне казалось что тест ибо слишком быстро все работает и решил сделать так )
func TestCreateAndRetrieveTodos(t *testing.T) {
	numTodos := 3
	todos := make([]*models.Todo, numTodos)

	for i := 0; i < numTodos; i++ {
		newTodo := &models.Todo{
			Title:     fmt.Sprintf("Test Task %d", i+1),
			ActiveAt:  time.Now(),
			Completed: false,
		}

		createdTodo, err := repo.CreateNewTodo(context.Background(), newTodo)
		if err != nil {
			t.Fatalf("Failed to create new task: %s", err)
		}
		todos[i] = createdTodo
	}

	for _, createdTodo := range todos {
		retrievedTodo, err := repo.GetTaskByID(context.Background(), createdTodo.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve task from database: %s", err)
		}

		assert.Equal(t, createdTodo.Title, retrievedTodo.Title)
		assert.Equal(t, createdTodo.ActiveAt.UTC().Format(time.RFC3339), retrievedTodo.ActiveAt.UTC().Format(time.RFC3339))
		assert.Equal(t, createdTodo.Completed, retrievedTodo.Completed)
	}
}

// обновление
func TestUpdateTodo(t *testing.T) {
	ctx := context.Background()

	newTodo := &models.Todo{
		Title:     "Test Task",
		ActiveAt:  time.Now(),
		Completed: false,
	}
	createdTodo, err := repo.CreateNewTodo(ctx, newTodo)
	if err != nil {
		t.Fatalf("Failed to create new task: %s", err)
	}

	updatedTodo := &models.Todo{
		ID:        createdTodo.ID,
		Title:     "Updated Test Task",
		ActiveAt:  time.Now().Add(24 * time.Hour),
		Completed: true,
	}
	_, err = repo.UpdateTodo(ctx, updatedTodo.ID, updatedTodo)
	if err != nil {
		t.Fatalf("Failed to update task: %s", err)
	}

	retrievedTodo, err := repo.GetTaskByID(ctx, createdTodo.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve task from database: %s", err)
	}

	assert.Equal(t, updatedTodo.Title, retrievedTodo.Title)
	assert.Equal(t, updatedTodo.ActiveAt.UTC().Format(time.RFC3339), retrievedTodo.ActiveAt.UTC().Format(time.RFC3339))
	assert.Equal(t, updatedTodo.Completed, retrievedTodo.Completed)
}

// удаление
func TestDeleteTodo(t *testing.T) {
	ctx := context.Background()

	newTodo := &models.Todo{
		Title:     "Test Task",
		ActiveAt:  time.Now(),
		Completed: false,
	}
	createdTodo, err := repo.CreateNewTodo(ctx, newTodo)
	if err != nil {
		t.Fatalf("Failed to create new task: %s", err)
	}

	err = repo.DeleteTodo(ctx, createdTodo.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %s", err)
	}

	_, err = repo.GetTaskByID(ctx, createdTodo.ID)
	assert.Error(t, err, "Expected error when retrieving deleted task")
}

// ну тут уже ждем ошибку
func TestUpdateTodo_ErrorNotFound(t *testing.T) {
	ctx := context.Background()

	nonExistentID, _ := primitive.ObjectIDFromHex("603f650f8b20bd000e1b857a") // Non-existent ID
	updatedTodo := &models.Todo{
		ID:        nonExistentID,
		Title:     "Updated Test Task",
		ActiveAt:  time.Now().Add(24 * time.Hour),
		Completed: true,
	}
	_, err := repo.UpdateTodo(ctx, updatedTodo.ID, updatedTodo)

	assert.NotNil(t, err, "Expected error for updating non-existent task")
}

// также ожидаем ошибку
func TestUpdateTodo_ErrorValidation(t *testing.T) {
	ctx := context.Background()

	newTodo := &models.Todo{
		Title:     "Test Task",
		ActiveAt:  time.Now(),
		Completed: false,
	}
	createdTodo, err := repo.CreateNewTodo(ctx, newTodo)
	if err != nil {
		t.Fatalf("Failed to create new task: %s", err)
	}

	updatedTodo := &models.Todo{
		ID:        createdTodo.ID,
		Title:     "", // Invalid title
		ActiveAt:  time.Now().Add(24 * time.Hour),
		Completed: true,
	}
	_, err = repo.UpdateTodo(ctx, updatedTodo.ID, updatedTodo)

	assert.NotNil(t, err, "Expected error for updating task with invalid data")
}

// удалем то чего нет и ожидаем ошибку
func TestDeleteTodo_ErrorNotFound(t *testing.T) {
	ctx := context.Background()

	nonExistentID, _ := primitive.ObjectIDFromHex("603f650f8b20bd000e1b857a") // Non-existent ID
	err := repo.DeleteTodo(ctx, nonExistentID)

	assert.NotNil(t, err, "Expected error for deleting non-existent task")
}
