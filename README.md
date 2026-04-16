# Task Service

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go `1.23+`
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

Если `postgres` уже запускался ранее со старой схемой, пересоздай volume:

```bash
docker compose down -v
docker compose up --build
```

Причина в том, что SQL-файл из `migrations/0001_create_tasks.up.sql` монтируется в `docker-entrypoint-initdb.d` и применяется только при инициализации пустого data volume.

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```

## API

Базовый префикс API:

```text
/api/v1
```

Основные маршруты:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`

## Решение тестового задания

Добавлена поддержка настройки периодичности для задач.
Реализация:
- Поля `scheduled_at`, `recurrence_type`, `recurrence_interval`, `recurrence_month_days`, `recurrence_specific_dates`, `next_generate_date`, `parent_task_id` были добавлены в таблицу `tasks`.
- Если задача создается с параметрами повторяемости (например, `recurrence_type = "daily"`), она рассматривается как "шаблон".
- В фоне (каждую минуту) работает воркер, который находит такие шаблоны с наступившим `next_generate_date` и генерирует новые инстансы задач. Инстанс имеет ссылку на родительскую задачу `parent_task_id`. У сгенерированной задачи поле `scheduled_at` установлено на дату выполнения.
- Поддерживаемые типы периодичности: `daily`, `monthly`, `specific_dates`, `even_days`, `odd_days`.
- При создании и обновлении задачи параметры повторяемости валидируются.
- Swagger/OpenAPI спецификация была обновлена, и добавлены примеры DTO для новых полей.
- Базовые Unit-тесты для логики расчета следующей даты вызова были добавлены в `internal/domain/task/task_test.go`.
