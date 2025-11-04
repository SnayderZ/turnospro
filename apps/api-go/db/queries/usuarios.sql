-- name: CreateUser :one
INSERT INTO usuarios (nombre, email, hash, rol)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM usuarios
WHERE email = $1;
