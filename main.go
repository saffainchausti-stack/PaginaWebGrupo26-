package main

import (
	sqlc "ServidorTrabajoWeb/db/sqlc"
	"ServidorTrabajoWeb/views"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

	http.HandleFunc("GET /", indexHandler)
	// Manejador POST /receta
	http.HandleFunc("POST /receta", postReceta)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %s\n", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	recetas, err := queries.ListReceta(r.Context())
	if err != nil {
		log.Printf("Error al listar recetas: %v\n", err)
		http.Error(w, "Error al cargar las recetas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	componente := views.Index(recetas)
	if err := componente.Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar vista: %v\n", err)
	}
}

func postReceta(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(r.FormValue("nombre"))
	descripcion := strings.TrimSpace(r.FormValue("descripcion"))
	ingredientes := strings.TrimSpace(r.FormValue("ingredientes"))
	pasos := strings.TrimSpace(r.FormValue("pasos"))

	if nombre == "" || descripcion == "" {
		http.Error(w, "Nombre y descripción son requeridos", http.StatusBadRequest)
		return
	}

	// Usamos el usuario inicial con ID 1
	const defaultUserID = int32(1)

	_, err := queries.CreateReceta(r.Context(), sqlc.CreateRecetaParams{
		Nombre:       nombre,
		Descripcion:  descripcion,
		Pasos:        pasos,
		Ingredientes: ingredientes,
		IDUsuario:    defaultUserID,
	})
	if err != nil {
		log.Printf("Fallo al crear receta: %v\n", err)
		http.Error(w, "Fallo al crear receta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Patrón PRG: Redirigir a "/" con código 303 (See Other)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
