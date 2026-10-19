#---------------------------------------------------------------------------------
# Crear un usuario con parametros incorrectos
POST http://localhost:8080/usuario
Content-Type: application/json

{
  "nombre": "Usuario",
  "apellido":  "",
  "email": "hurl@gmail.com",
  "contrasenia": "1234"
}
HTTP 400


# Crear un usuario
POST http://localhost:8080/usuario
Content-Type: application/json

{
  "nombre": "Usuario",
  "apellido":  "Hurl",
  "email": "hurl@gmail.com",
  "contrasenia": "1234"
}
HTTP 201
[Captures]
idU: jsonpath "$.id_usuario"
[Asserts]
jsonpath "$.nombre" == "Usuario"


# Listar usuarios
GET http://localhost:8080/usuario
HTTP 200
[Asserts]
jsonpath "$" isCollection


# Consultar el usuario creado
GET http://localhost:8080/usuario/{{idU}}
HTTP 200
[Asserts]
jsonpath "$.id_usuario" == {{idU}}


# Modificar el usuario
PUT http://localhost:8080/usuario/{{idU}}
Content-Type: application/json

{
  "nombre": "Usuario",
  "apellido": "Hurl Modificado",
  "email": "hurl@gmail.com",
  "contrasenia": "1234"
}
HTTP 200
[Asserts]
jsonpath "$.apellido" == "Hurl Modificado"


#---------------------------------------------------------------------------------
# Crear una receta con parametros incorrectos
POST http://localhost:8080/receta
Content-Type: application/json

{
  "nombre": "ensalada",
  "descripcion":  "",
  "ingredientes": "lechuga",
  "pasos": "cortar lechuga"
}
HTTP 400


# Crear una receta
POST http://localhost:8080/receta
Content-Type: application/json

{
  "nombre": "ensalada",
  "descripcion":  "buena ensalada",
  "ingredientes": "lechuga",
  "pasos": "cortar lechuga",
  "id_usuario": {{idU}}
}
HTTP 201
[Captures]
idR: jsonpath "$.id_receta"
[Asserts]
jsonpath "$.nombre" == "ensalada"


# Listar receta
GET http://localhost:8080/receta
HTTP 200
[Asserts]
jsonpath "$" isCollection


# Consultar la receta creada
GET http://localhost:8080/receta/{{idR}}
HTTP 200
[Asserts]
jsonpath "$.id_receta" == {{idR}}


# Modificar la receta
PUT http://localhost:8080/receta/{{idR}}
Content-Type: application/json

{
  "nombre": "ensalada",
  "descripcion":  "buena ensalada",
  "ingredientes": "lechuga",
  "pasos": "cortar lechuga"
}
HTTP 200
[Asserts]
jsonpath "$.descripcion" == "buena ensalada"

# Crear un comentario con parametros incorrectos (puntuación fuera de rango 1-5)
POST http://localhost:8080/comentario
Content-Type: application/json

{
  "id_usuario": {{idU}},
  "id_receta": {{idR}},
  "descripcion": "Comentario inválido",
  "puntuacion": 10
}
HTTP 400

#---------------------------------------------------------------------------------
# Crear un comentario válido
POST http://localhost:8080/comentario
Content-Type: application/json

{
  "id_usuario": {{idU}},
  "id_receta": {{idR}},
  "descripcion": "Excelente receta, muy rica",
  "puntuacion": 5
}
HTTP 201
[Captures]
idC: jsonpath "$.id_comentario"
[Asserts]
jsonpath "$.descripcion" == "Excelente receta, muy rica"
jsonpath "$.puntuacion" == 5
jsonpath "$.id_receta" == {{idR}}
jsonpath "$.id_usuario" == {{idU}}


# Listar todos los comentarios
GET http://localhost:8080/comentario
HTTP 200
[Asserts]
jsonpath "$" isCollection


# Listar comentarios filtrados por receta (query param)
GET http://localhost:8080/comentario?receta={{idR}}
HTTP 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[0].id_receta" == {{idR}}


# Consultar el comentario creado por ID
GET http://localhost:8080/comentario/{{idC}}
HTTP 200
[Asserts]
jsonpath "$.id_comentario" == {{idC}}
jsonpath "$.descripcion" == "Excelente receta, muy rica"


# Modificar el comentario con datos inválidos (puntuación menor a 1)
PUT http://localhost:8080/comentario/{{idC}}
Content-Type: application/json

{
  "descripcion": "Intento de modificación fallida",
  "puntuacion": 0
}
HTTP 400


# Modificar el comentario exitosamente
PUT http://localhost:8080/comentario/{{idC}}
Content-Type: application/json

{
  "descripcion": "Excelente receta, ahora me gusta más",
  "puntuacion": 4
}
HTTP 200
[Asserts]
jsonpath "$.descripcion" == "Excelente receta, ahora me gusta más"
jsonpath "$.puntuacion" == 4


#---------------------------------------------------------------------------------
# Eliminar el comentario
DELETE http://localhost:8080/comentario/{{idC}}
HTTP 204


# Verificar que el comentario ya no exista
GET http://localhost:8080/comentario/{{idC}}
HTTP 404


# Eliminar el receta
DELETE http://localhost:8080/receta/{{idR}}
HTTP 204


# Verificar que ya no exista
GET http://localhost:8080/receta/{{idR}}
HTTP 404


# Eliminar el usuario
DELETE http://localhost:8080/usuario/{{idU}}
HTTP 204


# Verificar que ya no exista
GET http://localhost:8080/usuario/{{idU}}
HTTP 404
