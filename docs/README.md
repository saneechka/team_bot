# Team Bot - Техническая документация

## Обзор проекта

Team Bot - это Telegram-бот для управления командой, написанный на Go. Бот предоставляет функции аутентификации, проверки прав администратора и базовой обработки команд.

## Структура проекта

```
team_bot/
├── cmd/
│   └── bot/
│       └── main.go              # Точка входа приложения
├── config/
│   ├── config.go               # Конфигурация приложения
│   └── config.yaml             # Файл конфигурации
├── internal/
│   ├── bot/
│   │   └── bot.go              # Основная логика бота
│   ├── handler/
│   │   └── admin_handler.go    # Обработчики команд
│   ├── model/
│   │   └── message.go          # Модели данных
│   └── repository/
│       ├── migrations/         # SQL миграции
│       └── sqlrepo/
│           └── auth.go         # Репозиторий для работы с БД
├── docs/                       # Документация
├── data/                       # База данных SQLite
└── go.mod                      # Зависимости Go
```

## Архитектура

Проект следует принципам Clean Architecture:

- **cmd/bot/main.go** - точка входа, инициализация зависимостей
- **internal/bot** - основная логика бота
- **internal/handler** - обработчики команд и бизнес-логика
- **internal/model** - модели данных
- **internal/repository** - слой доступа к данным
- **config** - конфигурация приложения

## Технологический стек

- **Go 1.24.4** - основной язык программирования
- **SQLite** - база данных
- **Telegram Bot API v5** - интеграция с Telegram
- **YAML** - формат конфигурации

## Зависимости

```go
require (
    github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
    github.com/mattn/go-sqlite3 v1.14.28
    gopkg.in/yaml.v3 v3.0.1
)
```
