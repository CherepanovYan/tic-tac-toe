# Tic-Tac-Toe

Веб-приложение «Крестики-нолики» с backend на Go, frontend на HTML/CSS/JavaScript и PostgreSQL для хранения пользователей и игр.

Приложение поддерживает:

* регистрацию пользователей;
* авторизацию через HTTP Basic Auth;
* создание игры с компьютером;
* создание игры между двумя игроками;
* просмотр доступных игр;
* подключение второго игрока;
* выполнение ходов;
* игру с компьютером на основе алгоритма Minimax;
* сохранение состояния игр и пользователей в PostgreSQL.

---

# Структура проекта

```text
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── api
│       └── main.go
├── config
│   └── config.go
├── docker-compose.yml
├── frontend
│   ├── index.html
│   ├── script.js
│   └── style.css
├── go.mod
├── go.sum
└── internal
    ├── di
    │   └── provider.go
    ├── domain
    │   ├── auth_service.go
    │   ├── domain_interfaces.go
    │   ├── game_service.go
    │   ├── model
    │   │   ├── errors.go
    │   │   ├── game.go
    │   │   ├── playingField.go
    │   │   └── user.go
    │   └── user_service.go
    ├── repository
    │   ├── mappers
    │   │   └── mappers.go
    │   ├── model
    │   │   └── model.go
    │   ├── postgres.go
    │   └── repository.go
    └── web
        ├── context.go
        ├── handler
        │   ├── auth_handler.go
        │   ├── game_handler.go
        │   ├── handler_interfaces.go
        │   └── user_handler.go
        ├── mappers
        │   └── mappers.go
        ├── middleware.go
        ├── model
        │   ├── errors.go
        │   └── model.go
        ├── router.go
        └── router_interfaces.go
```

---

# Настройка конфигурации

Перед запуском проекта необходимо создать `.env` файл.

В проекте находится пример конфигурации:

```text
./config/.env.example
```

Скопируйте его в директорию конфигурации:

```bash
cp ./config/.env.example ./.env
```

Проверьте значения переменных окружения перед запуском проекта.

В частности, должны быть настроены параметры подключения к PostgreSQL и порт HTTP-сервера.

---

# PostgreSQL

Проект использует PostgreSQL для хранения данных.

Основные таблицы:

* `users` — пользователи;
* `games` — игры.

Для запуска базы данных используется Docker Compose.

Запуск:

```bash
docker compose up -d
```

Проверить запущенные контейнеры:

```bash
docker compose ps
```

Остановить контейнеры:

```bash
docker compose down
```

---

# Запуск проекта

Для запуска backend и frontend используется:

```bash
make all
```

После запуска приложение будет доступно по адресу:

```text
http://localhost:8080
```

Порт зависит от значения, указанного в `.env`.

Frontend раздаётся backend-приложением.

> Если порт backend был изменён, необходимо изменить значение `API_BASE` в `frontend/script.js`.

Например:

```javascript
const API_BASE = "http://localhost:8080";
```

---

# Авторизация

Для авторизации используется HTTP Basic Authentication.

Frontend сохраняет логин и пароль в `localStorage` и при обращении к защищённым endpoints формирует заголовок:

```http
Authorization: Basic base64(login:password)
```

---

# API

## Регистрация

```http
POST /signup
```

Тело запроса:

```json
{
  "login": "user",
  "password": "password"
}
```

---

## Авторизация

```http
POST /login
```

Данные передаются в заголовке:

```http
Authorization: Basic base64(login:password)
```

Например:

```http
Authorization: Basic dXNlcjpwYXNzd29yZA==
```

При успешной авторизации сервер возвращает данные пользователя.

---

## Создание новой игры

Для создания игры с компьютером:

```http
GET /newGame?type=computer
```

Для создания игры между пользователями:

```http
GET /newGame?type=player
```

Также используется вариант:

```http
GET /newGame
```

который создаёт игру с компьютером.

Ответ:

```json
{
  "id": "uuid"
}
```

Endpoint требует авторизацию.

---

## Получение информации об игре

```http
GET /game/{UUID}
```

