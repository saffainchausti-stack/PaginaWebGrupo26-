package main

import (
	sqlc "ServidorTrabajoWeb/db/sqlc"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
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
	http.HandleFunc("/usuario", usuarioHandlerCRUD)
	http.HandleFunc("/usuario/", usuarioHandlerCRUD)

	http.HandleFunc("/receta", recetaHandlerCRUD)
	http.HandleFunc("/receta/", recetaHandlerCRUD)

	// Define el puerto y muestra un mensaje
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	// Establece la conexion con la base de datos
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=recetas sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v\n", err)
	}
	defer db.Close()
	queries = sqlc.New(db)

	// Inicia el servidor HTTP
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %s\n", err)
	}
}

/////////////////////CRUDs de tabla usuario/////////////////////

func usuarioHandlerCRUD(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	switch r.Method {
	case http.MethodPost:
		if len(parts) != 1 {
			http.Error(w, "Ruta inválida", http.StatusBadRequest)
			return
		}
		postUsuario(w, r)
		return

	case http.MethodGet:
		if len(parts) == 1 {
			listUsuario(w, r)
			return
		}

	case http.MethodPut, http.MethodDelete:
		if len(parts) != 2 {
			http.Error(w, "Ruta inválida", http.StatusBadRequest)
			return
		}

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if len(parts) != 2 {
		http.Error(w, "Ruta inválida", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un entero positivo", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getUsuario(w, r, id)
	case http.MethodPut:
		updateUsuario(w, r, id)
	case http.MethodDelete:
		deleteUsuario(w, r, id)
	}
}

func postUsuario(w http.ResponseWriter, r *http.Request) {
	var u sqlc.CreateUsuarioParams
	json.NewDecoder(r.Body).Decode(&u)
	uReturn, err := queries.CreateUsuario(context.Background(), u)
	if err != nil {
		log.Printf("Fallo al crear usuario: %v\n", err)
		http.Error(w, "Fallo al crear usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uReturn)
}

func listUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	usuarios, err := queries.ListUsuario(context.Background())
	if err != nil {
		log.Printf("Fallo al listar usuarios: %v\n", err)
		http.Error(w, "Fallo al listar usuarios", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(usuarios)
}

func getUsuario(w http.ResponseWriter, _ *http.Request, id int) {
	u, err := queries.GetUsuario(context.Background(), int32(id))
	if err != nil {
		//si la base devuelve un error por no encontrar al usuario
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Usuario con ID %d no encontrado", id)
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		//si la base devuelve un error por fallo interno en la consulta
		log.Printf("Fallo al obtener usuario con ID %d: %v", id, err)
		http.Error(w, "Error interno del servidor",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
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
		log.Printf("Fallo al actualizar usuario con ID %d: %v\n", u.IDUsuario, err)
		http.Error(w, "Fallo al actualizar usuario", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteUsuario(w http.ResponseWriter, _ *http.Request, id int) {
	err := queries.DeleteUsuario(context.Background(), int32(id))
	if err != nil {
		log.Printf("Fallo al eliminar usuario con ID %d: %v\n", id, err)
		http.Error(w, "Fallo al eliminar usuario", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/////////////////////CRUDs de tabla receta/////////////////////

func recetaHandlerCRUD(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var id int
	var err error
	switch r.Method {
	case http.MethodPost:
		if len(parts) != 1 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
	case http.MethodPut, http.MethodDelete:
		if len(parts) != 2 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		id, err = strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Invalid third parameter type", http.StatusBadRequest)
			return
		}
	case http.MethodGet:
		if len(parts) == 2 {
			id, err = strconv.Atoi(parts[1])
			if err != nil {
				http.Error(w, "Invalid third parameter type", http.StatusBadRequest)
				return
			}
		} else {
			if len(parts) != 1 {
				http.Error(w, "Invalid URL", http.StatusBadRequest)
				return
			}
		}
	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
	fmt.Print(id)

	switch r.Method {
	case http.MethodGet:
		if len(parts) == 2 {
			getReceta(w, r, id)
		} else {
			listReceta(w, r)
		}
	case http.MethodPost:
		postReceta(w, r)
	case http.MethodPut:
		updateReceta(w, r, id)
	case http.MethodDelete:
		deleteReceta(w, r, id)
	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
}

func postReceta(w http.ResponseWriter, r *http.Request) {
	var u sqlc.CreateRecetaParams
	json.NewDecoder(r.Body).Decode(&u)
	uReturn, err := queries.CreateReceta(context.Background(), u)

	if err != nil {
		log.Fatalf("Fallo al crear receta: %v", err)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uReturn)
}

func listReceta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	recetas, err := queries.ListReceta(context.Background())
	if err != nil {
		log.Fatalf("fallo al listar recetas: %v", err)
	}
	json.NewEncoder(w).Encode(recetas)
}

func getReceta(w http.ResponseWriter, _ *http.Request, id int) {
	u, err := queries.GetReceta(context.Background(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
	w.WriteHeader(http.StatusOK)
}

func updateReceta(w http.ResponseWriter, r *http.Request, id int) {
	var decoReceta sqlc.UpdateRecetaParams
	json.NewDecoder(r.Body).Decode(&decoReceta)
	u := sqlc.UpdateRecetaParams{
		IDReceta:     int32(id),
		Nombre:       decoReceta.Nombre,
		Descripcion:  decoReceta.Descripcion,
		Pasos:        decoReceta.Pasos,
		Ingredientes: decoReceta.Ingredientes,
	}
	_, err := queries.UpdateReceta(context.Background(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteReceta(w http.ResponseWriter, _ *http.Request, id int) {
	err := queries.DeleteReceta(context.Background(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
