# ImageProcessor

 ## ImageProcessor — сервис для фоновой обработки изображений с поддержкой очередей и хранения оригиналов, миниатюр и обработанных версий.

## Возможности

- Загрузка изображений через HTTP (`POST /upload`)
- Получение информации о изображении (`GET /image/{id}`)
- Удаление изображений (`DELETE /image/{id}`)
- Просмотр всех изображений (`GET /images`)
- Фоновая обработка через очередь (Kafka)
- Генерация:
  - уменьшенных версий (processed)
  - миниатюр (thumbnail)
  - водяного знака (watermark) при наличии
- Хранение:
  - оригинальные изображения (`data/uploads`)
  - обработанные (`data/processed`)
  - миниатюры (`data/thumbs`)
- Поддержка форматов: JPEG, PNG, GIF
- Простой веб-интерфейс для загрузки, просмотра и удаления изображений

## Технологии

- Go (Gin, sqlx, imaging)
- PostgreSQL (для хранения метаданных)
- Apache Kafka (фоновая обработка)
- HTML + JS + CSS (фронтенд)

## Структура проекта

/cmd/server - точка входа сервера
/internal/app - инициализация сущностей и запуск сервера
/internal/service - логика обработки изображений
/internal/storage - работа с файлами и БД
/internal/handlers - HTTP-эндпоинты и хэндлеры
/internal/queue_broker/kafka - инициализация и работа брокера сообщений
/web - фронтенд (HTML, JS, CSS)


## Быстрый старт

1. Настроить `.env` с параметрами БД, Kafka и путями к файлам.
2. Создать таблицу `images` в PostgreSQL:

```sql
CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    original_path TEXT NOT NULL,
    processed_path TEXT,
    thumbnail_path TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
либо при помощи миграций db/dumps

3. Запустить сервер:
go run cmd/server/main.go

4. Открыть веб-интерфейс: http://localhost:7575

5. Загрузить изображение и наблюдать его обработку.

Легкий и масштабируемый сервис для любых приложений, где нужно быстро обрабатывать изображения без блокировки пользователей.

Проект создан как демонстрация архитектуры многоуровневого Go-сервиса с PostgreSQL, KAFKA, REST API и минимальным UI без фреймворков

Автор
Разработчик: Vladimirmoscow84
Контакт: ccr1@yandex.ru
GitHub: github.com/Vladimirmoscow84


