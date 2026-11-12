#APP_NAME := PaginaWebGrupo26-

#DB_URL := postgres://usuario:password@localhost:5432/mydb?sslmode=disable

.PHONY: run-app test clean generate-all build-app up-db down

generate-all:
	@rm -rf db/sqlc
	@sqlc generate
	@templ generate

build-app: generate-all
	@go build -o bin/api .

down:
	@docker compose down -v

run-app: down generate-all
	@echo "==> Iniciando contenedores (Base de datos y Web)..."
	@trap 'docker compose down -v' EXIT INT TERM; \
	docker compose up --build

# Tests de persistencia
test: down generate-all
	@docker compose up -d database --wait; \
	go test -v -count=1 ./tests; \
	status=$$?; \
	docker compose down -v; \
	exit $$status