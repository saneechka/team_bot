# Установка и настройка

## 📋 Системные требования

### Обязательные требования
- **Go**: версия 1.24.4 или выше
- **SQLite3**: для работы с базой данных
- **Git**: для клонирования репозитория

### Рекомендуемые требования
- **RAM**: минимум 256MB
- **Диск**: 100MB свободного места
- **ОС**: Linux, macOS, Windows

## 🚀 Установка

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd team_bot
```

### 2. Установка зависимостей

```bash
go mod download
```

### 3. Проверка зависимостей

```bash
go mod verify
```

## ⚙️ Настройка

### 1. Создание Telegram бота

1. Откройте Telegram и найдите [@BotFather](https://t.me/botfather)
2. Отправьте команду `/newbot`
3. Следуйте инструкциям для создания бота
4. Сохраните полученный токен

### 2. Настройка конфигурации

Создайте файл `config/config.yaml`:

```yaml
bot:
  token: "YOUR_BOT_TOKEN_HERE"  # Замените на токен вашего бота
  debug: false                  # true для отладки

database:
  type: "sqlite3"
  path: "./data/bot.db"

admins:
  username:
    - "your_telegram_username"  # Замените на ваш username без @
```

### 3. Настройка переменных окружения (опционально)

Для продакшена рекомендуется использовать переменные окружения:

```bash
export TELEGRAM_BOT_TOKEN="your_bot_token_here"
```

Создайте `.env` файл (добавлен в `.gitignore`):
```bash
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

## 🗄️ Настройка базы данных

### 1. Создание директории для данных

```bash
mkdir -p data
```

### 2. Применение миграций

Применить миграции пользователей:
```bash
sqlite3 ./data/bot.db < internal/repository/migrations/000001_users.up.sql
```

Применить миграции сообщений:
```bash
sqlite3 ./data/bot.db < internal/repository/migrations/000002_messages.up.sql
```

### 3. Назначение администратора (опционально)

Получите ваш Telegram User ID и выполните:
```sql
sqlite3 ./data/bot.db "UPDATE users SET is_admin = true WHERE id = YOUR_USER_ID;"
```

## 🏗️ Сборка

### Сборка для разработки

```bash
go build -o bin/team_bot ./cmd/bot
```

### Сборка для продакшена

```bash
CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/team_bot ./cmd/bot
```

### Кросс-компиляция

Для Linux:
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/team_bot-linux ./cmd/bot
```

Для Windows:
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o bin/team_bot.exe ./cmd/bot
```

## 🚀 Запуск

### Запуск в режиме разработки

```bash
go run ./cmd/bot
```

### Запуск собранного бинарника

```bash
./bin/team_bot
```

### Запуск с логированием

```bash
./bin/team_bot 2>&1 | tee logs/bot.log
```

## 🐳 Docker (опционально)

### Dockerfile

Создайте `Dockerfile`:
```dockerfile
FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o team_bot ./cmd/bot

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/

COPY --from=builder /app/team_bot .
COPY --from=builder /app/config ./config
COPY --from=builder /app/internal/repository/migrations ./migrations

CMD ["./team_bot"]
```

### Docker Compose

Создайте `docker-compose.yml`:
```yaml
version: '3.8'

services:
  team_bot:
    build: .
    environment:
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
    volumes:
      - ./data:/root/data
      - ./config:/root/config
    restart: unless-stopped
```

Запуск:
```bash
docker-compose up -d
```

## ✅ Проверка установки

### 1. Проверка сборки

```bash
go build ./cmd/bot
echo $?  # Должно вывести 0
```

### 2. Проверка конфигурации

```bash
go run ./cmd/bot &
PID=$!
sleep 5
kill $PID
```

### 3. Проверка базы данных

```bash
sqlite3 ./data/bot.db ".tables"
# Должно показать: messages users
```

### 4. Тестирование бота

1. Найдите вашего бота в Telegram
2. Отправьте команду `/start`
3. Проверьте, что бот отвечает

## 🔧 Отладка

### Включение debug режима

В `config/config.yaml`:
```yaml
bot:
  debug: true
```

### Просмотр логов

```bash
tail -f logs/bot.log
```

### Проверка базы данных

```bash
sqlite3 ./data/bot.db "SELECT * FROM users;"
```

## 🆘 Устранение проблем

### "sqlite3: command not found"

**Ubuntu/Debian:**
```bash
sudo apt-get install sqlite3
```

**macOS:**
```bash
brew install sqlite3
```

**Windows:**
Скачайте с [официального сайта SQLite](https://sqlite.org/download.html)

### "permission denied"

```bash
chmod +x bin/team_bot
```

### "database is locked"

```bash
# Остановите все экземпляры бота
pkill team_bot

# Проверьте блокировку
lsof ./data/bot.db
```

### Проблемы с сертификатами SSL

```bash
# Linux
sudo apt-get update && sudo apt-get install ca-certificates

# macOS
brew install ca-certificates
```

## 📝 Следующие шаги

После успешной установки:

1. Изучите [Конфигурацию](./configuration.md)
2. Ознакомьтесь с [API Reference](./api-reference.md)
3. Посмотрите [Примеры использования](./examples.md)
4. Настройте [Развертывание](./deployment.md)
