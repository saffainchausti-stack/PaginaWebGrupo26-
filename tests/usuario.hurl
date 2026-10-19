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
id: jsonpath "$.id_usuario"
[Asserts]
jsonpath "$.nombre" == "Usuario"


# Listar usuarios
GET http://localhost:8080/usuario
HTTP 200
[Asserts]
jsonpath "$" isCollection


# Consultar el usuario creado
GET http://localhost:8080/usuario/{{id}}
HTTP 200
[Asserts]
jsonpath "$.id_usuario" == {{id}}


# Modificar el usuario
PUT http://localhost:8080/usuario/{{id}}
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


# Eliminar el usuario
DELETE http://localhost:8080/usuario/{{id}}
HTTP 204


# Verificar que ya no exista
GET http://localhost:8080/usuario/{{id}}
HTTP 404
