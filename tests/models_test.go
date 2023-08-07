package test

import (
	"RegionLabTZ/helpers"
	"RegionLabTZ/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateTodo(t *testing.T) {
	activationTime := time.Now().Add(time.Hour * 24) // Set activation for 24 hours from now
	title := "Test Todo"
	todo := models.NewTodo(title, activationTime)

	assert.Equal(t, title, todo.Title)
	assert.False(t, todo.Completed)
	assert.True(t, todo.ActiveAt.Equal(activationTime))
	assert.WithinDuration(t, time.Now(), todo.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)
}

func TestMarkAsCompleted(t *testing.T) {
	todo := models.NewTodo("Test Todo", time.Now())
	assert.False(t, todo.Completed)
	assert.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)

	todo.MarkAsCompleted()
	assert.True(t, todo.Completed)
	assert.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)
}

// валидность
func TestTodoValidation(t *testing.T) {
	validTodo := &models.Todo{
		Title:    "Valid Todo",
		ActiveAt: time.Now().Add(time.Hour * 24),
	}
	assert.NoError(t, validTodo.Validate())

	// Пустой заголовок
	emptyTitleTodo := &models.Todo{
		Title:    "",
		ActiveAt: time.Now().Add(time.Hour * 24),
	}
	assert.Error(t, emptyTitleTodo.Validate())
	assert.Equal(t, helpers.ErrTitleEmpty, emptyTitleTodo.Validate())

	// Длинна больше 200
	longTitleTodo := &models.Todo{
		Title:    "ААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААААА",
		ActiveAt: time.Now().Add(time.Hour * 24),
	}
	assert.Error(t, longTitleTodo.Validate())
	assert.Equal(t, helpers.ErrTitleLengthExceeded, longTitleTodo.Validate())

	// Test case: Past activation date
	pastActivationTodo := &models.Todo{
		Title:    "Past Activation Todo",
		ActiveAt: time.Now().Add(-time.Hour * 24),
	}
	assert.Error(t, pastActivationTodo.Validate())
	assert.Equal(t, helpers.ErrDateNotCurrent, pastActivationTodo.Validate())
}
