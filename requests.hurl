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