Пример:

```text
GET /game/550e8400-e29b-41d4-a716-446655440000
```

Ответ:

```json
{
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "playing_field": [
    [0, 0, 0],
    [0, 1, 0],
    [0, 0, 2]
  ],
  "status": "IN_PROGRESS",
  "user_x": "uuid",
  "user_o": "uuid",
  "current_player": "uuid"
}
```

---

## Получение доступных игр

```http
GET /games
```

Возвращает список игр между игроками, ожидающих второго игрока.

Endpoint требует авторизацию.

---

## Подключение к игре

```http
POST /game/{UUID}/join
```

Подключает текущего пользователя к существующей игре.

Endpoint требует авторизацию.

---

## Выполнение хода

```http
POST /game/{UUID}
```

Тело запроса:

```json
{
  "row": 0,
  "column": 1
}
```

`row` и `column` имеют значения от `0` до `2`.

Endpoint требует авторизацию.

После выполнения хода сервер возвращает актуальное состояние игры.

Для игры с компьютером после хода пользователя автоматически выполняется ход компьютера.

---

## Получение пользователя

```http
GET /user/{UUID}
```

Возвращает информацию о пользователе.

Endpoint требует авторизацию.

---

# Статусы игры

В зависимости от состояния игры сервер может возвращать следующие статусы:

```text
IN_PROGRESS
X_WON
O_WON
DRAW
```

---

# Архитектура проекта

Проект разделён на несколько слоёв.

## cmd/api

Точка входа приложения.

Содержит запуск HTTP-сервера и инициализацию зависимостей.

---

## internal/domain

Основной бизнес-слой приложения.

Содержит:

* игровую логику;
* создание игр;
* выполнение ходов;
* проверку победы;
* проверку ничьей;
* алгоритм Minimax;
* работу с пользователями;
* авторизацию;
* domain-модели;
* интерфейсы Repository.

---

## internal/repository

Слой работы с PostgreSQL.

Отвечает за:

* сохранение игр;
* получение игр по UUID;
* получение доступных игр;
* создание пользователей;
* получение пользователей по UUID;
* получение пользователей по логину.

Repository преобразует данные базы данных в domain-модели через mapper.

---

## internal/repository/model

DTO-модели, используемые для работы с данными PostgreSQL.

Например:

```text
GameDTO
UserDTO
BoardDTO
```

---

## internal/repository/mappers

Преобразование DTO-моделей Repository в domain-модели.

---

## internal/web

HTTP-слой приложения.

Содержит:

* HTTP handlers;
* HTTP middleware;
* регистрацию маршрутов;
* web-модели запросов и ответов;
* преобразование domain-моделей в web-модели.

---

## internal/web/handler

HTTP-обработчики приложения.

Основные handlers:

* `AuthHandler`
* `GameHandler`
* `UserHandler`

---

## internal/web/model

Модели HTTP-запросов и ответов.

Например:

```text
SignUpRequest
MoveRequest
GameResponse
ErrorResponse
```

---

## internal/web/mappers

Преобразование domain-моделей в модели HTTP-ответов.

---

## internal/di

Dependency Injection.

Для управления зависимостями используется Uber Fx.

---

# Frontend

Frontend находится в директории:

```text
frontend/
```

Состоит из:

```text
index.html
style.css
script.js
```

### index.html

Содержит интерфейс:

* авторизации;
* регистрации;
* главного меню;
* списка доступных игр;
* игрового поля.

### style.css

Отвечает за внешний вид приложения.

### script.js

Отвечает за:

* регистрацию;
* авторизацию;
* хранение credentials;
* Basic Auth;
* создание игр;
* получение доступных игр;
* подключение к игре;
* получение состояния игры;
* выполнение ходов;
* отображение игрового поля;
* взаимодействие с backend API.

---

# Используемые технологии

## Backend

* Go
* net/http
* PostgreSQL
* pgx
* Uber Fx
* UUID

## Frontend

* HTML5
* CSS3
* JavaScript
* Fetch API
* LocalStorage

## Инфраструктура

* Docker
* Docker Compose
* Make