package handlers

import (
	db "ServidorTrabajoWeb/db/sqlc"
	"ServidorTrabajoWeb/views"

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
	if r.URL.Path != "/recetas" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getRecetas(w, r)

	case http.MethodPost:
		h.postReceta(w, r)

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *RecetasHandler) getRecetas(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("No se encontró la cookie session_token: %v", err)
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

	recetas, err := h.queries.ListRecetaByUsuario(r.Context(), userSession.id_usuario)
	if err != nil {
		log.Printf("Error al listar recetas: %v\n", err)
		http.Error(w, "Error al cargar las recetas", http.StatusInternalServerError)
		return
	}

	componente := views.Index(recetas)
	if err := componente.Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar vista: %v\n", err)
	}
}
func (h *RecetasHandler) postReceta(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("No se encontró la cookie session_token: %v", err)
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

	// 2. Parsear los datos del formulario HTML
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(r.FormValue("nombre"))
	descripcion := strings.TrimSpace(r.FormValue("descripcion"))
	ingredientes := strings.TrimSpace(r.FormValue("ingredientes"))
	pasos := strings.TrimSpace(r.FormValue("pasos"))

	if nombre == "" || descripcion == "" || ingredientes == "" || pasos == "" {
		http.Error(w, "No pueden haber campos vacios", http.StatusBadRequest)
		return
	}

	// 3. Crear en base de datos
	_, err = h.queries.CreateReceta(r.Context(), db.CreateRecetaParams{
		Nombre:       nombre,
		Descripcion:  descripcion,
		Pasos:        pasos,
		Ingredientes: ingredientes,
		IDUsuario:    int32(userSession.id_usuario),
	})

	if err != nil {
		log.Printf("Fallo al crear receta: %v\n", err)
		http.Error(w, "Fallo al crear receta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/recetas", http.StatusSeeOther)
}
