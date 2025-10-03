package main

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	Id, Name, Email string
	CreatedUnix     int64
}

type Repo struct {
	DB *sql.DB
}

func (r *Repo) Create(ctx context.Context, u *User) error {
	query := `INSERT INTO users (name, email, created_unix) VALUES ($1,$2,$3,$4) RETURNING id`

	args := []any{u.Name, u.Email, time.Now().Unix()}

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&u.Id)
}

func (r *Repo) Get(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, name, email, created_unix FROM users WHERE id = $1`

	var user User
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.CreatedUnix,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
