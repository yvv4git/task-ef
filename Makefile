# Environment file
ENV_FILE := .env
ifneq (,$(wildcard ${ENV_FILE}))
    include ${ENV_FILE}
endif

up:
	docker compose up -d

down:
	docker compose down

image_build:
	docker build -t yvv4docker/task-ef .

image_push:
	docker push yvv4docker/task-ef:latest


migrations_setup:
	goose -dir="./migrations" postgres "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable" up

migrations_reset:
	goose -dir="./migrations" postgres "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable" reset

swagger_gen:
	swag init --parseDependency --parseInternal -g internal/infrastructure/api/web.go
