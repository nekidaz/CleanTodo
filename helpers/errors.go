package helpers

import "errors"

var (
	ErrTodoExists          = errors.New("Задача с таким заголовком и датой уже существует")
	ErrNotFound            = errors.New("Запись не найдена")
	ErrFailedToGetRecordID = errors.New("Не удалось получить идентификатор записи")
	ErrTitleEmpty          = errors.New("Заголовок не может быть пустым")
	ErrTitleLengthExceeded = errors.New("Длина заголовка не может превышать 200 символов")
	ErrDateNotCurrent      = errors.New("Дата должна быть актуальной и не раньше текущей даты")
	ErrParseActiveAt       = errors.New("Не удалось преобразовать ActiveAt")
	ErrInvalidID           = errors.New("Неверный ID")
	ErrTaskNotFound        = errors.New("Задача не найдена")
)
