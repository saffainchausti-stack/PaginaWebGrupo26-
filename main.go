package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	sqlc "ServidorTrabajoWeb/db/sqlc"
)

var queries *sqlc.Queries
var db *sql.DB

func main() {
	// Define el directorio que contiene los archivos estáticos.
	staticDir := "./static"

	// Crea un manejador (handler) de servidor de archivos.
	// http.Dir convierte la ruta del directorio en un sistema de archivos HTTP.
	// http.FileServer crea un manejador que sirve archivos desde ese sistema.
	// ¡Automáticamente sirve index.html para directorios!
	fileServer := http.FileServer(http.Dir(staticDir))

	// Para subir pagina web gratis: Netlify drop
	// Registra el manejador para que atienda todas las peticiones ("/").
	// Usamos http.Handle porque fileServer es un http.Handler.
	http.Handle("/", fileServer)

	// Define el puerto y muestra un mensaje
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	// Establece la conexion con la base de datos
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=recetas sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	queries = sqlc.New(db)

	// Inicia el servidor HTTP
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %s\n", err)
	}
}

func usuarioHandlerCRUD(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	var id int
	var err error
	switch r.Method {
	case http.MethodPost:
		if len(parts) != 1 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
	case http.MethodGet, http.MethodPut, http.MethodDelete:
		if len(parts) != 2 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		id, err = strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Invalid third parameter type", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
	fmt.Print(id)

	switch r.Method {
	case http.MethodGet:
		getUsuario(w, r, id)
	case http.MethodPost:
		postUsuario(w, r)
	case http.MethodPut:
		updateUsuario(w, r, id)
	case http.MethodDelete:
		deleteUsuario(w, r, id)
	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
}

func postUsuario(w http.ResponseWriter, r *http.Request) {
	var u sqlc.CreateUsuarioParams
	json.NewDecoder(r.Body).Decode(&u)
	uReturn, err := queries.CreateUsuario(context.Background(), u)
	if err != nil {
		log.Fatalf("Fallo al crear usuario: %v", err)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uReturn)
}

func listUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	usuarios, err := queries.ListUsuario(context.Background())
	if err != nil {
		log.Fatalf("fallo al listar usuarios: %v", err)
	}
	json.NewEncoder(w).Encode(usuarios)
}

func getUsuario(w http.ResponseWriter, _ *http.Request, id int) {
	u, err := queries.GetUsuario(context.Background(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
	w.WriteHeader(http.StatusOK)
}

func updateUsuario(w http.ResponseWriter, r *http.Request, id int) {
	var decoUsuario sqlc.UpdateUsuarioParams
	json.NewDecoder(r.Body).Decode(&decoUsuario)
	u := sqlc.UpdateUsuarioParams{
		IDUsuario: int32(id),
		Nombre:    decoUsuario.Nombre,
		Apellido:  decoUsuario.Apellido,
		Email:     decoUsuario.Email,
	}
	_, err := queries.UpdateUsuario(context.Background(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteProduct(w http.ResponseWriter, _ *http.Request, id int) {
	err := queries.DeleteUsuario(context.Background(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
