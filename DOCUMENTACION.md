# Documentación del proyecto

## Persistencia

La aplicación utiliza PostgreSQL como sistema de gestión de base de datos para mantener la información de la aplicación.

La base de datos contiene tres entidades principales:

- **Usuario:** almacena la información de los usuarios de la aplicación.
- **Receta:** almacena las recetas publicadas por los usuarios.
- **Comentario:** almacena los comentarios y puntuaciones realizados sobre las recetas.

Utilizamos docker compose para levantar la base de datos y que no dependa de la maquina en la que se ejecuta.

## Tablas

La base de datos está compuesta por las siguientes tablas:

### Usuario

| Campo | Tipo | Descripción |
|---|---|---|
| id_usuario | INTEGER | Identificador del usuario |
| nombre | VARCHAR | Nombre del usuario |
| apellido | VARCHAR | Apellido del usuario |
| email | VARCHAR | Correo electrónico |
| contrasenia | VARCHAR | Contraseña |

id_usuario es la clave primaria de la tabla.

### Receta

| Campo | Tipo | Descripción |
|---|---|---|
| id_receta | INTEGER | Identificador de la receta |
| nombre | VARCHAR | Nombre de la receta |
| descripcion | VARCHAR | Descripción de la receta |
| ingredientes | TEXT | Ingredientes utilizados |
| pasos | TEXT | Pasos de preparación |
| id_usuario | INTEGER | Usuario que creó la receta |

id_receta es la clave primaria.

id_usuario es una clave foránea que relaciona cada receta con el usuario que la creó.

### Comentario

| Campo | Tipo | Descripción |
|---|---|---|
| id_comentario | INTEGER | Identificador único del comentario |
| id_usuario | INTEGER | Usuario que realizó el comentario |
| id_receta | INTEGER | Receta comentada |
| descripcion | VARCHAR | Contenido del comentario |
| puntuacion | SMALLINT | Puntuación otorgada a la receta |

id_comentario es la clave primaria.

id_usuario es una clave foránea que relaciona el comentario con el usuario que lo realizó.

id_receta es una clave foránea que relaciona el comentario con la receta correspondiente.

La puntuación está restringida a valores entre 1 y 5.

### Relaciones

Relacion 1 a N entre usuario y receta (Un usuario puede crear múltiples recetas).

Relacion 1 a N entre usuario y comentario (Un usuario puede realizar múltiples comentarios).

Relacion 1 a N entre receta y comentario (Una receta puede tener múltiples comentarios).

## Inicialización de la base de datos

La estructura de la base de datos se encuentra definida en: db/schema/schema.sql
Este archivo contiene la definición de las tablas, claves primarias, claves foráneas y restricciones de la base de datos.
PostgreSQL se ejecuta mediante el archivo: docker-compose.yml
Este archivo descarga la imagen de postgres versión 16:alpine
El archivo schema.sql se monta en el contenedor para que PostgreSQL pueda utilizarlo durante la inicialización de una base de datos nueva.

## Acceso a la base de datos desde Go

La aplicación utiliza el paquete estándar database/sql para trabajar con la base de datos.
Como driver de PostgreSQL se utiliza ¨github.com/jackc/pgx/v5¨
La conexión se realiza mediante sql.Open y se verifica utilizando db.Ping().

## Consultas SQL y sqlc

Las consultas SQL utilizadas por el proyecto se encuentran en el archivo db/queries/queries.sql
En este archivo se definen las operaciones CRUD de las tablas usuarios, recetas y comentarios.
Para evitar implementar manualmente el código de acceso a los datos, se utiliza sqlc, mediante el comando sqlc generate.
La configuración de sqlc se encuentra en el archivo sqlc.yaml
A partir de la definición de las tablas y las consultas SQL, sqlc genera automáticamente el código Go necesario. El código generado se encuentra en db/sqlc/
El código generado se elimina y vuelve a generar automáticamente durante la ejecución de los tests mediante el Makefile.

## Tests 

Los tests se encuentran en test/db_test.go, utilizan el paquete estándar testing, y se conectan directamente a PostgreSQL para verificar que las operaciones funcionen correctamente.

### Usuario

Se prueban:

- creación del usuario;
- obtención del usuario;
- obtención de la contraseña;
- listado de usuarios;
- modificación del nombre;
- modificación del apellido;
- modificación del email;
- modificación de la contraseña;
- eliminación del usuario;
- comprobación de que el usuario ya no existe después de eliminarlo.

### Receta

Se prueban:

- creación de la receta;
- obtención de la receta;
- listado de recetas;
- listado de recetas pertenecientes a un usuario;
- modificación del nombre;
- modificación de la descripción;
- modificación de los ingredientes;
- modificación de los pasos;
- eliminación de la receta;
- comprobación de que la receta ya no existe después de eliminarla.

