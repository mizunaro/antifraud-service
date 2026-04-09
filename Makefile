DB_URL = postgres://user:password@localhost:5432/antifraud?sslmode=disable
MIGRATIONS_PATH = ./migrations/

# Запуск в Docker
docker-up:
	docker compose -f deploy/docker-compose.yaml up -d

# Остановка Docker 
docker-down:
	docker-compose -f deploy/docker-compose.yaml down

# Создание нового файла миграции (например: make migrate-create name=init_db)
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

# Применение миграций
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

# Откат последней миграции
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

# Установка golang-migrate
migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
