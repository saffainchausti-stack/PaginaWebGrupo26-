#APP_NAME := PaginaWebGrupo26-

#DB_URL := postgres://usuario:password@localhost:5432/mydb?sslmode=disable

.PHONY: up down generate build exetest test rmsqlc

up:
	@docker compose up -d --wait

down:
	@docker compose down -v

generate:
	@sqlc generate

build:
	@go build ./...

rmsqlc:
	@rm -rf db/sqlc

test: rmsqlc down generate build up 
	@go test -v -count=1 ./test; \
	status=$$?; \
	docker compose down -v; \
	exit $$status
