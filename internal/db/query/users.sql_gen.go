// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package query

import (
	"context"
)

const addUser = `-- name: AddUser :one
INSERT INTO users(name, email, password, role)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, password, role, created_at, updated_at
`

type AddUserParams struct {
	Name     string   `db:"name"`
	Email    string   `db:"email"`
	Password string   `db:"password"`
	Role     UserRole `db:"role"`
}

func (q *Queries) AddUser(ctx context.Context, arg AddUserParams) (User, error) {
	row := q.db.QueryRow(ctx, addUser,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const fetchAllUsers = `-- name: FetchAllUsers :many
SELECT id, name, email, password, role, created_at, updated_at
FROM users
`

func (q *Queries) FetchAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, fetchAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.Role,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchUser = `-- name: FetchUser :one
SELECT id, name, email, password, role, created_at, updated_at
FROM users
WHERE id = $1
`

func (q *Queries) FetchUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, fetchUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET name           = $1,
    email          = $2,
    role           = $3
WHERE id = $4
RETURNING id, name, email, password, role, created_at, updated_at
`

type UpdateUserParams struct {
	Name  string   `db:"name"`
	Email string   `db:"email"`
	Role  UserRole `db:"role"`
	ID    int32    `db:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Name,
		arg.Email,
		arg.Role,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users
SET password = $1
WHERE id = $2
RETURNING id, name, email, password, role, created_at, updated_at
`

type UpdateUserPasswordParams struct {
	Password string `db:"password"`
	ID       int32  `db:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserPassword, arg.Password, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
