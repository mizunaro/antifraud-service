# antifraud-service

```
antifraud-service/
├── cmd/
│   └── server/          # Точка входа
├── internal/            # Весь приватный код
│   ├── app/             # Собираем зависимости, стартуем Кафку, БД
│   ├── config/          # Чтение env/yaml (чистый конфиг, без логики)
│   ├── domain/          # Структуры и интерфейсы
│   ├── transport/       # Внешний мир
│   │   ├── http/        # Роуты, хендлеры, сериализация JSON
│   │   └── kafka/       # Продюсер (отправка в очередь) и Консьюмер (чтение)
│   ├── service/         # Бизнес-логика: проверка URL, принятие решений
│   └── repository/      # Работа с данными: Postgres (SQL) и Redis (Cache)
├── migrations/          # SQL-файлы для базы
├── deploy/              # Dockerfile, docker-compose.yaml
├── go.mod
└── go.sum
```

