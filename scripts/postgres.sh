#!/bin/bash

# Скрипт для управления PostgreSQL в Docker

case "$1" in
    start)
        echo "Запуск PostgreSQL..."
        docker-compose up -d postgres
        echo "PostgreSQL запущен. Ожидание готовности..."
        sleep 5
        docker-compose exec postgres pg_isready -U postgres
        ;;
    stop)
        echo "Остановка PostgreSQL..."
        docker-compose stop postgres
        ;;
    restart)
        echo "Перезапуск PostgreSQL..."
        docker-compose restart postgres
        ;;
    logs)
        echo "Логи PostgreSQL:"
        docker-compose logs postgres
        ;;
    connect)
        echo "Подключение к PostgreSQL..."
        docker-compose exec postgres psql -U postgres -d team_bot
        ;;
    status)
        echo "Статус PostgreSQL:"
        docker-compose ps postgres
        ;;
    migrate)
        echo "Запуск миграций..."
        # Здесь можно добавить команду для запуска миграций
        echo "Миграции выполнены автоматически при старте контейнера"
        ;;
    clean)
        echo "Очистка данных PostgreSQL..."
        docker-compose down postgres
        docker volume rm team_bot_postgres_data 2>/dev/null || true
        ;;
    *)
        echo "Использование: $0 {start|stop|restart|logs|connect|status|migrate|clean}"
        echo ""
        echo "Команды:"
        echo "  start   - Запустить PostgreSQL"
        echo "  stop    - Остановить PostgreSQL"
        echo "  restart - Перезапустить PostgreSQL"
        echo "  logs    - Показать логи PostgreSQL"
        echo "  connect - Подключиться к PostgreSQL"
        echo "  status  - Показать статус PostgreSQL"
        echo "  migrate - Запустить миграции"
        echo "  clean   - Удалить данные PostgreSQL"
        exit 1
        ;;
esac
