# Order Service

Микросервис для обработки и отображения данных о заказах, написанный на Go. Сервис получает данные о заказах из Kafka, сохраняет их в PostgreSQL и предоставляет REST API и веб-интерфейс для просмотра заказов.

## Функциональность

- Получение сообщений о заказах из Kafka
- Сохранение заказов в PostgreSQL
- In-memory кэширование заказов для быстрого доступа
- REST API для получения данных о заказе по ID
- Веб-интерфейс для поиска и просмотра заказов
- Автоматическое восстановление кэша из БД при перезапуске

## Архитектура

- **API Layer**: HTTP сервер на Go (Gorilla Mux)
- **Data Layer**: PostgreSQL для хранения данных
- **Message Broker**: Kafka для асинхронной обработки сообщений
- **Cache Layer**: In-memory кэш на основе map с RWMutex
- **Frontend**: Статический HTML/JS интерфейс

## Структура проекта
order-service/
├── cmd/
│ ├── producer/ # Тестовый продюсер для Kafka
│ └── server/ # Основной сервер
├── database/ # Работа с БД (PostgreSQL)
├── internal/
│ ├── cache/ # In-memory кэш
│ ├── handler/ # HTTP обработчики
│ ├── kafka/ # Kafka consumer
│ └── model/ # Модели данных
├── migrations/ # Миграции БД
├── static/ # Веб-интерфейс (HTML, JS, CSS)
├── docker-compose.yml # Docker Compose конфигурация
├── Dockerfile # Docker образ для сервиса
└── go.mod # Зависимости Go


## Требования

- Docker и Docker Compose
- Go 1.23+ (только для запуска producer из хоста)

## Быстрый старт

1. Клонируйте репозиторий:

```bash
git clone https://github.com/slinkeres/wbservice
```
2. Запустите все сервисы с помощью Docker Compose:

```bash
docker-compose up -d
```
3. Сервис будет доступен по адресу: http://localhost:8080

4. Отправьте тестовое сообщение в Kafka (из корневой директории)
```bash
go run cmd/producer/main.go
```
5. Откройте веб-интерфейс и введите ID заказа b563feb7b2b84b6test



