package unit_tests

import (
	"github.com/nekidaz/todolist/internal/entity"
	"github.com/nekidaz/todolist/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateTodo(t *testing.T) {
	activationTime := time.Now().Add(time.Hour * 24)
	title := "Test Todo"
	todo := entity.NewTodo(title, activationTime)

	assert.Equal(t, title, todo.Title)
	assert.False(t, todo.Completed)
	assert.True(t, todo.ActiveAt.Equal(activationTime))
	assert.WithinDuration(t, time.Now(), todo.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)
}

func TestMarkAsCompleted(t *testing.T) {
	todo := entity.NewTodo("Test Todo", time.Now())
	assert.False(t, todo.Completed)
	assert.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)

	todo.MarkAsCompleted()
	assert.True(t, todo.Completed)
	assert.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)
}

// валидность
func TestTodoValidation(t *testing.T) {
	validTodo := &entity.Todo{
		Title:    "Valid Todo",
		ActiveAt: time.Now().Add(time.Hour * 24),
	}
	assert.NoError(t, validTodo.Validate())

	// Пустой заголовок
	emptyTitleTodo := &entity.Todo{
		Title:    "",
		ActiveAt: time.Now().Add(time.Hour * 24),
	}
	assert.Error(t, emptyTitleTodo.Validate())
	assert.Equal(t, errors.ErrTitleEmpty, emptyTitleTodo.Validate())

	// Длинна больше 200
	longTitleTodo := &entity.Todo{
		Title:    "ААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААА",
		ActiveAt: time.Now().Add(time.Hour * 24),
	}
	assert.Error(t, longTitleTodo.Validate())
	assert.Equal(t, errors.ErrTitleLengthExceeded, longTitleTodo.Validate())

	// Test case: Past activation date
	pastActivationTodo := &entity.Todo{
		Title:    "Past Activation Todo",
		ActiveAt: time.Now().Add(-time.Hour * 24),
	}
	assert.Error(t, pastActivationTodo.Validate())
	assert.Equal(t, errors.ErrDateNotCurrent, pastActivationTodo.Validate())
}
