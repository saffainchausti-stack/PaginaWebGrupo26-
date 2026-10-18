
# Crear un usuario
POST http://localhost:8080/usuario
Content-Type: application/json

{
  "Nombre": "Usuario",
  "Apellido":  "Hurl",
  "Email": "hurl@gmail.com"
}
HTTP 201
[Captures]
id: jsonpath "$.IDUsuario"
[Asserts]
jsonpath "$.Nombre" == "Usuario"


# Listar usuarios
GET http://localhost:8080/usuario
HTTP 200
[Asserts]
jsonpath "$" isArray


# Consultar el usuario creado
GET http://localhost:8080/usuario/{{id}}
HTTP 200
[Asserts]
jsonpath "$.IDUsuario" == {{id}}


# Modificar el usuario
PUT http://localhost:8080/usuario/{{id}}
Content-Type: application/json

{
  "Nombre": "Usuario",
  "Apellido": "Hurl Modificado",
  "Email": "hurl@example.com"
}
HTTP 200
[Asserts]
jsonpath "$.Apellido" == "Hurl Modificado"


# Eliminar el usuario
DELETE http://localhost:8080/usuario/{{id}}
HTTP 204


# Verificar que ya no exista
GET http://localhost:8080/usuario/{{id}}
HTTP 404