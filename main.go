package main

import (
	sqlc "ServidorTrabajoWeb/db/sqlc"
	"ServidorTrabajoWeb/views"

	//"context"
	"database/sql"

	//"encoding/json"
	//"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var queries *sqlc.Queries

func abrirDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Valor por defecto para correr la API fuera de Docker
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

	// Define el directorio que contiene los archivos estáticos.
	staticDir := "./static"
	// Crea un manejador (handler) de servidor de archivos.
	// http.Dir convierte la ruta del directorio en un sistema de archivos HTTP.
	// http.FileServer crea un manejador que sirve archivos desde ese sistema.
	// ¡Automáticamente sirve index.html para directorios!
	//fileServer := http.FileServer(http.Dir(staticDir))
	// Para subir pagina web gratis: Netlify drop
	// Registra el manejador para que atienda todas las peticiones ("/").
	// Usamos http.Handle porque fileServer es un http.Handler.
	//http.Handle("/", fileServer)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Evitar que capture rutas no deseadas si no coinciden exactamente con "/"
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

		// Renderizar vista
		componente := views.Index(recetas)
		if err := componente.Render(r.Context(), w); err != nil {
			log.Printf("Error al renderizar vista: %v\n", err)
		}
	})

	//http.HandleFunc("/usuario", usuariosHandler)
	//http.HandleFunc("/usuario/", usuarioHandler)

	//http.HandleFunc("/receta", postReceta)
	http.HandleFunc("/receta/{id}", postReceta)
	/*
		http.HandleFunc("/comentario", comentariosHandler)
		http.HandleFunc("/comentario/", comentarioHandler)
	*/
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

//--DTO receta

type RecetaInput struct {
	Nombre       string `json:"nombre"`
	Descripcion  string `json:"descripcion"`
	Ingredientes string `json:"ingredientes"`
	Pasos        string `json:"pasos"`
	IDUsuario    int32  `json:"id_usuario"`
}

func postReceta(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer id_usuario desde la URL (/receta/1)
	idStr := r.PathValue("id")
	idUsuario, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || idUsuario < 0 {
		http.Error(w, "ID de usuario inválido en la URL", http.StatusBadRequest)
		return
	}

	// 2. Parsear los datos del formulario HTML
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

	// 3. Crear en base de datos
	_, err = queries.CreateReceta(r.Context(), sqlc.CreateRecetaParams{
		Nombre:       nombre,
		Descripcion:  descripcion,
		Pasos:        pasos,
		Ingredientes: ingredientes,
		IDUsuario:    int32(idUsuario),
	})
	if err != nil {
		log.Printf("Fallo al crear receta: %v\n", err)
		http.Error(w, "Fallo al crear receta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirigir al inicio para ver la lista actualizada
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
