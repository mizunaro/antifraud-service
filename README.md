# Antifraud URL Analysis Service

Распределенная система для проверки URL на благонадежность.

## 🚀 Основной стек

- **Language:** Go 1.26+
- **Database:** PostgreSQL 16 (хранение результатов)
- **Message Broker:** Kafka (буферизация задач на проверку)
- **Cache:** Redis 8 (кэширование результатов анализа)
- **Observability:** Prometheus & Grafana (мониторинг метрик и состояния системы)
- **Deployment:** Docker & Docker Compose

## 🏗 Архитектура

Проект реализован в соответствии с принципами **Clean Architecture** и **Standard Go Project Layout**:

- `cmd/`: Точки входа для API-сервера и Воркера.
- `internal/app/`: Инициализация зависимостей и запуск приложения.
- `internal/service/`: Бизнес-логика (анализ ссылок, управление кэшем).
- `internal/repository/`: Слой работы с БД (PostgreSQL) и Кэшем (Redis).
- `internal/transport/`: Реализация протоколов HTTP (chi) и Kafka (segmentio/kafka-go).

### Data Flow

1. **API Server** принимает POST-запрос с URL, сохраняет запись в БД со статусом `pending` и отправляет событие в Kafka. Клиент получает `202 Accepted`.
2. **Worker** вычитывает сообщение из Kafka.
3. Проверяется наличие результата в **Redis**:
   - Если есть — статус обновляется мгновенно (**Cache Hit**).
   - Если нет — проводится имитация тяжелого анализа, результат сохраняется в Redis и БД.
4. Метрики обработки (количество, статусы, кэш-хиты) экспортируются в **Prometheus**.

## 🛠️ Запуск проекта

### Требования

- Docker & Docker Compose

### Быстрый старт

1. Склонируйте репозиторий.
2. Запустите инфраструктуру и сервисы:

   ```bash
   docker-compose -f deploy/docker-compose.yaml up -d --build
   # or
   make docker-build
   ```

3. Выполните миграции:
   ```bash
   migrate -path migrations/ -database "postgres://user:password@localhost:5432/antifraud?sslmode=disable" up
   # or
   make migrate-up
   ```

## 📊 Мониторинг и Инструменты

После запуска доступны следующие интерфейсы:

- **API Server:** `http://localhost:8080`
- **Prometheus:** `http://localhost:9090 `(просмотр метрик `antifraud_processed_urls_total`)
- **Grafana:** `http://localhost:3000 `(логин/пароль: `admin/admin`)
- **Kafka UI:** `http://localhost:8082 `(мониторинг топиков и сообщений)
- **Redis Commander:** `http://localhost:8081 `(просмотр кэша)

## 🧪 Пример запроса

```bash
curl -X POST http://localhost:8080/api/v1/check \
     -H "Content-Type: application/json" \
     -d '{"url": "https://crypto-casino.win"}'
```
