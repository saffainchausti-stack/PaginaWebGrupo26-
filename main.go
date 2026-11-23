package main

import (
	sqlc "ServidorTrabajoWeb/db/sqlc"
	"ServidorTrabajoWeb/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var queries *sqlc.Queries

func abrirDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=recetas sslmode=disable"
	}
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(25)
	return db, nil
}

func main() {
	db, err := abrirDB()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v\n", err)
	}
	defer db.Close()
	queries = sqlc.New(db)

	http.Handle("/CSS/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/", handlers.IndexHandler)

	http.HandleFunc("/about", handlers.AboutHandler)

	handlerRecetas := handlers.NewRecetasHandler(queries)
	http.Handle("/myRecetas", handlerRecetas)
	http.Handle("/myRecetas/", handlerRecetas)

	handlerLogin := handlers.NewLoginHandler(queries)
	http.Handle("/login", handlerLogin)

	handlerSignOn := handlers.NewSignOnHandler(queries)
	http.Handle("/signOn", handlerSignOn)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %s\n", err)
	}
}
