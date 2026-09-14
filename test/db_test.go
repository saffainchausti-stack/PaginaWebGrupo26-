package test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	sqlc "ServidorTrabajoWeb/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) (*sqlc.Queries, *sql.DB) {
	t.Helper()
	dbStr := "host=localhost port=5432 user=postgres password=postgres dbname=recetas sslmode=disable"
	db, err := sql.Open("pgx", dbStr)

	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("No se pudo responder al ping de la base de datos: %v", err)
	}

	return sqlc.New(db), db
}

func TestUsuarioCRUD(t *testing.T) {
	queries, db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. CreateUsuario
	usuario, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
		Nombre:      "Carlos",
		Apellido:    sql.NullString{String: "Perez", Valid: true},
		Email:       "carlos.perez@gmail.com",
		Contrasenia: "contra123",
	})
	if err != nil {
		t.Fatalf("Fallo CreateUsuario: %v", err)
	}

	if usuario.IDUsuario == 0 {
		t.Errorf("Se esperaba IDUsuario mayor a 0, se obtuvo %d", usuario.IDUsuario)
	}
	if usuario.Nombre != "Carlos" {
		t.Errorf("nombre = %s; esperado Carlos", usuario.Nombre)
	}
	if usuario.Email != "carlos.perez@gmail.com" {
		t.Errorf("email = %s; esperado carlos.perez@gmail.com", usuario.Email)
	}

	// 2. GetUsuario
	obtenido, err := queries.GetUsuario(ctx, usuario.IDUsuario)
	if err != nil {
		t.Fatalf("Fallo GetUsuario: %v", err)
	}
	if obtenido.Email != "carlos.perez@gmail.com" {
		t.Errorf("GetUsuario: esperado email 'carlos.perez@gmail.com', se obtuvo '%s'", obtenido.Email)
	}

	// 3. GetContrasenia
	passRow, err := queries.GetContrasenia(ctx, usuario.IDUsuario)
	if err != nil {
		t.Fatalf("Fallo GetContrasenia: %v", err)
	}
	if passRow.Contrasenia != "contra123" {
		t.Errorf("GetContrasenia: esperada 'contra123', se obtuvo '%s'", passRow.Contrasenia)
	}

	// 4. ListUsuario
	usuarios, err := queries.ListUsuario(ctx)
	if err != nil {
		t.Fatalf("Fallo ListUsuario: %v", err)
	}
	if len(usuarios) == 0 {
		t.Errorf("ListUsuario: se esperaba al menos 1 usuario en la lista")
	}

	// 5. Update
	subtests := []struct {
		nombre    string
		ejecutar  func() error
		verificar func(u sqlc.GetUsuarioRow) bool
	}{
		{
			nombre: "UpdateUsuarioNombre",
			ejecutar: func() error {
				return queries.UpdateUsuarioNombre(ctx, sqlc.UpdateUsuarioNombreParams{
					IDUsuario: usuario.IDUsuario,
					Nombre:    "Ignacio Valentin Martin",
				})
			},
			verificar: func(u sqlc.GetUsuarioRow) bool { return u.Nombre == "Ignacio Valentin Martin" },
		},
		{
			nombre: "UpdateUsuarioApellido",
			ejecutar: func() error {
				return queries.UpdateUsuarioApellido(ctx, sqlc.UpdateUsuarioApellidoParams{
					IDUsuario: usuario.IDUsuario,
					Apellido:  sql.NullString{String: "Gomez", Valid: true},
				})
			},
			verificar: func(u sqlc.GetUsuarioRow) bool { return u.Apellido.String == "Gomez" },
		},
		{
			nombre: "UpdateUsuarioEmail",
			ejecutar: func() error {
				return queries.UpdateUsuarioEmail(ctx, sqlc.UpdateUsuarioEmailParams{
					IDUsuario: usuario.IDUsuario,
					Email:     "ignaciovalentinmartin@gmail.com",
				})
			},
			verificar: func(u sqlc.GetUsuarioRow) bool { return u.Email == "ignaciovalentinmartin@gmail.com" },
		},
		{
			nombre: "UpdateUsuarioContrasenia",
			ejecutar: func() error {
				return queries.UpdateUsuarioContrasenia(ctx, sqlc.UpdateUsuarioContraseniaParams{
					IDUsuario:   usuario.IDUsuario,
					Contrasenia: "nuevaClave456",
				})
			},
			verificar: func(u sqlc.GetUsuarioRow) bool {
				p, err := queries.GetContrasenia(ctx, usuario.IDUsuario)
				if err != nil {
					return false
				}
				return p.Contrasenia == "nuevaClave456"
			},
		},
	}

	for _, tc := range subtests {
		t.Run(tc.nombre, func(t *testing.T) {
			if err := tc.ejecutar(); err != nil {
				t.Fatalf("Fallo al ejecutar %s: %v", tc.nombre, err)
			}
			actualizado, err := queries.GetUsuario(ctx, usuario.IDUsuario)
			if err != nil {
				t.Fatalf("Error al recuperar usuario actualizado en %s: %v", tc.nombre, err)
			}
			if !tc.verificar(actualizado) {
				t.Errorf("La verificacion fallo en el subtest %s", tc.nombre)
			}
		})
	}

	// 6. DeleteUsuario
	err = queries.DeleteUsuario(ctx, usuario.IDUsuario)
	if err != nil {
		t.Fatalf("Fallo DeleteUsuario: %v", err)
	}

	_, err = queries.GetUsuario(ctx, usuario.IDUsuario)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows despues de DeleteUsuario, se obtuvo: %v", err)
	}
}

