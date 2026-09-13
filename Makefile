#APP_NAME := PaginaWebGrupo26-

#DB_URL := postgres://usuario:password@localhost:5432/mydb?sslmode=disable

.PHONY: up down generate run test rmsqlc

up:
	docker compose up -d

down:
	docker compose down -v

generate:
	@sqlc generate

test:
	go test -v ./test

rmsqlc:
	if [ -d "db/sqlc" ]; then rm -r "db/sqlc"; fi

run: down up rmsqlc generate test down
	@air # Ejecuta las tareas de generación configuradas por el proyecto