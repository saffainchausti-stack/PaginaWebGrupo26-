package main

import (
	"database/sql"
	"fmt"
	"net/http"

	sqlc "ServidorTrabajoWeb/db/sqlc"
)

func main() {
	// 1. Define el directorio que contiene los archivos estáticos.
	staticDir := "./static"

	// 2. Crea un manejador (handler) de servidor de archivos.
	// http.Dir convierte la ruta del directorio en un sistema de archivos HTTP.
	// http.FileServer crea un manejador que sirve archivos desde ese sistema.
	// ¡Automáticamente sirve index.html para directorios!
	fileServer := http.FileServer(http.Dir(staticDir))

	// Para subir pagina web gratis: Netlify drop
	// 3. Registra el manejador para que atienda todas las peticiones ("/").
	// Usamos http.Handle porque fileServer es un http.Handler.
	http.Handle("/", fileServer)

	// 4. Define el puerto y muestra un mensaje
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	// 5. Inicia el servidor HTTP
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}

func setupDB() (*sqlc.Queries, *sql.DB) {
	dbStr := "host=localhost port=5432 user=postgres password=postgres dbname=recetas sslmode=disable"
	db, err := sql.Open("pgx", dbStr)

	if err != nil {
		fmt.Errorf("Error al conectar con la base de datos: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		fmt.Errorf("No se pudo responder al ping de la base de datos: %v", err)
	}

	return sqlc.New(db), db
}
