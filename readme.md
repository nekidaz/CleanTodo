
# CleanTodo App

### Вместо Swagger использовал [Postman](https://documenter.getpostman.com/view/24551580/2s9XxztCKj) - Импортируйте Коллекцию 

Простое приложение для управления задачами (Todo List).




## Запуск приложения

1. Убедитесь, что у вас установлены Docker и Docker Compose.

2. Клонируйте репозиторий:

   ```sh
   git clone https://github.com/nekidaz/CleanTodo.git
   cd CleanTodo
   ```

3. Запустите контейнеры с помощью Docker Compose:

   ```sh
   docker-compose up -d
   ```

4. Ваше приложение будет доступно по адресу: [http://localhost:8080](http://localhost:8080)

5. Чтобы остановить контейнеры, выполните:

   ```sh
   docker-compose down
   ```

## API Endpoints

### Получение всех задач

```
GET /api/tasks
```

### Создание новой задачи

```
POST /api/tasks
```

### Обновление задачи

```
PUT /api/tasks/:id
```

### Удаление задачи

```
DELETE /api/tasks/:id
```

### Пометить задачу как выполненную

```
PUT /api/tasks/:id/done
```

### Получение задач по статусу

```
GET /api/tasks/status/:status
```

Где `:id` - идентификатор задачи, `:status` - статус задачи (`done` или `active`).



