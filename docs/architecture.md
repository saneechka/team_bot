# Архитектура Team Bot

## 🏛 Общая архитектура

Team Bot построен по принципам Clean Architecture с четким разделением слоев и ответственностей.

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
├─────────────────────────────────────────────────────────────┤
│  cmd/bot/main.go  │  internal/bot/bot.go                   │
│  - Инициализация  │  - Обработка updates                   │
│  - Graceful       │  - Маршрутизация команд                │
│    shutdown       │  - Контроль доступа                    │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                   Business Logic Layer                      │
├─────────────────────────────────────────────────────────────┤
│  internal/handler/                                          │
│  - AuthHandler: обработка авторизации                       │
│  - Бизнес-логика команд                                     │
│  - Валидация данных                                         │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                     Data Layer                              │
├─────────────────────────────────────────────────────────────┤
│  internal/model/     │  internal/repository/sqlrepo/        │
│  - User             │  - AuthRepository                     │
│  - Message          │  - CRUD операции                      │
│  - Business models  │  - SQL запросы                        │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                       │
├─────────────────────────────────────────────────────────────┤
│  config/            │  internal/repository/migrations/      │
│  - Configuration    │  - Database schema                    │
│  - Environment vars │  - Migration scripts                  │
└─────────────────────────────────────────────────────────────┘
```

## 🧩 Компоненты системы

### 1. Main Entry Point (`cmd/bot/main.go`)

**Ответственность:**
- Инициализация приложения
- Настройка зависимостей
- Graceful shutdown
- Управление lifecycle

**Основные функции:**
```go
func main() {
    // 1. Загрузка конфигурации
    // 2. Инициализация базы данных
    // 3. Создание зависимостей
    // 4. Запуск бота
    // 5. Обработка сигналов завершения
}
```

### 2. Bot Layer (`internal/bot/bot.go`)

**Ответственность:**
- Получение и обработка Telegram updates
- Маршрутизация команд к соответствующим обработчикам
- Контроль доступа на уровне сообщений
- Управление жизненным циклом бота

**Структура:**
```go
type Bot struct {
    api         *tgbotapi.BotAPI      // Telegram API клиент
    authRepo    *sqlrepo.AuthRepository // Репозиторий авторизации
    authHandler *handler.AuthHandler     // Обработчик авторизации
}
```

**Методы:**
- `Start(ctx)` - основной цикл обработки сообщений
- `handleMessage(ctx, update)` - маршрутизация сообщений
- `checkAdminAccess(ctx, userID, chatID)` - проверка прав доступа

### 3. Handler Layer (`internal/handler/`)

**Ответственность:**
- Обработка специфических команд
- Бизнес-логика операций
- Взаимодействие с репозиториями
- Формирование ответов пользователю

**AuthHandler структура:**
```go
type AuthHandler struct {
    bot  *tgbotapi.BotAPI
    repo *sqlrepo.AuthRepository
}
```

**Методы:**
- `HandleStart(ctx, update)` - обработка команды /start
- `HandleAdmin(ctx, update)` - обработка команды /admin
- `CheckAdminAccess(ctx, userID, chatID)` - проверка административных прав

### 4. Repository Layer (`internal/repository/sqlrepo/`)

**Ответственность:**
- Абстракция работы с базой данных
- CRUD операции для сущностей
- SQL запросы и их выполнение
- Обработка ошибок базы данных

**AuthRepository методы:**
```go
func (r *AuthRepository) SaveUser(ctx, user) error
func (r *AuthRepository) GetUserByID(ctx, id) (*User, error)
func (r *AuthRepository) GetUserByChatID(ctx, chatID) (*User, error)
func (r *AuthRepository) IsAdmin(ctx, userID) (bool, error)
func (r *AuthRepository) SetAdminStatus(ctx, userID, isAdmin) error
func (r *AuthRepository) GetUserByUsername(ctx, username) (*User, error)
```

### 5. Model Layer (`internal/model/`)

**Ответственность:**
- Определение бизнес-сущностей
- Структуры данных
- Бизнес-логика на уровне моделей
- Валидация данных

**Основные модели:**
```go
type User struct {
    ID          int64     `json:"id"`
    Username    string    `json:"username"`
    ChatID      int64     `json:"chat_id"`
    CreatedTime time.Time `json:"created_time"`
    IsAdmin     bool      `json:"is_admin"`
}

type Message struct {
    ID        int64     `json:"id"`
    ChatID    int64     `json:"chat_id"`
    Text      string    `json:"text"`
    UserID    int64     `json:"user_id"`
    Username  string    `json:"username"`
    Timestamp time.Time `json:"timestamp"`
    Type      string    `json:"type"`
}
```

### 6. Configuration Layer (`config/`)

**Ответственность:**
- Загрузка конфигурации из файлов
- Обработка переменных окружения
- Валидация конфигурации
- Предоставление настроек приложению

**Структура конфигурации:**
```go
type Config struct {
    Bot struct {
        Token string `yaml:"token"`
        Debug bool   `yaml:"debug"`
    } `yaml:"bot"`
    
    Database struct {
        Type string `yaml:"type"`
        Path string `yaml:"path"`
    } `yaml:"database"`
    
    TelegramAdmins struct {
        Usernames []string `yaml:"username"`
    } `yaml:"admins"`
}
```

## 🔄 Поток данных

### 1. Обработка входящего сообщения

```
Telegram → Bot.Start() → Bot.handleMessage() → Handler.Method() → Repository.Method() → Database
```

### 2. Проверка административных прав

```
User Request → Bot.checkAdminAccess() → AuthRepository.IsAdmin() → Database Query → Response
```

### 3. Регистрация нового пользователя

```
/start Command → AuthHandler.HandleStart() → AuthRepository.SaveUser() → Database Insert
```

## 🔒 Принципы безопасности архитектуры

### 1. Разделение ответственностей
- Каждый слой имеет четко определенную ответственность
- Бизнес-логика изолирована от инфраструктуры
- Данные валидируются на соответствующих уровнях

### 2. Dependency Injection
- Зависимости внедряются через конструкторы
- Легкое тестирование компонентов
- Слабая связанность между компонентами

### 3. Context Propagation
- Контекст передается через все слои
- Graceful cancellation для всех операций
- Таймауты и отмена операций

### 4. Error Handling
- Ошибки обрабатываются на каждом уровне
- Логирование ошибок без раскрытия внутренней информации
- Graceful degradation при сбоях

## 📈 Масштабируемость

### Горизонтальное масштабирование
- Stateless архитектура позволяет запускать несколько экземпляров
- База данных как единый источник истины
- Возможность добавления балансировщика нагрузки

### Вертикальное масштабирование
- Оптимизация SQL запросов
- Индексы для быстрого поиска
- Connection pooling для базы данных

### Расширяемость
- Легкое добавление новых handlers
- Модульная архитектура
- Возможность смены базы данных без изменения бизнес-логики