func TestRecetaCRUD(t *testing.T) {
	queries, db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Necesitamos un usuario para la Foreign Key
	autor, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
		Nombre:      "Chef",
		Apellido:    sql.NullString{String: "Gusteau", Valid: true},
		Email:       "gusteau@recetas.com",
		Contrasenia: "cocinorico",
	})
	if err != nil {
		t.Fatalf("Fallo al crear autor: %v", err)
	}
	defer queries.DeleteUsuario(ctx, autor.IDUsuario)

	// 1. CreateReceta
	receta, err := queries.CreateReceta(ctx, sqlc.CreateRecetaParams{
		Nombre:       "Ratatouille",
		Descripcion:  "Guiso clasico de verduras provenzales",
		Pasos:        "Cortar vegetales en rodajas finas, hornear a fuego lento.",
		Ingredientes: "Berenjenas, calabacin, pimientos, tomates, aceite de oliva",
		IDUsuario:    autor.IDUsuario,
	})
	if err != nil {
		t.Fatalf("Fallo CreateReceta: %v", err)
	}

	if receta.IDReceta == 0 {
		t.Errorf("Se esperaba IDReceta mayor a 0, se obtuvo %d", receta.IDReceta)
	}

	// 2. GetReceta
	recetaObtenida, err := queries.GetReceta(ctx, receta.IDReceta)
	if err != nil {
		t.Fatalf("Fallo GetReceta: %v", err)
	}
	if recetaObtenida.Nombre != "Ratatouille" {
		t.Errorf("GetReceta: esperado 'Ratatouille', obtenido '%s'", recetaObtenida.Nombre)
	}

	// 3. ListReceta
	recetas, err := queries.ListReceta(ctx)
	if err != nil {
		t.Fatalf("Fallo ListReceta: %v", err)
	}
	if len(recetas) == 0 {
		t.Errorf("ListReceta: no devolvio ninguna receta")
	}

	// 4. ListRecetaByUsuario
	recetasUsuario, err := queries.ListRecetaByUsuario(ctx, autor.IDUsuario)
	if err != nil {
		t.Fatalf("Fallo ListRecetaByUsuario: %v", err)
	}
	if len(recetasUsuario) != 1 {
		t.Errorf("ListRecetaByUsuario: esperada 1 receta para el usuario, se obtuvieron %d", len(recetasUsuario))
	}

	// 5. Subtests para los diferentes Updates de Receta
	casosUpdate := []struct {
		nombre   string
		ejecutar func() error
		validar  func(r sqlc.GetRecetaRow) bool
	}{
		{
			nombre: "UpdateRecetaNombre",
			ejecutar: func() error {
				return queries.UpdateRecetaNombre(ctx, sqlc.UpdateRecetaNombreParams{
					IDReceta: receta.IDReceta,
					Nombre:   "Ratatouille Tradicional",
				})
			},
			validar: func(r sqlc.GetRecetaRow) bool { return r.Nombre == "Ratatouille Tradicional" },
		},
		{
			nombre: "UpdateRecetaDescripcion",
			ejecutar: func() error {
				return queries.UpdateRecetaDescripcion(ctx, sqlc.UpdateRecetaDescripcionParams{
					IDReceta:    receta.IDReceta,
					Descripcion: "Nueva descripcion",
				})
			},
			validar: func(r sqlc.GetRecetaRow) bool { return r.Descripcion == "Nueva descripcion" },
		},
		{
			nombre: "UpdateRecetaIngredientes",
			ejecutar: func() error {
				return queries.UpdateRecetaIngredientes(ctx, sqlc.UpdateRecetaIngredientesParams{
					IDReceta:     receta.IDReceta,
					Ingredientes: "Berenjenas, calabacin, pimientos, tomates, aceite de oliva, chorizo",
				})
			},
			validar: func(r sqlc.GetRecetaRow) bool {
				return r.Ingredientes == "Berenjenas, calabacin, pimientos, tomates, aceite de oliva, chorizo"
			},
		},
		{
			nombre: "UpdateRecetaPasos",
			ejecutar: func() error {
				return queries.UpdateRecetaPasos(ctx, sqlc.UpdateRecetaPasosParams{
					IDReceta: receta.IDReceta,
					Pasos:    "Cortar vegetales en rodajas finas, hornear a fuego lento, hornear durante media hora",
				})
			},
			validar: func(r sqlc.GetRecetaRow) bool {
				return r.Pasos == "Cortar vegetales en rodajas finas, hornear a fuego lento, hornear durante media hora"
			},
		},
	}

	for _, tc := range casosUpdate {
		t.Run(tc.nombre, func(t *testing.T) {
			if err := tc.ejecutar(); err != nil {
				t.Fatalf("Fallo en %s: %v", tc.nombre, err)
			}
			actual, err := queries.GetReceta(ctx, receta.IDReceta)
			if err != nil {
				t.Fatalf("Error al obtener receta en %s: %v", tc.nombre, err)
			}
			if !tc.validar(actual) {
				t.Errorf("Validacion fallida para %s", tc.nombre)
			}
		})
	}

	// 6. DeleteReceta
	err = queries.DeleteReceta(ctx, receta.IDReceta)
	if err != nil {
		t.Fatalf("Fallo DeleteReceta: %v", err)
	}

	_, err = queries.GetReceta(ctx, receta.IDReceta)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows despues de borrar receta, se obtuvo: %v", err)
	}
}

