package handlers

import (
	db "ServidorTrabajoWeb/db/sqlc"
	"ServidorTrabajoWeb/views"
	"strconv"

	"log"
	"net/http"

	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type RecetasHandler struct {
	queries *db.Queries
}

func NewRecetasHandler(q *db.Queries) *RecetasHandler {
	return &RecetasHandler{queries: q}
}

func (h *RecetasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myRecetas" && !strings.HasPrefix(r.URL.Path, "/myRecetas/") {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getRecetas(w, r)
	case http.MethodPost:
		h.postReceta(w, r)
	case http.MethodDelete:
		h.deleteReceta(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *RecetasHandler) getRecetas(w http.ResponseWriter, r *http.Request) {

	//Lectura de cookies
	c, errC := r.Cookie("session_token")
	if errC != nil {
		log.Printf("No se encontró la cookie session_token: %v", errC)
		http.Error(w, "No hay una sesión iniciada", http.StatusUnauthorized)
		return
	}

	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		log.Println("La sesión no existe en el mapa sessions")
		http.Error(w, "Sesión inválida", http.StatusUnauthorized)
		return
	}

	if userSession.isExpired() {
		log.Println("La sesión expiró")
		http.Error(w, "La sesión expiró", http.StatusUnauthorized)
		return
	}
	//////////////////////////

	recetas, err := h.queries.ListRecetaByUsuario(r.Context(), userSession.id_usuario)
	if err != nil {
		log.Printf("Error al listar recetas: %v\n", err)
		http.Error(w, "Error al cargar las recetas", http.StatusInternalServerError)
		return
	}

	componente := views.MyRecetas(recetas)
	if err := componente.Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar vista: %v\n", err)
	}
}

func (h *RecetasHandler) postReceta(w http.ResponseWriter, r *http.Request) {
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

	//Lectura de cookies
	c, errC := r.Cookie("session_token")
	if errC != nil {
		log.Printf("No se encontró la cookie session_token: %v", errC)
		http.Error(w, "No hay una sesión iniciada", http.StatusUnauthorized)
		return
	}

	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		log.Println("La sesión no existe en el mapa sessions")
		http.Error(w, "Sesión inválida", http.StatusUnauthorized)
		return
	}

	if userSession.isExpired() {
		log.Println("La sesión expiró")
		http.Error(w, "La sesión expiró", http.StatusUnauthorized)
		return
	}
	//////////////////////////

	_, err := h.queries.CreateReceta(r.Context(), db.CreateRecetaParams{
		Nombre:       nombre,
		Descripcion:  descripcion,
		Pasos:        pasos,
		Ingredientes: ingredientes,
		IDUsuario:    userSession.id_usuario,
	})
	if err != nil {
		log.Printf("Fallo al crear receta: %v\n", err)
		http.Error(w, "Fallo al crear receta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	recetas, err := h.queries.ListRecetaByUsuario(r.Context(), userSession.id_usuario)
	if err != nil {
		http.Error(w, "Error al listar recetas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.RecetaByIdUsuarioList(recetas).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar lista de recetas: %v\n", err)
	}
}

func (h *RecetasHandler) deleteReceta(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/myRecetas/"))
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteReceta(r.Context(), int32(id)); err != nil {
		log.Printf("Error al eliminar receta: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Lectura de cookies
	c, errC := r.Cookie("session_token")
	if errC != nil {
		log.Printf("No se encontró la cookie session_token: %v", errC)
		http.Error(w, "No hay una sesión iniciada", http.StatusUnauthorized)
		return
	}

	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		log.Println("La sesión no existe en el mapa sessions")
		http.Error(w, "Sesión inválida", http.StatusUnauthorized)
		return
	}

	if userSession.isExpired() {
		log.Println("La sesión expiró")
		http.Error(w, "La sesión expiró", http.StatusUnauthorized)
		return
	}
	//////////////////////////

	recetas, err := h.queries.ListRecetaByUsuario(r.Context(), userSession.id_usuario)
	if err != nil {
		log.Printf("Error al listar recetas: %v\n", err)
		http.Error(w, "Error al listar recetas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.RecetaByIdUsuarioList(recetas).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar lista de recetas: %v\n", err)
	}
}

/////////////////////////////////////////////////////////////////////////////////
