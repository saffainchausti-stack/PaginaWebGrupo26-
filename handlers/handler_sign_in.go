package handlers

import (
	db "ServidorTrabajoWeb/db/sqlc"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type SignInHandler struct {
	queries *db.Queries
}

func NewSignInHandler(q *db.Queries) *SignInHandler {
	return &SignInHandler{queries: q}
}

func (h *SignInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "./static/login.html")
}
