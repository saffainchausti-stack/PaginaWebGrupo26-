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
1.Eliminar código generado anteriormente, chequeando antes si hay codigo generado por el sqlc generate (if [ -d "db/sqlc" ]; then rm -r "db/sqlc"; fi)
2.Eliminar contenedores y volúmenes anteriores (docker compose down -v)
3.Generar nuevamente el código con sqlc (slqc generate)
4. Compilar el proyecto (go build ./...)
5. Levantar PostgreSQL (docker compose up -d)
6. Esperar a que PostgreSQL esté disponible (mediante la opcion --wait de docker compose up -d y el healtcheck que pusimos en el docker-compose.yml)
7. Ejecutar los tests (go test -v ./test)
8. Eliminar contenedor y volumenes (docker compose down -v)

Para estos dos ultimos pasos tuvimos que poner 
    "go test -v ./test; \
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