func TestComentarioCRUD(t *testing.T) {
	queries, db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Crear autor y receta para las Foreign Keys
	usuario, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
		Nombre:      "Critico",
		Apellido:    sql.NullString{String: "Ego", Valid: true},
		Email:       "antonEgo@critica.com",
		Contrasenia: "claveCritico1",
	})
	if err != nil {
		t.Fatalf("Error al crear usuario para comentario: %v", err)
	}
	defer queries.DeleteUsuario(ctx, usuario.IDUsuario)

	receta, err := queries.CreateReceta(ctx, sqlc.CreateRecetaParams{
		Nombre:       "Sopa de Cebolla",
		Descripcion:  "Sopa clasica francesa con queso gratinado",
		Pasos:        "Caramelizar cebollas, agregar caldo y servir con pan tostado.",
		Ingredientes: "Cebollas, manteca, caldo de carne, queso gruyere, pan",
		IDUsuario:    usuario.IDUsuario,
	})
	if err != nil {
		t.Fatalf("Error al crear receta para comentario: %v", err)
	}
	defer queries.DeleteReceta(ctx, receta.IDReceta)

	// 1. CreateComentario
	comentario, err := queries.CreateComentario(ctx, sqlc.CreateComentarioParams{
		IDUsuario:   usuario.IDUsuario,
		IDReceta:    receta.IDReceta,
		Descripcion: "Sabor excelente, textura perfecta.",
		Puntuacion:  5,
	})
	if err != nil {
		t.Fatalf("Fallo CreateComentario: %v", err)
	}
	if comentario.IDComentario == 0 {
		t.Errorf("Se esperaba IDComentario mayor a 0, se obtuvo %d", comentario.IDComentario)
	}

	// 2. GetComentario
	comentarioObtenido, err := queries.GetComentario(ctx, comentario.IDComentario)
	if err != nil {
		t.Fatalf("Fallo GetComentario: %v", err)
	}
	if comentarioObtenido.Descripcion != "Sabor excelente, textura perfecta." {
		t.Errorf("GetComentario: descripcion inesperada '%s'", comentarioObtenido.Descripcion)
	}

	// 3. ListComentarioByReceta
	listaComentarios, err := queries.ListComentarioByReceta(ctx, receta.IDReceta)
	if err != nil {
		t.Fatalf("Fallo ListComentarioByReceta: %v", err)
	}
	if len(listaComentarios) != 1 {
		t.Errorf("ListComentarioByReceta: se esperaba 1 comentario, se encontraron %d", len(listaComentarios))
	}

	//4. CreateComentario Puntuacion invalida
	_, err = queries.CreateComentario(ctx, sqlc.CreateComentarioParams{
		IDUsuario:   usuario.IDUsuario,
		IDReceta:    receta.IDReceta,
		Descripcion: "Puntuación inválida",
		Puntuacion:  6,
	})

	if err == nil {
		t.Errorf("Se esperaba un error al crear un comentario con puntuación 6")
	}

	// 5. Updates de Comentario
	subtests := []struct {
		nombre   string
		ejecutar func() error
		validar  func(c sqlc.GetComentarioRow) bool
	}{
		{
			nombre: "UpdateComentarioDescripcion",
			ejecutar: func() error {
				return queries.UpdateComentarioDescripcion(ctx, sqlc.UpdateComentarioDescripcionParams{
					IDComentario: comentario.IDComentario,
					Descripcion:  "Esta horrible, lo peor que comi en mi vida",
				})
			},
			validar: func(c sqlc.GetComentarioRow) bool {
				return c.Descripcion == "Esta horrible, lo peor que comi en mi vida"
			},
		},
		{
			nombre: "UpdateComentarioPuntuacion",
			ejecutar: func() error {
				return queries.UpdateComentarioPuntuacion(ctx, sqlc.UpdateComentarioPuntuacionParams{
					IDComentario: comentario.IDComentario,
					Puntuacion:   1,
				})
			},
			validar: func(c sqlc.GetComentarioRow) bool {
				return c.Puntuacion == 1
			},
		},
	}

	for _, tc := range subtests {
		t.Run(tc.nombre, func(t *testing.T) {
			if err := tc.ejecutar(); err != nil {
				t.Fatalf("Fallo en %s: %v", tc.nombre, err)
			}
			actual, err := queries.GetComentario(ctx, comentario.IDComentario)
			if err != nil {
				t.Fatalf("Error al obtener comentario tras update en %s: %v", tc.nombre, err)
			}
			if !tc.validar(actual) {
				t.Errorf("Validacion fallo en %s", tc.nombre)
			}
		})
	}

	// 6. DeleteComentario
	err = queries.DeleteComentario(ctx, comentario.IDComentario)
	if err != nil {
		t.Fatalf("Fallo DeleteComentario: %v", err)
	}

	_, err = queries.GetComentario(ctx, comentario.IDComentario)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows luego de borrar el comentario, se obtuvo: %v", err)
	}
}
