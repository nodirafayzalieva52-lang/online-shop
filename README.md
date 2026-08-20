# Production-Ready E-Commerce Backend (Go Clean Architecture)

[![CI](https://github.com/nodirafayzalieva52-lang/online-shop/actions/workflows/ci.yml/badge.svg)](https://github.com/nodirafayzalieva52-lang/online-shop/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue?style=flat&logo=postgresql)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Высокопроизводительный, надежный и масштабируемый RESTful API бэкенд интернет-магазина (маркетплейса) на языке **Go**, построенный с соблюдением принципов **Clean Architecture** (Чистая Архитектура), строгого разделения слоев, транзакционной целостности бизнес-логики и ролевого разграничения доступа (RBAC).

---

## 📋 Содержание

- [Технологический стек](#-технологический-стек)
- [Архитектура проекта](#-архитектура-проекта)
- [Структура директорий](#-структура-директорий)
- [Быстрый запуск](#-быстрый-запуск)
  - [Запуск через Docker Compose](#1-запуск-через-docker-compose-рекомендуется)
  - [Локальный запуск](#2-локальный-запуск-без-docker)
- [Переменные окружения](#-переменные-окружения)
- [Схема базы данных](#-схема-базы-данных)
- [Спецификация API (Endpoints)](#-спецификация-api-endpoints)
- [Примеры cURL запросов](#-примеры-curl-запросов)
- [Тестирование и CI/CD](#-тестирование-и-cicd)
- [Выполненный рефакторинг и исправления](#-выполненный-рефакторинг-и-исправления)

---

## 🛠 Технологический стек

- **Язык**: Go (1.22+)
- **HTTP Сервер и Роутинг**: Стандартный `net/http` (Go 1.22+ `ServeMux` с поддержкой методов и параметров путей)
- **База данных**: PostgreSQL 16
- **Драйвер БД & Пул соединений**: `pgx/v5` (`pgxpool`)
- **Миграции**: `golang-migrate`
- **Аутентификация**: JWT (JSON Web Tokens) via `golang-jwt/jwt/v5`
- **Логирование**: Структурированное логирование `go.uber.org/zap` с контекстным `X-Request-ID`
- **Конфигурация**: `cleanenv` (поддержка `.env` и переменных окружения)
- **Контейнеризация**: Docker (Multi-stage build) & Docker Compose
- **Тестирование**: `testing`, `net/http/httptest`, `github.com/stretchr/testify`

---

## 📐 Архитектура проекта

Проект организован в соответствии со стандартами **Clean Architecture**:

1. **Cmd Layer (`cmd/api/main.go`)**: Точка входа в приложение. Отвечает за инициализацию конфигурации, логгера, пула соединений БД, внедрение зависимостей (Dependency Injection) и Graceful Shutdown.
2. **Delivery Layer (`internal/delivery/http`)**:
   - `handler`: Переводит HTTP-запросы в вызовы сервисного слоя и возвращает единый стандартизированный JSON-ответ.
   - `middleware`: Реализует сквозную логику — аутентификацию (JWT), авторизацию по ролям (RBAC), логирование HTTP-запросов и CORS.
   - `dto`: Определение DTO (Data Transfer Objects) для входящих запросов и исходящих ответов.
   - `router`: Централизованный роутер с явным монтированием эндпоинтов и middleware.
3. **Service Layer (`internal/service`)**: Бизнес-логика приложения. Валидация входных данных, проверки прав доступа, координация работы репозиториев.
4. **Repository Layer (`internal/repository/postgres`)**: Слой работы с базой данных PostgreSQL. Содержит SQL-запросы, атомарные транзакции (ACID) и сканирование строк.
5. **Domain Layer (`internal/domain`)**: Ядро приложения — базовые доменные сущности (`User`, `Store`, `Category`, `Product`, `Order`, `OrderItem`).
6. **Infrastructure / PKG (`pkg/`)**: Вспомогательные пакеты общего назначения (`jwt`, `logger`, `hash`, `errors`).

---

## 📂 Структура директорий

```text
.
├── cmd/
│   └── api/
│       └── main.go               # Точка входа в приложение
├── internal/
│   ├── config/                   # Загрузка и валидация конфигурации
│   ├── delivery/
│   │   └── http/
│   │       ├── dto/              # DTO модели запросов и ответов
│   │       ├── handler/          # HTTP хендлеры
│   │       ├── middleware/       # JWT Auth, RequireRole, RequestLogger, CORS
│   │       └── router/           # Конфигурация роутов и HTTP сервера
│   ├── domain/                   # Доменные модели
│   ├── repository/
│   │   ├── postgres/             # Реализация репозиториев для PostgreSQL
│   │   └── repository.go         # Интерфейсы репозиториев
│   └── service/                  # Бизнес-логика и сценарии использования
├── migrations/                   # SQL скрипты миграций golang-migrate
│   ├── 000001_init.up.sql
│   ├── 000001_init.down.sql
│   ├── 000002_add_indexes.up.sql
│   └── 000002_add_indexes.down.sql
├── pkg/
│   ├── errors/                   # Доменные ошибки приложения
│   ├── hash/                     # Хеширование паролей (bcrypt)
│   ├── jwt/                      # Генерация и валидация JWT токенов
│   └── logger/                   # Настройка Zap логгера
├── config.env.example            # Пример файла конфигурации
├── docker-compose.yml            # Docker Compose манифест
├── Dockerfile                    # Оптимизированный multi-stage Dockerfile
├── Makefile                      # Скрипты сборки, запуска и миграций
└── README.md                     # Документация проекта
```

---

## 🚀 Быстрый запуск

### 1. Запуск через Docker Compose (Рекомендуется)

Запустит PostgreSQL 16 и API-сервис в изолированных контейнерах:

```bash
# Сборка и запуск контейнеров
make docker-up

# Проверка статуса сервиса
curl http://localhost:8080/health

# Остановка контейнеров
make docker-down
```

### 2. Локальный запуск (Без Docker)

**Требования**: Установленный Go 1.22+ и запущенный PostgreSQL.

```bash
# 1. Скопировать пример конфигурации
cp config.env.example config.env

# 2. Применить миграции базы данных
make migrate-up DB_URL="postgres://postgres:postgres_password@localhost:5432/shop_db?sslmode=disable"

# 3. Запустить приложение
make run
```

---

## ⚙️ Переменные окружения

Все настройки можно переопределить через переменные окружения или файл `config.env`:

| Переменная | Описание | Значение по умолчанию |
| :--- | :--- | :--- |
| `HTTP_PORT` | Порт HTTP-сервера | `:8080` |
| `DB_HOST` | Хост PostgreSQL | `localhost` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `postgres_password` |
| `DB_NAME` | Название базы данных | `shop_db` |
| `DB_SSLMODE` | Режим SSL для БД | `disable` |
| `JWT_SECRET` | Секретный ключ подписи JWT | `super-secret-key-at-least-32-chars` |
| `JWT_TTL` | Время жизни токена | `24h` |

---

## 🗄 Схема базы данных

```mermaid
erDiagram
    users ||--o{ stores : "owns (seller)"
    users ||--o{ orders : "places (customer)"
    stores ||--o{ products : "contains"
    categories ||--o{ products : "classifies"
    orders ||--|{ order_items : "contains"
    products ||--o{ order_items : "referenced in"

    users {
        bigserial id PK
        varchar email UK
        varchar password_hash
        varchar role
        timestamp created_at
        timestamp updated_at
    }

    stores {
        bigserial id PK
        bigint seller_id FK
        varchar name
        text description
        timestamp created_at
    }

    categories {
        bigserial id PK
        varchar name UK
    }

    products {
        bigserial id PK
        bigint store_id FK
        bigint category_id FK
        varchar name
        text description
        numeric price
        integer stock
        timestamp created_at
        timestamp updated_at
    }

    orders {
        bigserial id PK
        bigint customer_id FK
        bigint store_id FK
        numeric total_price
        varchar status
        timestamp created_at
    }

    order_items {
        bigserial id PK
        bigint order_id FK
        bigint product_id FK
        bigint store_id FK
        integer quantity
        numeric price
    }
```

---

## 📡 Спецификация API (Endpoints)

### Формат ответов об ошибках

Все ошибки API возвращаются в едином формате:
```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "product stock cannot be negative"
  }
}
```

### Таблица эндпоинтов

| Метод | URL | Роль / Авторизация | Описание | Коды ответов |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/health` | Публичный | Проверка доступности БД и сервиса | `200`, `503` |
| `POST` | `/register` | Публичный | Регистрация нового пользователя | `201`, `400`, `409` |
| `POST` | `/login` | Публичный | Вход в систему и получение JWT токена | `200`, `400`, `401` |
| `GET` | `/me` | Auth | Профиль текущего пользователя | `200`, `401`, `404` |
| `PATCH` | `/me` | Auth | Обновление email / пароля пользователя | `200`, `400`, `401`, `409` |
| `POST` | `/stores` | Auth (Customer/Seller) | Создание магазина (Customer становится Seller) | `201`, `400`, `401` |
| `GET` | `/stores/{id}` | Публичный | Получение информации о магазине | `200`, `400`, `404` |
| `GET` | `/stores/seller` | Auth (Seller/Admin) | Получение магазина продавца | `200`, `401`, `404` |
| `PATCH` | `/stores/{id}` | Store Owner / Admin | Обновление названия/описания магазина | `200`, `400`, `401`, `403`, `404` |
| `DELETE`| `/stores/{id}` | Store Owner / Admin | Удаление магазина | `200`, `401`, `403`, `404` |
| `GET` | `/categories` | Публичный | Получение списка категорий | `200` |
| `GET` | `/categories/{id}` | Публичный | Получение категории по ID | `200`, `400`, `404` |
| `POST` | `/categories` | Admin | Создание новой категории | `201`, `400`, `401`, `403` |
| `PATCH` | `/categories/{id}` | Admin | Изменение категории | `200`, `400`, `401`, `403`, `404` |
| `DELETE`| `/categories/{id}` | Admin | Удаление категории | `200`, `401`, `403`, `404` |
| `GET` | `/products` | Публичный | Список всех товаров (`?limit=10&offset=0`) | `200` |
| `GET` | `/products/{id}` | Публичный | Получение товара по ID | `200`, `400`, `404` |
| `GET` | `/stores/{store_id}/products` | Публичный | Товары магазина (`?limit=10&offset=0`) | `200`, `400` |
| `POST` | `/products` | Store Owner / Admin | Добавление нового товара | `201`, `400`, `401`, `403`, `404` |
| `PATCH` | `/products/{id}` | Store Owner / Admin | Редактирование товара | `200`, `400`, `401`, `403`, `404` |
| `DELETE`| `/products/{id}` | Store Owner / Admin | Удаление товара | `200`, `401`, `403`, `404` |
| `POST` | `/orders` | Auth (Customer) | Оформление заказа (атомарное списание) | `201`, `400`, `401`, `409`, `422` |
| `GET` | `/orders` | Auth | Заказы пользователя | `200`, `401` |
| `GET` | `/orders/{id}` | Customer / Store Owner / Admin | Детали заказа | `200`, `401`, `403`, `404` |
| `PATCH` | `/orders/{id}/status` | Store Owner / Admin | Изменение статуса заказа (с возвратом склада при `cancelled`) | `200`, `400`, `401`, `403`, `404`, `409` |

---

## 💻 Примеры cURL запросов

### 1. Регистрация нового пользователя
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "seller@example.com",
    "password": "password123"
  }'
```

### 2. Вход и получение JWT токена
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "seller@example.com",
    "password": "password123"
  }'
```
*Ответ*: `{"token": "YOUR_JWT_TOKEN_HERE"}`

### 3. Создание магазина
```bash
curl -X POST http://localhost:8080/stores \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech Galaxy Store",
    "description": "Электроника и гаджеты по лучшим ценам"
  }'
```

### 4. Добавление товара в магазин
```bash
curl -X POST http://localhost:8080/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "category_id": 1,
    "name": "Флагманский Смартфон 2026",
    "description": "AMOLED 120Hz, 12GB RAM, 512GB",
    "price": 89900.00,
    "stock": 15
  }'
```

### 5. Оформление заказа
```bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'
```

---

## 🧪 Тестирование и CI/CD

Проект содержит модульные (Unit) и интеграционные тесты HTTP-слоя.

```bash
# Запуск всех тестов
make test

# Проверка кода линтером и go vet
make lint
```

В репозиторий встроен **GitHub Actions Workflow** (`.github/workflows/ci.yml`), который автоматически проверяет сборку, запускает `go vet` и выполняет все тесты при каждом Push и Pull Request.

---

## 🛠 Выполненный рефакторинг и исправления

1. **Архитектурный рефакторинг и структура**:
   - Перенесены все хендлеры из устаревшего пакета `/handler` в `internal/delivery/http/handler`.
   - Выделен изолированный слой роутинга в `internal/delivery/http/router`.
   - Точка входа перенесена в `cmd/api/main.go`.
   - Устранены дублирующиеся пакеты `usecase` и `handler`.

2. **Безопасность и Бизнес-логика**:
   - **Уязвимость с ценой заказа полностью устранена**: Итоговая стоимость заказа и цены позиций рассчитываются исключительно на стороне бэкенда из базы данных под транзакционной блокировкой `FOR UPDATE`. Переданные клиентом цены игнорируются.
   - **Атомарное списание склада**: Списание количества товара `stock` происходит в единой базе данных транзакции `pgx.Tx` с проверкой остатка `(stock >= quantity)`.
   - **Возврат склада при отмене заказа**: При переводе заказа в статус `cancelled` складские остатки автоматически возвращаются в базу данных.
   - **Валидация принадлежности товара магазину**: Добавлена проверка того, что товар действительно принадлежит указанному в заказе `store_id`.
   - **Проверка прав владельцев (Ownership)**: Исключена возможность редактирования/удаления чужих магазинов и товаров.

3. **Инфраструктурные улучшения**:
   - Реализован Graceful Shutdown сервера и закрытие пула БД по сигналам `SIGINT`/`SIGTERM`.
   - Добавлены `RequestLogger` с генерацией `X-Request-ID` и `CORS` middleware.
   - Стандартизирован формат HTTP ошибок (`{"error": {"code": "...", "message": "..."}}`).
   - Добавлены B-Tree индексы PostgreSQL для частых выборок (`products.store_id`, `orders.customer_id`, и др.).
   - Добавлены `Dockerfile` (Multi-stage build), `docker-compose.yml`, `Makefile` и `config.env.example`.
