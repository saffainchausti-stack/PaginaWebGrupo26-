-- name: GetUsuario :one
SELECT id_usuario, nombre, apellido, email
FROM usuario
WHERE id_usuario = $1;

-- name: GetContrasenia :one
SELECT id_usuario, contrasenia
FROM usuario
WHERE id_usuario = $1;

-- name: ListUsuario :many
SELECT id_usuario, nombre, apellido, email
FROM usuario
ORDER BY id_usuario;


-- name: CreateUsuario :one
INSERT INTO usuario (nombre, apellido, email, contrasenia)
VALUES ($1, $2, $3, $4)
RETURNING id_usuario, nombre, apellido, email;


-- name: UpdateUsuarioNombre :exec
UPDATE usuario
SET nombre = $2
WHERE id_usuario = $1;

-- name: UpdateUsuarioEmail :exec
UPDATE usuario
SET email = $2
WHERE id_usuario = $1;

-- name: UpdateUsuarioApellido :exec
UPDATE usuario
SET apellido = $2
WHERE id_usuario = $1;

-- name: UpdateUsuarioContrasenia :exec
UPDATE usuario
SET contrasenia = $2
WHERE id_usuario = $1;

-- name: DeleteUsuario :exec
DELETE FROM usuario
WHERE id_usuario = $1;



-- name: GetReceta :one
SELECT id_receta, nombre, descripcion, pasos, ingredientes, id_usuario
FROM receta
WHERE id_receta = $1;


-- name: ListRecetaByUsuario :many
SELECT id_receta, nombre, descripcion, pasos, ingredientes, id_usuario
FROM receta
WHERE id_usuario=$1
ORDER BY id_receta;

-- name: ListReceta :many
SELECT id_receta, nombre, descripcion, pasos, ingredientes, id_usuario
FROM receta
ORDER BY id_receta;


-- name: CreateReceta :one
INSERT INTO receta (nombre, descripcion, pasos, ingredientes, id_usuario)
VALUES ($1, $2, $3, $4, $5)
RETURNING id_receta, nombre, descripcion, pasos, ingredientes, id_usuario;


-- name: UpdateRecetaNombre :exec
UPDATE receta
SET nombre = $2
WHERE id_receta = $1;

-- name: UpdateRecetaDescripcion :exec
UPDATE receta
SET descripcion = $2
WHERE id_receta = $1;

-- name: UpdateRecetaPasos :exec
UPDATE receta
SET pasos = $2
WHERE id_receta = $1;

-- name: UpdateRecetaIngredientes :exec
UPDATE receta
SET ingredientes = $2
WHERE id_receta = $1;

-- name: DeleteReceta :exec
DELETE FROM receta
WHERE id_receta = $1;



-- name: GetComentario :one
SELECT id_comentario, id_usuario, id_receta, descripcion, puntuacion
FROM comentario
WHERE id_comentario = $1;

-- name: ListComentarioByReceta :many
SELECT id_comentario, id_usuario, descripcion, puntuacion, id_receta
FROM comentario
WHERE id_receta=$1
ORDER BY id_comentario;


-- name: CreateComentario :one
INSERT INTO comentario (id_usuario, id_receta, descripcion, puntuacion)
VALUES ($1, $2, $3, $4)
RETURNING id_comentario, id_usuario, id_receta, descripcion, puntuacion;

-- name: UpdateComentarioDescripcion :exec
UPDATE comentario
SET descripcion = $2
WHERE id_comentario = $1;

-- name: UpdateComentarioPuntuacion :exec
UPDATE comentario
SET puntuacion = $2
WHERE id_comentario = $1;


-- name: DeleteComentario :exec
DELETE FROM comentario
WHERE id_comentario = $1;