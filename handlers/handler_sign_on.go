package handlers

import (
	db "ServidorTrabajoWeb/db/sqlc"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type SignOnHandler struct {
	queries *db.Queries
}

func NewSignOnHandler(q *db.Queries) *SignOnHandler {
	return &SignOnHandler{queries: q}
}

func (h *SignOnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./static/signOn.html")
		return
	} else if r.Method == http.MethodPost {
		nombre := r.FormValue("name")
		apellido := r.FormValue("surname")
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, err := h.queries.CreateUsuario(r.Context(),
			db.CreateUsuarioParams{Nombre: nombre,
				Apellido:    apellido,
				Email:       email,
				Contrasenia: password,
			})

		if err != nil {
			http.Error(w, "Hubo un problema con el registro del usuario", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
}
