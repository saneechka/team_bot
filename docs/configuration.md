# Конфигурация

## 📁 Структура конфигурации

Team Bot использует YAML-файл для конфигурации с возможностью переопределения через переменные окружения.

## 🎛️ Основной конфигурационный файл

### Расположение
```
config/config.yaml
```

### Полная структура конфигурации

```yaml
# Настройки Telegram бота
bot:
  token: "YOUR_BOT_TOKEN"     # Токен бота от @BotFather
  debug: false                # Режим отладки (true/false)

# Настройки базы данных
database:
  type: "sqlite3"            # Тип БД (пока только sqlite3)
  path: "./data/bot.db"      # Путь к файлу базы данных

# Администраторы системы
admins:
  username:                  # Список username администраторов
    - "admin_username1"      # без символа @
    - "admin_username2"
```

## 🔧 Детальное описание параметров

### Bot секция

| Параметр | Тип | Обязательный | Описание |
|----------|-----|--------------|----------|
| `token` | string | ✅ | Токен бота, полученный от [@BotFather](https://t.me/botfather) |
| `debug` | boolean | ❌ | Включает детальное логирование API запросов |

**Пример:**
```yaml
bot:
  token: "1234567890:ABCdefGhIJKlmNoPQRsTuVwXYz"
  debug: true  # для разработки
```

### Database секция

| Параметр | Тип | Обязательный | Описание |
|----------|-----|--------------|----------|
| `type` | string | ✅ | Тип базы данных (поддерживается только `sqlite3`) |
| `path` | string | ✅ | Относительный или абсолютный путь к файлу БД |

**Примеры:**
```yaml
database:
  type: "sqlite3"
  path: "./data/bot.db"          # Относительный путь
  
# или
database:
  type: "sqlite3"
  path: "/var/lib/team_bot/bot.db"  # Абсолютный путь
```

### Admins секция

| Параметр | Тип | Обязательный | Описание |
|----------|-----|--------------|----------|
| `username` | []string | ❌ | Массив Telegram username администраторов |

**Пример:**
```yaml
admins:
  username:
    - "john_doe"        # username без @
    - "admin_user"
    - "team_lead"
```

## 🌍 Переменные окружения

### Приоритет переменных окружения

Переменные окружения имеют **высший приоритет** и переопределяют значения из YAML файла.

### Поддерживаемые переменные

| Переменная | Переопределяет | Пример |
|------------|----------------|---------|
| `TELEGRAM_BOT_TOKEN` | `bot.token` | `export TELEGRAM_BOT_TOKEN="1234567890:ABC..."` |

### Примеры использования

**Linux/macOS:**
```bash
# Временно для текущей сессии
export TELEGRAM_BOT_TOKEN="your_token_here"

# Постоянно в ~/.bashrc или ~/.zshrc
echo 'export TELEGRAM_BOT_TOKEN="your_token_here"' >> ~/.bashrc
```

**Windows:**
```cmd
# Временно
set TELEGRAM_BOT_TOKEN=your_token_here

# Постоянно
setx TELEGRAM_BOT_TOKEN "your_token_here"
```

**Docker:**
```yaml
version: '3.8'
services:
  team_bot:
    image: team_bot:latest
    environment:
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
```

## 🔒 Безопасность конфигурации

### ❌ НЕ РЕКОМЕНДУЕТСЯ

```yaml
# НЕ храните токены в репозитории!
bot:
  token: "1234567890:ABCdefGhIJKlmNoPQRsTuVwXYz"  # ❌ ПЛОХО
```

### ✅ РЕКОМЕНДУЕТСЯ

**config/config.yaml (для разработки):**
```yaml
bot:
  token: ""  # Пустое значение, будет взято из переменной окружения
  debug: true

database:
  type: "sqlite3"
  path: "./data/bot.db"
```

**Использование:**
```bash
export TELEGRAM_BOT_TOKEN="your_real_token"
go run ./cmd/bot
```

### .env файл (для разработки)

Создайте `.env` файл (уже добавлен в `.gitignore`):
```bash
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

Загрузка перед запуском:
```bash
source .env
go run ./cmd/bot
```

## 🎯 Конфигурации для разных сред

### Разработка (development)

**config/config.yaml:**
```yaml
bot:
  token: ""              # Из переменной окружения
  debug: true            # Включить отладку

database:
  type: "sqlite3"
  path: "./data/dev.db"  # Отдельная БД для разработки

admins:
  username:
    - "your_dev_username"
```

### Тестирование (testing)

**config/config.test.yaml:**
```yaml
bot:
  token: ""
  debug: false

database:
  type: "sqlite3"
  path: ":memory:"      # In-memory база для тестов

admins:
  username:
    - "test_admin"
```

### Продакшен (production)

**config/config.yaml:**
```yaml
bot:
  token: ""              # Только из переменной окружения
  debug: false           # Отключить отладку

database:
  type: "sqlite3"
  path: "/var/lib/team_bot/bot.db"  # Системный путь

admins:
  username:
    - "production_admin"
    - "team_lead"
```

## 🔄 Загрузка конфигурации

### Код загрузки

```go
// config/config.go
func LoadConfig(path string) (*Config, error) {
    config := &Config{}
    
    // Загружаем из YAML файла
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    if err := yaml.Unmarshal(data, config); err != nil {
        return nil, err
    }
    
    // Переопределяем из переменных окружения
    if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
        config.Bot.Token = token
    }
    
    return config, nil
}
```

### Использование в main.go

```go
func main() {
    cfg, err := config.LoadConfig("config/config.yaml")
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }
    
    // Валидация обязательных параметров
    if cfg.Bot.Token == "" {
        log.Fatal("Bot token is required")
    }
}
```

## ✅ Валидация конфигурации

### Встроенная валидация

```go
func (c *Config) Validate() error {
    if c.Bot.Token == "" {
        return errors.New("bot token is required")
    }
    
    if c.Database.Type != "sqlite3" {
        return errors.New("only sqlite3 database is supported")
    }
    
    if c.Database.Path == "" {
        return errors.New("database path is required")
    }
    
    return nil
}
```

### Проверка конфигурации

```bash
go run ./cmd/config-check
```

## 🛠️ Дополнительные возможности

### Конфигурация через флаги командной строки

```go
import "flag"

