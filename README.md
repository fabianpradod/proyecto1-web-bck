# Series Tracker — Backend

API REST para gestionar series de televisión y sus ratings.

Construida en Go con PostgreSQL.

## Repositorios

- **Frontend:** https://github.com/fabianpradod/proyecto1-web-frt

## Tech Stack

- Go 1.21
- Chi (router)
- PostgreSQL
- godotenv

## Correr localmente

### Requisitos

- Go 1.21+
- PostgreSQL corriendo localmente

### Pasos

1. Clonar el repositorio
```bash
git clone https://github.com/fabianpradod/proyecto1-web-bck
cd proyecto1-web-bck
```

2. Crear la base de datos
```bash
psql -U TU_USUARIO -c "CREATE DATABASE seriestracker;"
```

3. Copiar el archivo de variables de entorno
```bash
cp .env.example .env
```

4. Editar `.env` con tu usuario de PostgreSQL
```bash
DATABASE_URL=postgresql://TU_USUARIO@localhost:5432/seriestracker?sslmode=disable
```

5. Correr el servidor
```bash
go run main.go
```

El servidor corre en `http://localhost:8080`. Las tablas se crean automáticamente al iniciar.

## Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /series | Listar series (soporta ?q=, ?sort=, ?order=, ?page=, ?limit=) |
| GET | /series/:id | Obtener serie por ID |
| POST | /series | Crear serie |
| PUT | /series/:id | Editar serie |
| DELETE | /series/:id | Eliminar serie |
| GET | /series/:id/rating | Obtener rating promedio |
| POST | /series/:id/rating | Agregar rating |
| GET | /docs | Swagger UI |
| GET | /openapi.yaml | Especificación OpenAPI |

## CORS

CORS está configurado para permitir todos los orígenes (`*`) durante desarrollo. Esto es necesario porque el cliente y el servidor corren en orígenes distintos (distintos puertos o dominios), y el navegador bloquea las peticiones fetch() por defecto si el servidor no lo permite explícitamente.

## Challenges implementados

- Códigos HTTP correctos (201 al crear, 204 al eliminar, 404 si no existe, 400 en input inválido)
- Validación server-side con respuestas de error en JSON
- Paginación con `?page=` y `?limit=`
- Búsqueda por nombre con `?q=`
- Ordenamiento con `?sort=` y `?order=`
- Especificación OpenAPI/Swagger escrita y precisa
- Swagger UI corriendo y siendo servido desde el backend
- Sistema de ratings con tabla propia y endpoints REST propios

## Reflexión

Go fue una experiencia interesante y la verdad comploicada para construir una API REST. El lenguaje es extremadamente especifico y requiere de mucho para hacer algo simple, cosas que en Python con SQLAlchemy las puedo hacer en una línea por su ORM, en Go con SQL toman como diez. Un ejemplo concreto: para hacer un `ORDER BY` dinámico sin vulnerabilidades de SQL injection, no se pueden usar parámetros, así que hay que construir la query como string y validar manualmente contra una lista blanca de columnas permitidas. Es explícito y seguro, pero detallado.

PostgreSQL fue la elección correcta, robusto, con soporte nativo en Go via `lib/pq`, y fácil de hostear en Render.

## Servidor en producción

API disponible en: https://proyecto1-web-bck.onrender.com
Swagger UI: https://proyecto1-web-bck.onrender.com/docs