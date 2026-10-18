package main

import (
	sqlc "ServidorTrabajoWeb/db/sqlc"
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
	// Establece la conexion con la base de datos
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=recetas sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v\n", err)
	}
	defer db.Close()
	queries = sqlc.New(db)

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

	http.HandleFunc("/comentario", comentarioHandlerCRUD)
	http.HandleFunc("/comentario/", comentarioHandlerCRUD)

	// Define el puerto y muestra un mensaje
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	// Inicia el servidor HTTP
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %s\n", err)
	}
}

///////////////////// DTOs para las tablas /////////////////////

type UsuarioInput struct {
	Nombre      string `json:"nombre"`
	Apellido    string `json:"apellido"`
	Email       string `json:"email"`
	Contrasenia string `json:"contrasenia"`
}

type RecetaInput struct {
	Nombre       string `json:"nombre"`
	Descripcion  string `json:"descripcion"`
	Ingredientes string `json:"ingredientes"`
	Pasos        string `json:"pasos"`
	IDUsuario    int32  `json:"id_usuario"`
}

type ComentarioInput struct {
	IDUsuario   int32  `json:"id_usuario"`
	IDReceta    int32  `json:"id_receta"`
	Descripcion string `json:"descripcion"`
	Puntuacion  int16  `json:"puntuacion"`
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
	var input UsuarioInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Nombre) == "" || strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Contrasenia) == "" {
		http.Error(w, "Nombre, email y contraseña son obligatorios", http.StatusBadRequest)
		return
	}

	uReturn, err := queries.CreateUsuario(r.Context(), sqlc.CreateUsuarioParams{
		Nombre:      input.Nombre,
		Apellido:    sql.NullString{String: input.Apellido, Valid: input.Apellido != ""},
		Email:       input.Email,
		Contrasenia: input.Contrasenia,
	})
	if err != nil {
		log.Printf("Fallo al crear usuario: %v\n", err)
		http.Error(w, "Fallo al crear usuario: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uReturn)
}

func listUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	usuarios, err := queries.ListUsuario(r.Context())
	if err != nil {
		log.Printf("Fallo al listar usuarios: %v\n", err)
		http.Error(w, "Fallo al listar usuarios", http.StatusInternalServerError)
		return
	}
	if usuarios == nil {
		usuarios = []sqlc.ListUsuarioRow{}
	}
	json.NewEncoder(w).Encode(usuarios)
}

