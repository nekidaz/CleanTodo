// models.go

package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Completed bool               `bson:"completed" json:"completed"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	ActiveAt  time.Time          `bson:"active_at" json:"active_at"`
}

func NewTodo(title string, activeAt time.Time) *Todo {
	return &Todo{
		Title:     title,
		Completed: false,
		ActiveAt:  activeAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *Todo) MarkAsCompleted() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

func (t *Todo) Validate() error {
	if t.Title == "" {
		return errors.New("Заголовок не может быть пустым")
	}

	t.Title = strings.ReplaceAll(t.Title, " ", "")
	if len(t.Title) > 200 {
		return errors.New("Длина заголовка не может превышать 200 символов")
	}
	today := time.Now().Truncate(24 * time.Hour) // Обрезаем время, оставляя только дату

	if t.ActiveAt.Before(today) {
		return errors.New("Дата должна быть актуальной и не раньше текущей даты")
	}
	return nil
}