### Comentario

Se prueban:

- creación del comentario;
- obtención del comentario;
- listado de comentarios de una receta;
- modificación de la descripción;
- modificación de la puntuación;
- eliminación del comentario;
- comprobación de que el comentario ya no existe después de eliminarlo.

## Automatización mediante Makefile

El proyecto incluye un Makefile para automatizar las tareas necesarias para ejecutar los tests.

El comando principal para ejecutar las pruebas es `make test`
El proceso realiza las siguientes tareas:
1.Eliminar código generado anteriormente, chequeando antes si hay codigo generado por el sqlc generate (rm -fr "db/sqlc")
2.Eliminar contenedores y volúmenes anteriores (docker compose down -v)
3.Generar nuevamente el código con sqlc (slqc generate)
4. Compilar el proyecto (go build ./...)
5. Levantar PostgreSQL (docker compose up -d)
6. Esperar a que PostgreSQL esté disponible (mediante la opcion --wait de docker compose up -d y el healtcheck que pusimos en el docker-compose.yml)
7. Ejecutar los tests, haciendo -count=1 para asegurarnos que haga el test y no utilice cache (go test -v -count=1 ./test)
8. Eliminar contenedor y volumenes (docker compose down -v)

Para estos dos ultimos pasos tuvimos que poner 
    "go test -v -count=1 ./test; \
    status=$$?; \
    docker compose down -v; \
    exit $$status" 
Ya que sino no se ejecutaba el docker compose down por segunda vez. De esta forma, nos aseguramos de que si lo haga, y de conservar el error si ocurre en el test. 

## Entorno de ejecución

Para ejecutar los tests de persistencia se requiere contar con:

- Go.
- Docker.
- sqlc.

El código generado por sqlc no necesita estar previamente generado para ejecutar el proceso de testing, ya que el `Makefile` lo genera automáticamente.

## 2. Capa de Aplicación y API REST

### Especificación de Endpoints

#### 1. Entidad: Usuario (`/usuario`)

| Método | Ruta | Descripción | Payload Request | Respuesta Exitosa | Respuestas de Error |
|---|---|---|---|---|---|
| **POST** | `/usuario` | Crea un nuevo usuario | `{"nombre", "apellido", "email", "contrasenia"}` | `201 Created` (datos del usuario sin password) | `400 Bad Request` |
| **GET** | `/usuario` | Lista todos los usuarios | Ninguno | `200 OK` (array JSON) | `500 Internal Server Error` |
| **GET** | `/usuario/{id}` | Obtiene el detalle de un usuario | Ninguno | `200 OK` (objeto usuario) | `400 Bad Request`, `404 Not Found` |
| **PUT** | `/usuario/{id}` | Actualiza nombre, apellido o email | `{"nombre", "apellido", "email"}` | `200 OK` (objeto actualizado) | `400 Bad Request`, `404 Not Found` |
| **DELETE** | `/usuario/{id}` | Elimina un usuario por ID | Ninguno | `204 No Content` | `400 Bad Request`, `500 Internal Server Error` |

#### 2. Entidad: Receta (`/receta`)

| Método | Ruta | Descripción | Payload Request | Respuesta Exitosa | Respuestas de Error |
|---|---|---|---|---|---|
| **POST** | `/receta` | Crea una nueva receta | `{"nombre", "descripcion", "ingredientes", "pasos", "id_usuario"}` | `201 Created` (objeto receta) | `400 Bad Request` |
| **GET** | `/receta` | Lista todas las recetas registradas | Ninguno | `200 OK` (array JSON) | `500 Internal Server Error` |
| **GET** | `/receta?usuario={id}` | Filtra recetas creadas por un usuario específico | Query param `usuario` | `200 OK` (array JSON) | `500 Internal Server Error` |
| **GET** | `/receta/{id}` | Obtiene el detalle completo de una receta | Ninguno | `200 OK` (objeto receta) | `400 Bad Request`, `404 Not Found` |
| **PUT** | `/receta/{id}` | Modifica el contenido de una receta | `{"nombre", "descripcion", "ingredientes", "pasos"}` | `200 OK` (objeto actualizado) | `400 Bad Request`, `404 Not Found` |
| **DELETE** | `/receta/{id}` | Elimina una receta por ID | Ninguno | `204 No Content` | `400 Bad Request`, `500 Internal Server Error` |

#### 3. Entidad: Comentario (`/comentario`)

