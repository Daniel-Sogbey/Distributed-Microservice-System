package main

import (
	"context"
	"database/sql"
	"time"
)

type Payment struct {
	Id          string `json:"id"`
	OrderId     string `json:"order_id"`
	Status      string `json:"status"`
	AmountCents int64  `json:"amount_cents"`
	CreatedUnix int64  `json:"created_unix"`
}

type Repo struct {
	DB *sql.DB
}

func (r Repo) Create(ctx context.Context, p *Payment) error {
	query := `INSERT INTO payments (order_id, status, amount_cents, created_unix) VALUES ($1,$2,$3,$4) RETURNING id`

	args := []any{p.OrderId, p.Status, p.AmountCents, time.Now().Unix()}
	return r.DB.QueryRowContext(ctx, query, args...).Scan(&p.Id)
}

func (r Repo) Get(ctx context.Context, id string) (*Payment, error) {
	query := `SELECT id, order_id, status, amount_cents, created_unix FROM payments WHERE id=$1`

	var payment Payment

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.Status,
		&payment.AmountCents,
		&payment.CreatedUnix,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}
