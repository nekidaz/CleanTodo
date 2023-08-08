package entity

import (
	"github.com/nekidaz/todolist/pkg/errors"
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
		return errors.ErrTitleEmpty
	}

	//для проверки длинны
	subTitile := strings.ReplaceAll(t.Title, " ", "")

	if len(subTitile) > 200 {
		return errors.ErrTitleLengthExceeded
	}

	//чтобы мог создавать задачи на сегодня
	now := time.Now().UTC().Truncate(24 * time.Hour)

	// проверка времени
	if t.ActiveAt.UTC().Truncate(24 * time.Hour).Before(now) {
		return errors.ErrDateNotCurrent
	}

	return nil
}
