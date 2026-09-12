APP_NAME := PaginaWebGrupo26-

DB_URL := postgres://usuario:password@localhost:5432/mydb?sslmode=disable

.PHONY: up down generate run test

up:
	docker compose up -d

down:
	docker compose down -v

generate:

	@sqlc generate

test:
	go test -v ./test


run: up generate test
	@air # Ejecuta las tareas de generación configuradas por el proyecto







