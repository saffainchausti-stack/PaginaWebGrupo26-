package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	db "ServidorTrabajoWeb/db/sqlc"
)

func setupTestDB(t *testing.T) (*db.Queries, *sql.DB) {
	connStr := "postgres://postgres:postgres@localhost:5432/recetas?sslmode=disable"
	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	if err := conn.Ping(); err != nil {
		t.Fatalf("No se pudo hacer ping a la base de datos: %v", err)
	}

	return db.New(conn), conn
}

func TestDominioRecetasCRUD(t *testing.T) {
	queries, conn := setupTestDB(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. TEST USUARIO: Crear un usuario
	usuarioCreado, err := queries.CreateUsuario(ctx, db.CreateUsuarioParams{
		Nombre:      "Chefcito",
		Apellido:    sql.NullString{String: "Ratatouille", Valid: true},
		Email:       "chef@recetas.com",
		Contrasenia: "supersecret123",
	})
	if err != nil {
		t.Fatalf("Fallo CreateUsuario: %v", err)
	}

	if usuarioCreado.IDUsuario == 0 {
		t.Errorf("Se esperaba IDUsuario autoincremental > 0, se obtuvo 0")
	}

	// 2. TEST RECETA: Crear una receta vinculada al usuario
	recetaCreada, err := queries.CreateReceta(ctx, db.CreateRecetaParams{
		Nombre:       "Guiso de Lentejas",
		Descripcion:  "Guiso tradicional de invierno",
		Pasos:        "1. Sofreír verduras. 2. Agregar lentejas y caldo. 3. Cocinar 40 min.",
		Ingredientes: "Lentejas, cebolla, zanahoria, chorizo colorado, caldo",
		IDUsuario:    usuarioCreado.IDUsuario,
	})
	if err != nil {
		t.Fatalf("Fallo CreateReceta: %v", err)
	}

	if recetaCreada.IDReceta == 0 {
		t.Errorf("Se esperaba IDReceta > 0")
	}

	// 3. TEST GET RECETA
	recetaObtenida, err := queries.GetReceta(ctx, recetaCreada.IDReceta)
	if err != nil {
		t.Fatalf("Fallo GetReceta: %v", err)
	}
	if recetaObtenida.Nombre != "Guiso de Lentejas" {
		t.Errorf("Se esperaba 'Guiso de Lentejas', se obtuvo: %s", recetaObtenida.Nombre)
	}

	// 4. TEST LIST RECETA
	recetas, err := queries.ListReceta(ctx)
	if err != nil {
		t.Fatalf("Fallo ListReceta: %v", err)
	}
	if len(recetas) == 0 {
		t.Errorf("La lista de recetas no debería estar vacía")
	}

	// 5. TEST COMENTARIO: Crear comentario en la receta
	comentarioCreado, err := queries.CreateComentario(ctx, db.CreateComentarioParams{
		IDUsuario:   usuarioCreado.IDUsuario,
		IDReceta:    recetaCreada.IDReceta,
		Descripcion: "¡Quedó espectacular! Muy bien explicados los pasos.",
		Puntuacion:  5,
	})
	if err != nil {
		t.Fatalf("Fallo CreateComentario: %v", err)
	}

	// 6. TEST LIST COMENTARIOS POR RECETA
	comentarios, err := queries.ListComentarioByReceta(ctx, recetaCreada.IDReceta)
	if err != nil {
		t.Fatalf("Fallo ListComentarioByReceta: %v", err)
	}
	if len(comentarios) != 1 {
		t.Errorf("Se esperaba 1 comentario, se encontraron %d", len(comentarios))
	}

	// 7. TEST UPDATE: Modificar nombre de la receta
	err = queries.UpdateRecetaNombre(ctx, db.UpdateRecetaNombreParams{
		IDReceta: recetaCreada.IDReceta,
		Nombre:   "Guiso de Lentejas Casero",
	})
	if err != nil {
		t.Fatalf("Fallo UpdateRecetaNombre: %v", err)
	}

	recetaModificada, err := queries.GetReceta(ctx, recetaCreada.IDReceta)
	if err != nil {
		t.Fatalf("Fallo al obtener receta luego de update: %v", err)
	}
	if recetaModificada.Nombre != "Guiso de Lentejas Casero" {
		t.Errorf("Se esperaba nombre actualizado, se obtuvo: %s", recetaModificada.Nombre)
	}

	// 8. TEST DELETE (en orden inverso de FK): Comentario -> Receta -> Usuario
	err = queries.DeleteComentario(ctx, comentarioCreado.IDComentario)
	if err != nil {
		t.Fatalf("Fallo DeleteComentario: %v", err)
	}

	err = queries.DeleteReceta(ctx, recetaCreada.IDReceta)
	if err != nil {
		t.Fatalf("Fallo DeleteReceta: %v", err)
	}

	// Verificar que no exista más la receta
	_, err = queries.GetReceta(ctx, recetaCreada.IDReceta)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows para receta eliminada, se obtuvo: %v", err)
	}

	err = queries.DeleteUsuario(ctx, usuarioCreado.IDUsuario)
	if err != nil {
		t.Fatalf("Fallo DeleteUsuario: %v", err)
	}
}
