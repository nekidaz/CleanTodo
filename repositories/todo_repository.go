// repositories.swaggerRouter

package repositories

import (
	"RegionLabTZ/helpers"
	"RegionLabTZ/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TodoRepository interface {
	CreateNewTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error)
	UpdateTodo(ctx context.Context, id primitive.ObjectID, todo *models.Todo) (*models.Todo, error)
	DeleteTodo(ctx context.Context, id primitive.ObjectID) error
	MarkAsCompleted(ctx context.Context, id primitive.ObjectID) error
	GetTasksByStatus(ctx context.Context, status string) ([]*models.Todo, error)
	GetAllTasks(ctx context.Context) ([]*models.Todo, error)
	GetTaskByID(ctx context.Context, id primitive.ObjectID) (*models.Todo, error)
	Close() error
}

type repository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewRepository(connectionString, dbName, collectionName string) (TodoRepository, error) {
	// Подключение к MongoDB
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Пингуем сервер MongoDB
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)
	collection := database.Collection(collectionName)

	return &repository{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (r *repository) CreateNewTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error) {
	if err := todo.Validate(); err != nil {
		return nil, err
	}

	// Проверка уникальности записи по полям title и activeAt
	filter := bson.D{
		{"title", todo.Title},
		{"active_at", todo.ActiveAt},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, helpers.ErrTodoExists
	}

	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, helpers.ErrFailedToGetRecordID
	}

	todo.ID = insertedID
	return todo, nil
}

func (r *repository) UpdateTodo(ctx context.Context, id primitive.ObjectID, todo *models.Todo) (*models.Todo, error) {
	if err := todo.Validate(); err != nil {
		return nil, err
	}

	// Проверка существования задачи по ID
	existingTodo, err := r.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверка уникальности записи по полям title и activeAt (за исключением текущей задачи)
	filter := bson.D{
		{"title", todo.Title},
		{"active_at", todo.ActiveAt},
		{"_id", bson.D{{"$ne", id}}},
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, helpers.ErrNotFound
	}

	todo.ID = id
	todo.CreatedAt = existingTodo.CreatedAt
	todo.UpdatedAt = time.Now()

	update := bson.M{
		"$set": todo,
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *repository) DeleteTodo(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) MarkAsCompleted(ctx context.Context, id primitive.ObjectID) error {
	existingTodo, err := r.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}

	// Если задача уже выполнена, ничего не делаем
	if existingTodo.Completed {
		return nil
	}

	// Помечаем задачу как выполненную
	update := bson.M{
		"$set": bson.M{"completed": true, "updated_at": time.Now()},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetTasksByStatus(ctx context.Context, status string) ([]*models.Todo, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var filter bson.M
	if status == "done" {
		filter = bson.M{"completed": true}
	} else {
		// Получить задачи, которые не завершены и имеют activeAt <= today
		filter = bson.M{"completed": false, "active_at": bson.M{"$lte": today}}
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []*models.Todo
	if err = cursor.All(ctx, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

// Вспомогательный метод для поиска задачи по ID
func (r *repository) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*models.Todo, error) {
	var todo models.Todo
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, helpers.ErrNotFound
		}
		return nil, err
	}
	return &todo, nil
}

func (r *repository) GetAllTasks(ctx context.Context) ([]*models.Todo, error) {
	filter := bson.M{}

	// Получаем список всех задач
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"active_at": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []*models.Todo
	for cursor.Next(ctx) {
		var todo models.Todo
		if err := cursor.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, nil
}

func (r *repository) Close() error {
	if r.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := r.client.Disconnect(ctx); err != nil {
			return err
		}
	}
	return nil
}
