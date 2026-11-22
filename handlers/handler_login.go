package handlers

import (
	db "ServidorTrabajoWeb/db/sqlc"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type LoginHandler struct {
	queries *db.Queries
}

func NewLoginHandler(q *db.Queries) *LoginHandler {
	return &LoginHandler{queries: q}
}

// Mapa para almacenar las sesiones activas
// La clave es el "session_token" (un UUID)
// NOTA: sin sincronización ni persistencia. En producción, proteger con sync.Mutex
// o usar un almacén compartido (ej. Redis).
var sessions = map[string]session{}

// La estructura de una sesión
type session struct {
	id_usuario int32
	expiry     time.Time
}

// Método para verificar si la sesión ha expirado
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./static/login.html")

	} else if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		id, err := h.queries.GetIdUsuarioEmailContrasenia(r.Context(),
			db.GetIdUsuarioEmailContraseniaParams{Email: email,
				Contrasenia: password,
			})

		if err != nil {
			log.Printf("Error al obtener usuario: %v", err)
			http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
			return
		}

		// Crear un nuevo token de sesión aleatorio (paquete estándar `uuid`, Go 1.27)
		sessionToken := uuid.New().String()
		expiresAt := time.Now().Add(3600 * time.Second)

		// Almacenar la sesión en el mapa del servidor
		sessions[sessionToken] = session{
			id_usuario: id,
			expiry:     expiresAt,
		}

		// Enviar la cookie al cliente
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			Expires:  expiresAt,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/recetas", http.StatusSeeOther)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
}
