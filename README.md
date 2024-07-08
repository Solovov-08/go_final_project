# Планировщик задач

Планировщик задач — это веб-приложение на Go, предназначенное для управления задачами с возможностью повторяющихся событий. Приложение использует SQLite для хранения задач и предоставляет API для их управления.

### Требования

- Go версии 1.16 или выше
- SQLite3

### Клонирование репозитория

```bash
git clone https://github.com/ваш-репозиторий/task-scheduler.git
cd task-scheduler
go mod tidy
```

### Параметры окружения

Приложение использует два параметра окружения:

- `TODO_PORT`: Порт, на котором будет работать сервер. По умолчанию используется порт 7540.
- `TODO_DBFILE`: Путь к файлу базы данных SQLite. По умолчанию используется `./scheduler.db`.

## Инициализация базы данных

При первом запуске приложение автоматически создает файл базы данных и необходимые таблицы, если файл базы данных не существует.

### Запуск сервера

```bash
go run main.go
```

Сервер будет доступен по адресу http://localhost:7540

### Примеры использования API

#### Добавление задачи

```bash
curl -X POST "http://localhost:7540/api/task" -H "Content-Type: application/json" -d '{
"date": "20240201",
"title": "Подвести итог",
"comment": "Мой комментарий",
"repeat": "d 5"
}'
```

#### Получение всех задач

```bash
curl -X GET "http://localhost:7540/api/tasks"
```

#### Получение задачи по идентификатору

```bash
curl -X GET "http://localhost:7540/api/task?id=185"
```

#### Обновление задачи

```bash
curl -X PUT "http://localhost:7540/api/task" -H "Content-Type: application/json" -d '{
"id": 185,
"date": "20240201",
"title": "Обновленный заголовок",
"comment": "Обновленный комментарий",
"repeat": "d 10"
}'
```

#### Пометка задачи как выполненной

```bash
curl -X POST "http://localhost:7540/api/task/done?id=185"
```

#### Удаление задачи

```bash
curl -X DELETE "http://localhost:7540/api/task?id=185"
```

### Тестирование

Для запуска тестов выполните следующую команду:

```bash
go test ./...
```

Эта команда запустит все тесты в проекте и выведет результаты.

### Структура проекта

- `main.go`: Главный файл приложения, точка входа сервера.
- `database/`: Пакет для инициализации базы данных.
- `handlers/`: Пакет с обработчиками API запросов.
- `models/`: Пакет с моделями данных.
- `tasks/`: Пакет с логикой обработки задач и вычисления следующей даты выполнения.
- `web/`: Директория для статических файлов (HTML, CSS, JS).