func getUsuario(w http.ResponseWriter, r *http.Request, id int) {
	u, err := queries.GetUsuario(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func updateUsuario(w http.ResponseWriter, r *http.Request, id int) {
	var input UsuarioInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Nombre) == "" || strings.TrimSpace(input.Email) == "" {
		http.Error(w, "Nombre y email no pueden estar vacíos", http.StatusBadRequest)
		return
	}

	actualizado, err := queries.UpdateUsuario(r.Context(), sqlc.UpdateUsuarioParams{
		IDUsuario: int32(id),
		Nombre:    input.Nombre,
		Apellido:  sql.NullString{String: input.Apellido, Valid: input.Apellido != ""},
		Email:     input.Email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Fallo al actualizar usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actualizado)
}

func deleteUsuario(w http.ResponseWriter, r *http.Request, id int) {
	err := queries.DeleteUsuario(r.Context(), int32(id))
	if err != nil {
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
			http.Error(w, "Ruta inválida", http.StatusBadRequest)
			return
		}
		postReceta(w, r)
		return
	case http.MethodGet:
		if len(parts) == 1 {
			listReceta(w, r)
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

	id, err = strconv.Atoi(parts[1])
	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un entero positivo", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getReceta(w, r, id)
	case http.MethodPut:
		updateReceta(w, r, id)
	case http.MethodDelete:
		deleteReceta(w, r, id)
	}
}

func postReceta(w http.ResponseWriter, r *http.Request) {
	var input RecetaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Nombre) == "" || strings.TrimSpace(input.Descripcion) == "" || input.IDUsuario <= 0 {
		http.Error(w, "Nombre, descripción y un id_usuario válido son requeridos", http.StatusBadRequest)
		return
	}

	uReturn, err := queries.CreateReceta(r.Context(), sqlc.CreateRecetaParams{
		Nombre:       input.Nombre,
		Descripcion:  input.Descripcion,
		Pasos:        input.Pasos,
		Ingredientes: input.Ingredientes,
		IDUsuario:    input.IDUsuario,
	})
	if err != nil {
		http.Error(w, "Fallo al crear receta: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(uReturn)
}

func listReceta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	recetas, err := queries.ListReceta(r.Context())
	if err != nil {
		http.Error(w, "Fallo al listar recetas", http.StatusInternalServerError)
		return
	}
	// Si está vacía o es nil, escribimos directamente los corchetes vacíos []
	if len(recetas) == 0 {
		w.Write([]byte("[]\n"))
		return
	}
	json.NewEncoder(w).Encode(recetas)
}

func getReceta(w http.ResponseWriter, r *http.Request, id int) {
	u, err := queries.GetReceta(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Receta no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func updateReceta(w http.ResponseWriter, r *http.Request, id int) {
	var input RecetaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Nombre) == "" || strings.TrimSpace(input.Descripcion) == "" {
		http.Error(w, "Nombre y descripción no pueden estar vacíos", http.StatusBadRequest)
		return
	}

	actualizada, err := queries.UpdateReceta(r.Context(), sqlc.UpdateRecetaParams{
		IDReceta:     int32(id),
		Nombre:       input.Nombre,
		Descripcion:  input.Descripcion,
		Pasos:        input.Pasos,
		Ingredientes: input.Ingredientes,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Receta no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actualizada)
}

func deleteReceta(w http.ResponseWriter, r *http.Request, id int) {
	err := queries.DeleteReceta(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

///////////////////// CRUDs de tabla comentario /////////////////////

func comentarioHandlerCRUD(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var id int
	var err error

	switch r.Method {
	case http.MethodPost:
		if len(parts) != 1 {
			http.Error(w, "Ruta inválida", http.StatusBadRequest)
			return
		}
		postComentario(w, r)
		return
	case http.MethodGet:
		if len(parts) == 1 {
			listComentario(w, r)
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

	id, err = strconv.Atoi(parts[1])
	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un entero positivo", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getComentario(w, r, id)
	case http.MethodPut:
		updateComentario(w, r, id)
	case http.MethodDelete:
		deleteComentario(w, r, id)
	}
}

func postComentario(w http.ResponseWriter, r *http.Request) {
	var input ComentarioInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Descripcion) == "" || input.IDUsuario <= 0 || input.IDReceta <= 0 {
		http.Error(w, "Descripcion, id_usuario y id_receta válidos son obligatorios", http.StatusBadRequest)
		return
	}
	if input.Puntuacion < 1 || input.Puntuacion > 5 {
		http.Error(w, "La puntuación debe estar entre 1 y 5", http.StatusBadRequest)
		return
	}

	cReturn, err := queries.CreateComentario(r.Context(), sqlc.CreateComentarioParams{
		IDUsuario:   input.IDUsuario,
		IDReceta:    input.IDReceta,
		Descripcion: input.Descripcion,
		Puntuacion:  input.Puntuacion,
	})
	if err != nil {
		log.Printf("Fallo al crear comentario: %v", err)
		http.Error(w, "Fallo al crear comentario: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cReturn)
}

func listComentario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	recetaParam := r.URL.Query().Get("receta")
	if recetaParam != "" {
		recetaID, err := strconv.Atoi(recetaParam)
		if err == nil && recetaID > 0 {
			comentariosByReceta, err := queries.ListComentarioByReceta(r.Context(), int32(recetaID))
			if err != nil {
				http.Error(w, "Fallo al listar comentarios de la receta", http.StatusInternalServerError)
				return
			}
			if len(comentariosByReceta) == 0 {
				w.Write([]byte("[]\n"))
				return
			}
			json.NewEncoder(w).Encode(comentariosByReceta)
			return
		}
	}

	comentarios, err := queries.ListComentario(r.Context())
	if err != nil {
		log.Printf("Fallo al listar comentarios: %v", err)
		http.Error(w, "Fallo al listar comentarios", http.StatusInternalServerError)
		return
	}
	if len(comentarios) == 0 {
		w.Write([]byte("[]\n"))
		return
	}
	json.NewEncoder(w).Encode(comentarios)
}

func getComentario(w http.ResponseWriter, r *http.Request, id int) {
	c, err := queries.GetComentario(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Comentario no encontrado", http.StatusNotFound) // 404
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func updateComentario(w http.ResponseWriter, r *http.Request, id int) {
	var decoComentario sqlc.UpdateComentarioParams
	if err := json.NewDecoder(r.Body).Decode(&decoComentario); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validación
	if strings.TrimSpace(decoComentario.Descripcion) == "" {
		http.Error(w, "La descripción no puede estar vacía", http.StatusBadRequest)
		return
	}
	if decoComentario.Puntuacion < 1 || decoComentario.Puntuacion > 5 {
		http.Error(w, "La puntuación debe estar entre 1 y 5", http.StatusBadRequest)
		return
	}

	u := sqlc.UpdateComentarioParams{
		IDComentario: int32(id),
		Descripcion:  decoComentario.Descripcion,
		Puntuacion:   decoComentario.Puntuacion,
	}

	actualizado, err := queries.UpdateComentario(r.Context(), u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Comentario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actualizado)
}

func deleteComentario(w http.ResponseWriter, r *http.Request, id int) {
	err := queries.DeleteComentario(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
