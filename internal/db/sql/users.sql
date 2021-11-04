-- name: AddUser :one
INSERT INTO users(name, email, password, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET password = $1
WHERE id = $2
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name           = $1,
    email          = $2,
    role           = $3
WHERE id = $4
RETURNING *;

-- name: UpdateUserName :one
UPDATE users
SET name = $1
WHERE id = $2
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $1
WHERE id = $2
RETURNING *;

-- name: UpdateUserRole :one
UPDATE users
SET role = $1
WHERE id = $2
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;

-- name: FetchUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: FetchAllUsers :many
SELECT *
FROM users;