| Método | Ruta | Descripción | Payload Request | Respuesta Exitosa | Respuestas de Error |
|---|---|---|---|---|---|
| **POST** | `/comentario` | Agrega un comentario y puntuación a una receta | `{"id_usuario", "id_receta", "descripcion", "puntuacion"}` | `201 Created` (objeto comentario) | `400 Bad Request` |
| **GET** | `/comentario` | Lista todos los comentarios de la plataforma | Ninguno | `200 OK` (array JSON) | `500 Internal Server Error` |
| **GET** | `/comentario?receta={id}` | Lista comentarios pertenecientes a una receta | Query param `receta` | `200 OK` (array JSON) | `500 Internal Server Error` |
| **GET** | `/comentario/{id}` | Obtiene un comentario por su ID | Ninguno | `200 OK` (objeto comentario) | `400 Bad Request`, `404 Not Found` |
| **PUT** | `/comentario/{id}` | Modifica la descripción y/o puntuación | `{"descripcion", "puntuacion"}` | `200 OK` (objeto actualizado) | `400 Bad Request`, `404 Not Found` |
| **DELETE** | `/comentario/{id}` | Elimina un comentario por ID | Ninguno | `204 No Content` | `400 Bad Request`, `500 Internal Server Error` |

### Validaciones y Manejo de Errores

La API valida exhaustivamente las entradas antes de interactuar con la base de datos:
1. **Campos Requeridos y Espacios en Blanco:** Se sanitizan cadenas mediante `strings.TrimSpace`. Campos esenciales como nombre, email, contraseña, títulos de recetas o descripciones no pueden enviarse vacíos.
2. **Restricción de Puntuación:** La puntuación en comentarios está validada tanto en código Go como por restricción `CHECK` en PostgreSQL para aceptar exclusivamente valores enteros entre **1 y 5**.
3. **Validación de Identificadores:** Cualquier ID pasado en la URL (`/{id}`) debe ser un entero positivo mayor a cero; de lo contrario, se responde de inmediato con `400 Bad Request`.
4. **Respuestas de Colecciones Vacías:** Para evitar que un listado vacío devuelva `null`, la API garantiza el retorno de corchetes vacíos `[]` con código `200 OK`.
5. **Códigos de Estado Estándar:**
   - `200 OK`: Operación de consulta o modificación exitosa.
   - `201 Created`: Recurso creado satisfactoriamente.
   - `204 No Content`: Eliminación exitosa sin cuerpo de retorno.
   - `400 Bad Request`: JSON malformado, parámetros faltantes o valores fuera de rango.
   - `404 Not Found`: Recurso solicitado no encontrado en la base de datos.
   - `405 Method Not Allowed`: Verbo HTTP no soportado en la ruta.
   - `500 Internal Server Error`: Falla interna no controlada.

### Pruebas de la API (`requests.hurl`)
Utilizando la herramienta **Hurl**, se ejecutan 23 casos de prueba:
1. Rechazo de creación de usuario por parámetros inválidos (`400`).
2. Creación exitosa de usuario (`201`) y captura de `id_usuario`.
3. Listado de usuarios (`200`) y verificación de tipo colección.
4. Consulta del usuario recién creado (`200`).
5. Actualización de datos del usuario (`200`).
6. Rechazo de creación de receta con datos faltantes (`400`).
7. Creación exitosa de receta asociada al usuario capturado (`201`) y captura de `id_receta`.
8. Listado de recetas (`200`) y detalle de receta (`200`).
9. Modificación de receta (`200`).
10. Rechazo de comentario con puntuación inválida (> 5) (`400`).
11. Creación exitosa de comentario (`201`) con captura de `id_comentario`.
12. Listado general de comentarios y listado filtrado por receta (`?receta=...`) (`200`).
13. Consulta de comentario por ID (`200`).
14. Rechazo de modificación de comentario con puntuación menor a 1 (`400`).
15. Modificación válida de comentario (`200`).
16. Eliminación de comentario (`204`) y verificación de inexistencia (`404`).
17. Eliminación de receta (`204`) y verificación de inexistencia (`404`).
18. Eliminación de usuario (`204`) y verificación de inexistencia (`404`).

## Automatización mediante Makefile

| Comando | Descripción |
|---|---|
| `make test-api` | Elimina código generado anterior, detiene y elimina contenedores y volúmenes, regenera código con `sqlc`, compila el proyecto, levanta la base y la API en Docker esperando que ambos estén listos (`--wait`), ejecuta las pruebas con `hurl --test requests.hurl`, y borra contenedores y volúmenes al finalizar. |
| `make test` | Ejecuta las pruebas unitarias de persistencia en Go (`tests/db_test.go`) contra el contenedor de base de datos. |
| `make up` | Levanta los contenedores. |
| `make down` | Detiene y borra los contenedores y volúmenes creados. |
| `make generate` | Ejecuta `sqlc generate`. |
| `make build` | Compila el código fuente de Go. |