var (
    configPath = flag.String("config", "config/config.yaml", "Path to config file")
    debug      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
    flag.Parse()
    
    cfg, err := config.LoadConfig(*configPath)
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }
    
    // Переопределение из флагов
    if *debug {
        cfg.Bot.Debug = true
    }
}
```

### Горячая перезагрузка конфигурации

```go
func watchConfig(configPath string) {
    // Мониторинг изменений файла конфигурации
    // Перезагрузка при изменениях
}
```

## 📝 Лучшие практики

### 1. Безопасность
- ✅ Никогда не коммитьте токены в репозиторий
- ✅ Используйте переменные окружения для чувствительных данных
- ✅ Используйте разные токены для разных сред

### 2. Организация
- ✅ Создавайте отдельные конфигурации для каждой среды
- ✅ Документируйте все параметры конфигурации
- ✅ Валидируйте конфигурацию при запуске

### 3. Мониторинг
- ✅ Логируйте загрузку конфигурации
- ✅ Проверяйте корректность параметров
- ✅ Отслеживайте изменения конфигурации

## 🆘 Устранение проблем

### "config file not found"
```bash
# Проверьте путь к файлу
ls -la config/config.yaml

# Создайте файл если его нет
cp config/config.yaml.example config/config.yaml
```

### "invalid yaml format"
```bash
# Проверьте синтаксис YAML
go run -c 'import yaml; yaml.safe_load(open("config/config.yaml"))'
```

### "bot token is empty"
```bash
# Проверьте переменную окружения
echo $TELEGRAM_BOT_TOKEN

# Установите если не установлена
export TELEGRAM_BOT_TOKEN="your_token"
```
